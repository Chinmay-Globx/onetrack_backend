package domain

import "time"

// CreationMode drives how the bid was initialized — not its lifecycle
const (
	CreationModeManual       = "MANUAL"
	CreationModeIntelligence = "INTELLIGENCE"
)

// Workflow stages — shared by MANUAL and INTELLIGENCE modes
const (
	StageDiscovered            = "DISCOVERED"
	StageUnderAnalysis         = "UNDER_ANALYSIS"        // Intelligence mode only
	StageQualificationPending  = "QUALIFICATION_PENDING" // Intelligence mode only
	StageQualificationReview   = "QUALIFICATION_REVIEW"
	StageDocumentCompilation   = "DOCUMENT_COMPILATION"
	StageOEMCoordination       = "OEM_COORDINATION"
	StageCommercialPreparation = "COMMERCIAL_PREPARATION"
	StageInternalReview        = "INTERNAL_REVIEW"
	StageFinalApproval         = "FINAL_APPROVAL"
	StageReadyForSubmission    = "READY_FOR_SUBMISSION"
	StageSubmitted             = "SUBMITTED"
	StageRAActive              = "RA_ACTIVE"
	StageAwaitingResult        = "AWAITING_RESULT"
	StageWon                   = "WON"
	StageLost                  = "LOST"
	StageCancelled             = "CANCELLED"
)

// Bid status
const (
	BidStatusActive    = "ACTIVE"
	BidStatusCancelled = "CANCELLED"
	BidStatusArchived  = "ARCHIVED"
	BidStatusWon       = "WON"
	BidStatusLost      = "LOST"
)

// Bid outcome
const (
	OutcomeWon       = "WON"
	OutcomeLost      = "LOST"
	OutcomeCancelled = "CANCELLED"
)

// ────────────────────────────────────────
// Core entity
// ────────────────────────────────────────

type BidWorkspace struct {
	ID               string  `json:"id"`
	BidNo            *string `json:"bid_no,omitempty"`
	GemBidNo         *string `json:"gem_bid_no,omitempty"`
	Title            string  `json:"title"`
	OrganizationName *string `json:"organization_name,omitempty"`
	DepartmentName   *string `json:"department_name,omitempty"`
	PortalSource     string  `json:"portal_source"`
	CreationMode     string  `json:"creation_mode"`
	WorkflowStage    string  `json:"workflow_stage"`
	BidStatus        string  `json:"bid_status"`
	BidOwnerID       string  `json:"bid_owner_id"`
	CreatedBy        string  `json:"created_by"`

	// Financial
	EstimatedValue *float64 `json:"estimated_value,omitempty"`
	EMDAmount      *float64 `json:"emd_amount,omitempty"`
	EMDType        *string  `json:"emd_type,omitempty"`
	EMDExempted    bool     `json:"emd_exempted"`
	FinalBidValue  *float64 `json:"final_bid_value,omitempty"`
	L1Price        *float64 `json:"l1_price,omitempty"`
	QuotedPrice    *float64 `json:"quoted_price,omitempty"`

	// Dates
	OpeningDate    *time.Time `json:"opening_date,omitempty"`
	ClosingDate    *time.Time `json:"closing_date,omitempty"`
	SubmissionDate *time.Time `json:"submission_date,omitempty"`
	ResultDate     *time.Time `json:"result_date,omitempty"`
	RADate         *time.Time `json:"ra_date,omitempty"`

	// Categorization
	Category    *string `json:"category,omitempty"`
	BidType     *string `json:"bid_type,omitempty"`
	GemBidType  *string `json:"gem_bid_type,omitempty"`
	OEMRequired bool    `json:"oem_required"`
	HasTechEval bool    `json:"has_tech_eval"`

	// Outcome
	QualificationStatus  string  `json:"qualification_status"`
	BidOutcome           *string `json:"bid_outcome,omitempty"`
	OutcomeReason        *string `json:"outcome_reason,omitempty"`
	TechComplianceStatus *string `json:"tech_compliance_status,omitempty"`

	// Additional
	Remarks        *string `json:"remarks,omitempty"`
	CompetitorInfo []byte  `json:"-"` // raw JSONB
	Metadata       []byte  `json:"-"` // raw JSONB

	// AI fields (populated only for INTELLIGENCE mode)
	AISourceDocumentID     *string  `json:"ai_source_document_id,omitempty"`
	AIExtractionConfidence *float64 `json:"ai_extraction_confidence,omitempty"`

	// Timestamps
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at,omitempty"`
}

type BidWorkspaceMember struct {
	ID      string    `json:"id"`
	BidID   string    `json:"bid_id"`
	UserID  string    `json:"user_id"`
	Role    string    `json:"role"`
	AddedBy *string   `json:"added_by,omitempty"`
	AddedAt time.Time `json:"added_at"`
}

type BidStageHistory struct {
	ID               string    `json:"id"`
	BidID            string    `json:"bid_id"`
	FromStage        *string   `json:"from_stage,omitempty"`
	ToStage          string    `json:"to_stage"`
	TransitionReason *string   `json:"transition_reason,omitempty"`
	TransitionedBy   string    `json:"transitioned_by"`
	CreatedAt        time.Time `json:"created_at"`
}

// ────────────────────────────────────────
// Request DTOs
// ────────────────────────────────────────

type CreateBidRequest struct {
	CreationMode     string   `json:"creation_mode" binding:"required,oneof=MANUAL INTELLIGENCE"`
	Title            string   `json:"title" binding:"required"`
	BidNo            *string  `json:"bid_no"`
	GemBidNo         *string  `json:"gem_bid_no"`
	OrganizationName *string  `json:"organization_name"`
	DepartmentName   *string  `json:"department_name"`
	PortalSource     *string  `json:"portal_source"`
	BidType          *string  `json:"bid_type"`
	GemBidType       *string  `json:"gem_bid_type"`
	Category         *string  `json:"category"`
	EstimatedValue   *float64 `json:"estimated_value"`
	EMDAmount        *float64 `json:"emd_amount"`
	EMDType          *string  `json:"emd_type"`
	EMDExempted      *bool    `json:"emd_exempted"`
	OEMRequired      *bool    `json:"oem_required"`
	HasTechEval      *bool    `json:"has_tech_eval"`
	OpeningDate      *string  `json:"opening_date"`
	ClosingDate      *string  `json:"closing_date"`
	BidOwnerID       string   `json:"bid_owner_id" binding:"required"`
	Remarks          *string  `json:"remarks"`
	Metadata         *string  `json:"metadata"` // raw JSON string
}

type UpdateBidRequest struct {
	Title            *string  `json:"title"`
	BidNo            *string  `json:"bid_no"`
	GemBidNo         *string  `json:"gem_bid_no"`
	OrganizationName *string  `json:"organization_name"`
	DepartmentName   *string  `json:"department_name"`
	PortalSource     *string  `json:"portal_source"`
	BidType          *string  `json:"bid_type"`
	GemBidType       *string  `json:"gem_bid_type"`
	Category         *string  `json:"category"`
	EstimatedValue   *float64 `json:"estimated_value"`
	EMDAmount        *float64 `json:"emd_amount"`
	EMDType          *string  `json:"emd_type"`
	EMDExempted      *bool    `json:"emd_exempted"`
	OEMRequired      *bool    `json:"oem_required"`
	HasTechEval      *bool    `json:"has_tech_eval"`
	OpeningDate      *string  `json:"opening_date"`
	ClosingDate      *string  `json:"closing_date"`
	SubmissionDate   *string  `json:"submission_date"`
	ResultDate       *string  `json:"result_date"`
	RADate           *string  `json:"ra_date"`
	Remarks          *string  `json:"remarks"`
}

type TransitionStageRequest struct {
	TargetStage string  `json:"target_stage" binding:"required"`
	Reason      *string `json:"reason"`
}

type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=OWNER MANAGER MEMBER REVIEWER OBSERVER"`
}

type RecordOutcomeRequest struct {
	BidOutcome    string           `json:"bid_outcome" binding:"required,oneof=WON LOST CANCELLED"`
	FinalBidValue *float64         `json:"final_bid_value"`
	L1Price       *float64         `json:"l1_price"`
	QuotedPrice   *float64         `json:"quoted_price"`
	OutcomeReason *string          `json:"outcome_reason"`
	ResultDate    *string          `json:"result_date"`
	Competitors   []CompetitorInfo `json:"competitor_info"`
}

type CompetitorInfo struct {
	Name        string   `json:"name"`
	QuotedPrice *float64 `json:"quoted_price,omitempty"`
	Rank        *string  `json:"rank,omitempty"`
}

// ────────────────────────────────────────
// Response DTOs
// ────────────────────────────────────────

type UserSummary struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
}

type MemberResponse struct {
	UserID   string    `json:"user_id"`
	FullName string    `json:"full_name"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	AddedAt  time.Time `json:"added_at"`
}

type StageHistoryResponse struct {
	ID               string      `json:"id"`
	FromStage        *string     `json:"from_stage"`
	ToStage          string      `json:"to_stage"`
	TransitionReason *string     `json:"transition_reason"`
	TransitionedBy   UserSummary `json:"transitioned_by"`
	CreatedAt        time.Time   `json:"created_at"`
}

type BidResponse struct {
	ID                     string           `json:"id"`
	BidNo                  *string          `json:"bid_no"`
	GemBidNo               *string          `json:"gem_bid_no"`
	Title                  string           `json:"title"`
	OrganizationName       *string          `json:"organization_name"`
	DepartmentName         *string          `json:"department_name"`
	PortalSource           string           `json:"portal_source"`
	CreationMode           string           `json:"creation_mode"`
	WorkflowStage          string           `json:"workflow_stage"`
	BidStatus              string           `json:"bid_status"`
	EstimatedValue         *float64         `json:"estimated_value"`
	EMDAmount              *float64         `json:"emd_amount"`
	EMDType                *string          `json:"emd_type"`
	EMDExempted            bool             `json:"emd_exempted"`
	FinalBidValue          *float64         `json:"final_bid_value"`
	L1Price                *float64         `json:"l1_price"`
	QuotedPrice            *float64         `json:"quoted_price"`
	OpeningDate            *time.Time       `json:"opening_date"`
	ClosingDate            *time.Time       `json:"closing_date"`
	SubmissionDate         *time.Time       `json:"submission_date"`
	ResultDate             *time.Time       `json:"result_date"`
	RADate                 *time.Time       `json:"ra_date"`
	Category               *string          `json:"category"`
	BidType                *string          `json:"bid_type"`
	GemBidType             *string          `json:"gem_bid_type"`
	OEMRequired            bool             `json:"oem_required"`
	HasTechEval            bool             `json:"has_tech_eval"`
	QualificationStatus    string           `json:"qualification_status"`
	BidOutcome             *string          `json:"bid_outcome"`
	OutcomeReason          *string          `json:"outcome_reason"`
	TechComplianceStatus   *string          `json:"tech_compliance_status"`
	Remarks                *string          `json:"remarks"`
	CompetitorInfo         interface{}      `json:"competitor_info"`
	Metadata               interface{}      `json:"metadata"`
	BidOwner               UserSummary      `json:"bid_owner"`
	Members                []MemberResponse `json:"members"`
	CreatedBy              string           `json:"created_by"`
	AISourceDocumentID     *string          `json:"ai_source_document_id,omitempty"`
	AIExtractionConfidence *float64         `json:"ai_extraction_confidence,omitempty"`
	CreatedAt              time.Time        `json:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at"`
	ArchivedAt             *time.Time       `json:"archived_at"`
}

type BidListItem struct {
	ID               string      `json:"id"`
	BidNo            *string     `json:"bid_no"`
	GemBidNo         *string     `json:"gem_bid_no"`
	Title            string      `json:"title"`
	OrganizationName *string     `json:"organization_name"`
	Category         *string     `json:"category"`
	CreationMode     string      `json:"creation_mode"`
	WorkflowStage    string      `json:"workflow_stage"`
	BidStatus        string      `json:"bid_status"`
	BidOutcome       *string     `json:"bid_outcome"`
	EstimatedValue   *float64    `json:"estimated_value"`
	ClosingDate      *time.Time  `json:"closing_date"`
	OEMRequired      bool        `json:"oem_required"`
	BidOwner         UserSummary `json:"bid_owner"`
	CreatedAt        time.Time   `json:"created_at"`
}

type BidListResponse struct {
	Bids       []BidListItem `json:"bids"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

// ────────────────────────────────────────
// Repository-level insert type
// ────────────────────────────────────────

type CreateBidParams struct {
	BidNo            *string
	GemBidNo         *string
	Title            string
	OrganizationName *string
	DepartmentName   *string
	PortalSource     string
	CreationMode     string
	BidOwnerID       string
	CreatedBy        string
	EstimatedValue   *float64
	EMDAmount        *float64
	EMDType          *string
	EMDExempted      bool
	OEMRequired      bool
	HasTechEval      bool
	OpeningDate      *time.Time
	ClosingDate      *time.Time
	Category         *string
	BidType          *string
	GemBidType       *string
	Remarks          *string
	Metadata         []byte
	// AI-mode fields (nil for MANUAL)
	AISourceDocumentID     *string
	AIExtractionConfidence *float64
}

// ────────────────────────────────────────
// Query params
// ────────────────────────────────────────

type ListBidsParams struct {
	Page          int
	Limit         int
	Search        string
	WorkflowStage string
	BidStatus     string
	BidOutcome    string
	BidOwnerID    string
	Category      string
	CreationMode  string
	ClosingBefore *time.Time
	ClosingAfter  *time.Time
	OEMRequired   *bool
}
