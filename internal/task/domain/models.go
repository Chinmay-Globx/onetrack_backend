package domain

import "time"

// Task status values — shared by manual and AI-generated tasks
const (
	StatusOpen            = "OPEN"
	StatusAssigned        = "ASSIGNED"
	StatusInProgress      = "IN_PROGRESS"
	StatusWaitingExternal = "WAITING_EXTERNAL"
	StatusBlocked         = "BLOCKED"
	StatusUnderReview     = "UNDER_REVIEW"
	StatusCompleted       = "COMPLETED"
	StatusCancelled       = "CANCELLED"
	StatusEscalated       = "ESCALATED"
	StatusReopened        = "REOPENED"
)

// Task priority values
const (
	PriorityLow      = "LOW"
	PriorityMedium   = "MEDIUM"
	PriorityHigh     = "HIGH"
	PriorityCritical = "CRITICAL"
)

// Task types
const (
	TypeGeneral             = "GENERAL"
	TypeDocumentCollection  = "DOCUMENT_COLLECTION"
	TypeOEMCoordination     = "OEM_COORDINATION"
	TypeQualification       = "QUALIFICATION"
	TypeCommercial          = "COMMERCIAL"
	TypeApproval            = "APPROVAL"
	TypeReview              = "REVIEW"
	TypeSubmission          = "SUBMISSION"
	TypeChecklist           = "CHECKLIST"
	TypeCustom              = "CUSTOM"
	TypeAIGenerated         = "AI_GENERATED" // Used by intelligence mode
)

// Task source
const (
	SourceManual      = "MANUAL"
	SourceAIGenerated = "AI_GENERATED"
	SourceTemplate    = "TEMPLATE"
)

// Activity types
const (
	ActivityComment         = "COMMENT"
	ActivityStatusChanged   = "STATUS_CHANGED"
	ActivityAssigned        = "ASSIGNED"
	ActivityAttachmentAdded = "ATTACHMENT_ADDED"
	ActivityCallLogged      = "CALL_LOGGED"
	ActivityEmailLogged     = "EMAIL_LOGGED"
	ActivityEscalated       = "ESCALATED"
	ActivityDeadlineChanged = "DEADLINE_CHANGED"
)

// ────────────────────────────────────────
// Core entities
// ────────────────────────────────────────

type Task struct {
	ID                   string     `json:"id"`
	BidID                string     `json:"bid_id"`
	ParentTaskID         *string    `json:"parent_task_id,omitempty"`
	TaskType             string     `json:"task_type"`
	TaskCategory         *string    `json:"task_category,omitempty"`
	Title                string     `json:"title"`
	Description          *string    `json:"description,omitempty"`
	Status               string     `json:"status"`
	Priority             string     `json:"priority"`
	AssignedTo           *string    `json:"assigned_to,omitempty"`
	CreatedBy            string     `json:"created_by"`
	DueDate              *time.Time `json:"due_date,omitempty"`
	SLADeadline          *time.Time `json:"sla_deadline,omitempty"`
	CompletionPercentage float64    `json:"completion_percentage"`
	Source               string     `json:"source"`
	AIConfidence         *float64   `json:"ai_confidence,omitempty"`
	Metadata             []byte     `json:"-"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type TaskActivity struct {
	ID           string    `json:"id"`
	TaskID       string    `json:"task_id"`
	ActivityType string    `json:"activity_type"`
	ActivityData []byte    `json:"-"`
	PerformedBy  string    `json:"performed_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type TaskChecklist struct {
	ID        string     `json:"id"`
	TaskID    string     `json:"task_id"`
	Title     string     `json:"title"`
	IsDone    bool       `json:"is_done"`
	DoneBy    *string    `json:"done_by,omitempty"`
	DoneAt    *time.Time `json:"done_at,omitempty"`
	SortOrder int        `json:"sort_order"`
	CreatedAt time.Time  `json:"created_at"`
}

// ────────────────────────────────────────
// Request DTOs
// ────────────────────────────────────────

type CreateTaskRequest struct {
	TaskType    string   `json:"task_type" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Description *string  `json:"description"`
	Priority    *string  `json:"priority"`
	AssignedTo  *string  `json:"assigned_to"`
	DueDate     *string  `json:"due_date"`
	SLADeadline *string  `json:"sla_deadline"`
	Metadata    *string  `json:"metadata"`
}

type CreateSubtaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Priority    *string `json:"priority"`
	AssignedTo  *string `json:"assigned_to"`
	DueDate     *string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title                *string  `json:"title"`
	Description          *string  `json:"description"`
	Priority             *string  `json:"priority"`
	DueDate              *string  `json:"due_date"`
	SLADeadline          *string  `json:"sla_deadline"`
	CompletionPercentage *float64 `json:"completion_percentage"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AssignTaskRequest struct {
	AssignedTo string `json:"assigned_to" binding:"required"`
}

type AddActivityRequest struct {
	ActivityType string                 `json:"activity_type" binding:"required"`
	ActivityData map[string]interface{} `json:"activity_data" binding:"required"`
}

type AddChecklistRequest struct {
	Title     string `json:"title" binding:"required"`
	SortOrder *int   `json:"sort_order"`
}

type UpdateChecklistRequest struct {
	IsDone bool `json:"is_done"`
}

// ────────────────────────────────────────
// Response DTOs
// ────────────────────────────────────────

type UserSummary struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
}

type TaskResponse struct {
	ID                   string       `json:"id"`
	BidID                string       `json:"bid_id"`
	ParentTaskID         *string      `json:"parent_task_id"`
	TaskType             string       `json:"task_type"`
	TaskCategory         *string      `json:"task_category"`
	Title                string       `json:"title"`
	Description          *string      `json:"description"`
	Status               string       `json:"status"`
	Priority             string       `json:"priority"`
	AssignedTo           *UserSummary `json:"assigned_to"`
	CreatedBy            UserSummary  `json:"created_by"`
	DueDate              *time.Time   `json:"due_date"`
	SLADeadline          *time.Time   `json:"sla_deadline"`
	CompletionPercentage float64      `json:"completion_percentage"`
	Source               string       `json:"source"`
	AIConfidence         *float64     `json:"ai_confidence,omitempty"`
	Metadata             interface{}  `json:"metadata"`
	SubtaskCount         int          `json:"subtask_count"`
	ChecklistCount       int          `json:"checklist_count"`
	ChecklistDoneCount   int          `json:"checklist_done_count"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
}

type ActivityResponse struct {
	ID           string      `json:"id"`
	TaskID       string      `json:"task_id"`
	ActivityType string      `json:"activity_type"`
	ActivityData interface{} `json:"activity_data"`
	PerformedBy  UserSummary `json:"performed_by"`
	CreatedAt    time.Time   `json:"created_at"`
}

type ChecklistResponse struct {
	ID        string      `json:"id"`
	TaskID    string      `json:"task_id"`
	Title     string      `json:"title"`
	IsDone    bool        `json:"is_done"`
	DoneBy    *UserSummary `json:"done_by,omitempty"`
	DoneAt    *time.Time  `json:"done_at,omitempty"`
	SortOrder int         `json:"sort_order"`
	CreatedAt time.Time   `json:"created_at"`
}

type TaskListResponse struct {
	Tasks      []TaskResponse `json:"tasks"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type ActivityListResponse struct {
	Activities []ActivityResponse `json:"activities"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
}

// ────────────────────────────────────────
// Query params
// ────────────────────────────────────────

type ListTasksParams struct {
	BidID      string
	Status     string
	Priority   string
	AssignedTo string
	TaskType   string
	ParentOnly bool
	Page       int
	Limit      int
}

type ListActivitiesParams struct {
	Page  int
	Limit int
}

// ────────────────────────────────────────
// Repository-level insert types
// ────────────────────────────────────────

type CreateTaskParams struct {
	BidID        string
	ParentTaskID *string
	TaskType     string
	Title        string
	Description  *string
	Priority     string
	AssignedTo   *string
	CreatedBy    string
	DueDate      *time.Time
	SLADeadline  *time.Time
	Source       string
	AIConfidence *float64
	Metadata     []byte
}
