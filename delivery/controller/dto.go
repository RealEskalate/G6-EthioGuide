package controller

import (
	"EthioGuide/domain"
	"time"
)

type RegisterRequest struct {
	Name          string `json:"name" binding:"required"`
	Username      string `json:"username" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone"`
	Password      string `json:"password" binding:"required"`
	PreferredLang string `json:"preferredLang"`
}

type SocialLoginRequest struct {
	Provider domain.AuthProvider `json:"provider" binding:"required"`
	Code     string              `json:"code" binding:"required"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`

	ProfilePicture string      `json:"profile_picture,omitempty"`
	Role           domain.Role `json:"role"`
	IsVerified     bool        `json:"is_verified"`
	CreatedAt      time.Time   `json:"created_at"`
}

func toUserResponse(a *domain.Account) UserResponse {
	return UserResponse{
		ID:       a.ID,
		Name:     a.Name,
		Username: a.UserDetail.Username,
		Email:    a.Email,

		ProfilePicture: a.ProfilePicURL,
		Role:           a.Role,
		IsVerified:     a.UserDetail.IsVerified,
		CreatedAt:      a.CreatedAt,
	}
}

type TranslateDTO struct {
	Content string `json:"content" binding:"required"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ProcedureCreateRequest struct {
	Name           string `json:"name"`
	GroupID        string `json:"groupId"`
	OrganizationID string `json:"organizationId,omitempty"`

	// content
	Prerequisites []string `json:"prerequisites"`
	Steps         []string `json:"steps"`
	Result        string   `json:"result"`

	// Fees
	Label    string  `json:"label"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`

	// ProcessingTime
	MinDays int `json:"minDays"`
	MaxDays int `json:"maxDays"`
}

func toDomainProcedure(p *ProcedureCreateRequest) *domain.Procedure {
	content := domain.ProcedureContent{
		Prerequisites: p.Prerequisites,
		Steps:         p.Steps,
		Result:        p.Result,
	}
	fees := domain.ProcedureFee{
		Label:    p.Label,
		Currency: p.Currency,
		Amount:   p.Amount,
	}
	processingTime := domain.ProcessingTime{
		MinDays: p.MinDays,
		MaxDays: p.MaxDays,
	}

	return &domain.Procedure{
		GroupID:        &p.GroupID,
		OrganizationID: p.OrganizationID,
		Name:           p.Name,
		Content:        content,
		Fees:           fees,
		ProcessingTime: processingTime,
	}
}

type CreateCategoryRequest struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id" binding:"required"`
	ParentID       string `json:"parent_id"`
	Title          string `json:"title" binding:"required"`
}

type CategoryResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	ParentID       string `json:"parent_id,omitempty"`
	Title          string `json:"title"`
}

type PaginatedCategoryResponse struct {
	Data  []*CategoryResponse `json:"data"`
	Total int64               `json:"total"`
	Page  int64               `json:"page"`
	Limit int64               `json:"limit"`
}

// --- Feedback DTOs ---

// --- Request DTOs ---

type FeedbackCreateRequest struct {
	Content string   `json:"content" binding:"required"`
	Type    string   `json:"type" binding:"required,oneof=inaccuracy outdated thanks missing"`
	Tags    []string `json:"tags,omitempty"`
}

type FeedbackUpdateRequest struct {
	Content       *string   `json:"content,omitempty"`
	Type          *string   `json:"type,omitempty" binding:"omitempty,oneof=inaccuracy outdated thanks missing"`
	Status        *string   `json:"status,omitempty" binding:"omitempty,oneof=new in_progress resolved declined"`
	AdminResponse *string   `json:"admin_response,omitempty"`
	Tags          *[]string `json:"tags,omitempty"`
}

// --- Response DTOs ---

type FeedbackResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	ProcedureID   string    `json:"procedure_id"`
	Content       string    `json:"content"`
	LikeCount     int       `json:"like_count"`
	DislikeCount  int       `json:"dislike_count"`
	ViewCount     int       `json:"view_count"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	AdminResponse string    `json:"admin_response,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	CreatedAT     time.Time `json:"created_at"`
	UpdatedAT     time.Time `json:"updated_at"`
}

type FeedbackListResponse struct {
	Feedbacks []*FeedbackResponse `json:"feedbacks"`
	Total     int64               `json:"total"`
	Page      int64               `json:"page"`
	Limit     int64               `json:"limit"`
}

type FeedbackStatePatchRequest struct {
	Status        string `json:"status" binding:"required,oneof=new in_progress resolved declined"`
	AdminResponse string `json:"admin_response,omitempty"`
}

func fromDomainFeedback(f *domain.Feedback) *FeedbackResponse {
	return &FeedbackResponse{
		ID:            f.ID,
		UserID:        f.UserID,
		ProcedureID:   f.ProcedureID,
		Content:       f.Content,
		LikeCount:     f.LikeCount,
		DislikeCount:  f.DislikeCount,
		ViewCount:     f.ViewCount,
		Type:          string(f.Type),
		Status:        string(f.Status),
		AdminResponse: f.AdminResponse,
		Tags:          f.Tags,
		CreatedAT:     f.CreatedAT,
		UpdatedAT:     f.UpdatedAT,
	}
}

func toFeedbackListResponse(feedbacks []*domain.Feedback, total, page, limit int64) FeedbackListResponse {
	respFeedbacks := make([]*FeedbackResponse, len(feedbacks))
	for i, f := range feedbacks {
		respFeedbacks[i] = fromDomainFeedback(f)
	}
	return FeedbackListResponse{
		Feedbacks: respFeedbacks,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}
}

type CreatePostDTO struct {
	Title 	 string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Procedures []string `json:"procedures,omitempty"`
	Tags		[]string `json:"tags,omitempty"`
}

type UpdatePostDTO struct{
	Title 	 	string 		`json:"title"`
	Content     string   `json:"content"`
	Procedures []string `json:"procedures,omitempty"`
	Tags		[]string `json:"tags,omitempty"`
}