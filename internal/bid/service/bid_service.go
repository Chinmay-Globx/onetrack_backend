package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/onetrack/backend/internal/bid/domain"
)

type bidService struct {
	repo domain.BidRepository
}

func NewBidService(repo domain.BidRepository) domain.BidService {
	return &bidService{repo: repo}
}

func (s *bidService) CreateBid(ctx context.Context, req *domain.CreateBidRequest, createdBy string) (*domain.BidResponse, error) {
	params := &domain.CreateBidParams{
		BidNo:            req.BidNo,
		GemBidNo:         req.GemBidNo,
		Title:            req.Title,
		OrganizationName: req.OrganizationName,
		DepartmentName:   req.DepartmentName,
		PortalSource:     "GeM",
		CreationMode:     req.CreationMode,
		BidOwnerID:       req.BidOwnerID,
		CreatedBy:        createdBy,
		EstimatedValue:   req.EstimatedValue,
		EMDAmount:        req.EMDAmount,
		EMDType:          req.EMDType,
		Category:         req.Category,
		BidType:          req.BidType,
		GemBidType:       req.GemBidType,
		Remarks:          req.Remarks,
		Metadata:         []byte("{}"),
	}

	if req.PortalSource != nil {
		params.PortalSource = *req.PortalSource
	}
	if req.EMDExempted != nil {
		params.EMDExempted = *req.EMDExempted
	}
	if req.OEMRequired != nil {
		params.OEMRequired = *req.OEMRequired
	}
	if req.HasTechEval != nil {
		params.HasTechEval = *req.HasTechEval
	}
	if req.OpeningDate != nil {
		t, err := time.Parse(time.RFC3339, *req.OpeningDate)
		if err != nil {
			return nil, fmt.Errorf("invalid opening_date format: %w", err)
		}
		params.OpeningDate = &t
	}
	if req.ClosingDate != nil {
		t, err := time.Parse(time.RFC3339, *req.ClosingDate)
		if err != nil {
			return nil, fmt.Errorf("invalid closing_date format: %w", err)
		}
		params.ClosingDate = &t
	}
	if req.Metadata != nil {
		params.Metadata = []byte(*req.Metadata)
	}

	id, err := s.repo.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create bid: %w", err)
	}

	// Record initial stage history
	_ = s.repo.AddStageHistory(ctx, &domain.BidStageHistory{
		BidID:          id,
		ToStage:        domain.StageDiscovered,
		TransitionedBy: createdBy,
	})

	return s.GetBid(ctx, id)
}

func (s *bidService) GetBid(ctx context.Context, id string) (*domain.BidResponse, error) {
	bid, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	owner, err := s.repo.GetUserSummary(ctx, bid.BidOwnerID)
	if err != nil {
		owner = &domain.UserSummary{ID: bid.BidOwnerID}
	}

	members, err := s.repo.GetMembers(ctx, id)
	if err != nil {
		members = []domain.MemberResponse{}
	}

	return buildBidResponse(bid, owner, members), nil
}

func (s *bidService) ListBids(ctx context.Context, params domain.ListBidsParams) (*domain.BidListResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	bids, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	items := make([]domain.BidListItem, 0, len(bids))
	for _, b := range bids {
		owner, _ := s.repo.GetUserSummary(ctx, b.BidOwnerID)
		if owner == nil {
			owner = &domain.UserSummary{ID: b.BidOwnerID}
		}
		items = append(items, buildBidListItem(&b, owner))
	}

	return &domain.BidListResponse{
		Bids:       items,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}, nil
}

func (s *bidService) UpdateBid(ctx context.Context, id string, req *domain.UpdateBidRequest) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *bidService) TransitionStage(ctx context.Context, id string, req *domain.TransitionStageRequest, actorID string) (*domain.TransitionResult, error) {
	bid, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if IsTerminalStage(bid.WorkflowStage) {
		return nil, fmt.Errorf("bid is in a terminal stage: %s", bid.WorkflowStage)
	}

	if !IsTransitionAllowed(bid.CreationMode, bid.WorkflowStage, req.TargetStage) {
		return nil, fmt.Errorf("transition from %s to %s is not allowed for %s mode",
			bid.WorkflowStage, req.TargetStage, bid.CreationMode)
	}

	// Determine bid_status update
	newStatus := domain.BidStatusActive
	switch req.TargetStage {
	case domain.StageWon:
		newStatus = domain.BidStatusWon
	case domain.StageLost:
		newStatus = domain.BidStatusLost
	case domain.StageCancelled:
		newStatus = domain.BidStatusCancelled
	}

	prevStage := bid.WorkflowStage

	if err := s.repo.UpdateStage(ctx, id, req.TargetStage, newStatus); err != nil {
		return nil, fmt.Errorf("update stage: %w", err)
	}

	_ = s.repo.AddStageHistory(ctx, &domain.BidStageHistory{
		BidID:            id,
		FromStage:        &prevStage,
		ToStage:          req.TargetStage,
		TransitionReason: req.Reason,
		TransitionedBy:   actorID,
	})

	return &domain.TransitionResult{
		BidID:          id,
		PreviousStage:  prevStage,
		CurrentStage:   req.TargetStage,
		TransitionedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s *bidService) GetStageHistory(ctx context.Context, id string) ([]domain.StageHistoryResponse, error) {
	history, err := s.repo.GetStageHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]domain.StageHistoryResponse, 0, len(history))
	for _, h := range history {
		actor, _ := s.repo.GetUserSummary(ctx, h.TransitionedBy)
		if actor == nil {
			actor = &domain.UserSummary{ID: h.TransitionedBy}
		}
		result = append(result, domain.StageHistoryResponse{
			ID:               h.ID,
			FromStage:        h.FromStage,
			ToStage:          h.ToStage,
			TransitionReason: h.TransitionReason,
			TransitionedBy:   *actor,
			CreatedAt:        h.CreatedAt,
		})
	}
	return result, nil
}

func (s *bidService) AddMember(ctx context.Context, bidID string, req *domain.AddMemberRequest, actorID string) error {
	_, err := s.repo.GetByID(ctx, bidID)
	if err != nil {
		return err
	}
	return s.repo.AddMember(ctx, bidID, req.UserID, req.Role, actorID)
}

func (s *bidService) RemoveMember(ctx context.Context, bidID string, userID string) error {
	return s.repo.RemoveMember(ctx, bidID, userID)
}

func (s *bidService) RecordOutcome(ctx context.Context, id string, req *domain.RecordOutcomeRequest) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.UpdateOutcome(ctx, id, req)
}

func (s *bidService) ArchiveBid(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, id)
}

// ────────────────────────────────────────
// Response builders
// ────────────────────────────────────────

func buildBidResponse(bid *domain.BidWorkspace, owner *domain.UserSummary, members []domain.MemberResponse) *domain.BidResponse {
	var competitorInfo interface{} = []interface{}{}
	var metadata interface{} = map[string]interface{}{}

	if len(bid.CompetitorInfo) > 0 {
		_ = json.Unmarshal(bid.CompetitorInfo, &competitorInfo)
	}
	if len(bid.Metadata) > 0 {
		_ = json.Unmarshal(bid.Metadata, &metadata)
	}

	return &domain.BidResponse{
		ID:               bid.ID,
		BidNo:            bid.BidNo,
		GemBidNo:         bid.GemBidNo,
		Title:            bid.Title,
		OrganizationName: bid.OrganizationName,
		DepartmentName:   bid.DepartmentName,
		PortalSource:     bid.PortalSource,
		CreationMode:     bid.CreationMode,
		WorkflowStage:    bid.WorkflowStage,
		BidStatus:        bid.BidStatus,
		EstimatedValue:   bid.EstimatedValue,
		EMDAmount:        bid.EMDAmount,
		EMDType:          bid.EMDType,
		EMDExempted:      bid.EMDExempted,
		FinalBidValue:    bid.FinalBidValue,
		L1Price:          bid.L1Price,
		QuotedPrice:      bid.QuotedPrice,
		OpeningDate:      bid.OpeningDate,
		ClosingDate:      bid.ClosingDate,
		SubmissionDate:   bid.SubmissionDate,
		ResultDate:       bid.ResultDate,
		RADate:           bid.RADate,
		Category:         bid.Category,
		BidType:          bid.BidType,
		GemBidType:       bid.GemBidType,
		OEMRequired:      bid.OEMRequired,
		HasTechEval:      bid.HasTechEval,
		QualificationStatus:    bid.QualificationStatus,
		BidOutcome:       bid.BidOutcome,
		OutcomeReason:    bid.OutcomeReason,
		TechComplianceStatus:   bid.TechComplianceStatus,
		Remarks:          bid.Remarks,
		CompetitorInfo:   competitorInfo,
		Metadata:         metadata,
		BidOwner:         *owner,
		Members:          members,
		CreatedBy:        bid.CreatedBy,
		AISourceDocumentID:     bid.AISourceDocumentID,
		AIExtractionConfidence: bid.AIExtractionConfidence,
		CreatedAt:        bid.CreatedAt,
		UpdatedAt:        bid.UpdatedAt,
		ArchivedAt:       bid.ArchivedAt,
	}
}

func buildBidListItem(bid *domain.BidWorkspace, owner *domain.UserSummary) domain.BidListItem {
	return domain.BidListItem{
		ID:               bid.ID,
		BidNo:            bid.BidNo,
		GemBidNo:         bid.GemBidNo,
		Title:            bid.Title,
		OrganizationName: bid.OrganizationName,
		Category:         bid.Category,
		CreationMode:     bid.CreationMode,
		WorkflowStage:    bid.WorkflowStage,
		BidStatus:        bid.BidStatus,
		BidOutcome:       bid.BidOutcome,
		EstimatedValue:   bid.EstimatedValue,
		ClosingDate:      bid.ClosingDate,
		OEMRequired:      bid.OEMRequired,
		BidOwner:         *owner,
		CreatedAt:        bid.CreatedAt,
	}
}
