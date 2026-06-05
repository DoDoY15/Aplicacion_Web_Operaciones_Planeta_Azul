package models

import (
	"time"

	"github.com/google/uuid"
)

// Enums

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleChefGeral  UserRole = "chefe_geral"
	RoleChefArea   UserRole = "chefe_area"
	RoleSupervisor UserRole = "supervisor"
	RoleMembro     UserRole = "membro"
)

type ItemStatus string

const (
	StatusDraft      ItemStatus = "draft"
	StatusPending    ItemStatus = "pending"
	StatusInProgress ItemStatus = "in_progress"
	StatusWaiting    ItemStatus = "waiting"
	StatusDone       ItemStatus = "done"
	StatusRejected   ItemStatus = "rejected"
)

type ItemVisibility string

const (
	VisibilityPrivate ItemVisibility = "private"
	VisibilityTeam    ItemVisibility = "team"
	VisibilityPublic  ItemVisibility = "public"
)

type ItemPriority string

const (
	PriorityLow    ItemPriority = "low"
	PriorityMedium ItemPriority = "medium"
	PriorityHigh   ItemPriority = "high"
	PriorityUrgent ItemPriority = "urgent"
)

type ApprovalDecision string

const (
	DecisionApproved ApprovalDecision = "approved"
	DecisionRejected ApprovalDecision = "rejected"
)

type InteractionStatus string

const (
	InteractionOpen     InteractionStatus = "open"
	InteractionResolved InteractionStatus = "resolved"
)

// Domain models

type Area struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type User struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         UserRole   `json:"role"`
	Cargo        string     `json:"cargo,omitempty"`
	AreaID       *uuid.UUID `json:"area_id,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`

	// Populated joins
	Area *Area `json:"area,omitempty"`
}

type Item struct {
	ID               uuid.UUID      `json:"id"`
	ParentID         *uuid.UUID     `json:"parent_id,omitempty"`
	Title            string         `json:"title"`
	Description      string         `json:"description,omitempty"`
	CreatedBy        uuid.UUID      `json:"created_by"`
	AssignedTo       *uuid.UUID     `json:"assigned_to,omitempty"`
	AreaID           uuid.UUID      `json:"area_id"`
	DeletedBy        *uuid.UUID     `json:"deleted_by,omitempty"`
	Status           ItemStatus     `json:"status"`
	Visibility       ItemVisibility `json:"visibility"`
	RequiresApproval bool           `json:"requires_approval"`
	Priority         ItemPriority   `json:"priority"`
	Deadline         *time.Time     `json:"deadline,omitempty"`
	CompletedAt      *time.Time     `json:"completed_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        *time.Time     `json:"deleted_at,omitempty"`

	// Populated joins
	Creator    *User        `json:"creator,omitempty"`
	Assignee   *User        `json:"assignee,omitempty"`
	Children   []*Item      `json:"children,omitempty"`
	Comments   []*Comment   `json:"comments,omitempty"`
	Approvals  []*Approval  `json:"approvals,omitempty"`
}

type ItemAssignment struct {
	ID         uuid.UUID `json:"id"`
	ItemID     uuid.UUID `json:"item_id"`
	UserID     uuid.UUID `json:"user_id"`
	RoleInItem string    `json:"role_in_item"`
	AssignedAt time.Time `json:"assigned_at"`
	AssignedBy uuid.UUID `json:"assigned_by"`
}

type Interaction struct {
	ID          uuid.UUID         `json:"id"`
	ItemID      uuid.UUID         `json:"item_id"`
	OpenedBy    uuid.UUID         `json:"opened_by"`
	AddressedTo uuid.UUID         `json:"addressed_to"`
	Message     string            `json:"message"`
	Response    string            `json:"response,omitempty"`
	Status      InteractionStatus `json:"status"`
	ResolvedAt  *time.Time        `json:"resolved_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

type Comment struct {
	ID        uuid.UUID `json:"id"`
	ItemID    uuid.UUID `json:"item_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	Author *User `json:"author,omitempty"`
}

type Approval struct {
	ID         uuid.UUID        `json:"id"`
	ItemID     uuid.UUID        `json:"item_id"`
	Reviewer   uuid.UUID        `json:"reviewer"`
	Decision   ApprovalDecision `json:"decision"`
	Note       string           `json:"note,omitempty"`
	DecidedAt  time.Time        `json:"decided_at"`
}

type AreaAccess struct {
	ID        uuid.UUID  `json:"id"`
	GrantedBy uuid.UUID  `json:"granted_by"`
	GrantedTo uuid.UUID  `json:"granted_to"`
	AreaID    uuid.UUID  `json:"area_id"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	RefID     uuid.UUID `json:"ref_id"`
	RefType   string    `json:"ref_type"`
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// Request/Response DTOs

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type CreateItemRequest struct {
	ParentID         *uuid.UUID     `json:"parent_id"`
	Title            string         `json:"title" binding:"required,min=3,max=255"`
	Description      string         `json:"description"`
	AssignedTo       *uuid.UUID     `json:"assigned_to"`
	Visibility       ItemVisibility `json:"visibility" binding:"required"`
	RequiresApproval bool           `json:"requires_approval"`
	Priority         ItemPriority   `json:"priority"`
	Deadline         *time.Time     `json:"deadline"`
}

type UpdateItemRequest struct {
	Title       *string        `json:"title"`
	Description *string        `json:"description"`
	AssignedTo  *uuid.UUID     `json:"assigned_to"`
	Status      *ItemStatus    `json:"status"`
	Priority    *ItemPriority  `json:"priority"`
	Deadline    *time.Time     `json:"deadline"`
	Visibility  *ItemVisibility `json:"visibility"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type CreateInteractionRequest struct {
	AddressedTo uuid.UUID `json:"addressed_to" binding:"required"`
	Message     string    `json:"message" binding:"required,min=10"`
}

type ResolveInteractionRequest struct {
	Response string `json:"response" binding:"required,min=1"`
}

type ApprovalRequest struct {
	Decision ApprovalDecision `json:"decision" binding:"required"`
	Note     string           `json:"note"`
}

// JWT Claims

type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Role   UserRole `json:"role"`
	AreaID string   `json:"area_id"`
}
