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

func ToUserResponse(a *domain.Account) UserResponse {
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

type UpdateOrgRequest struct {
	Name          *string            `json:"name"`
	ProfilePicURL *string            `json:"profile_pic_url"`
	Description   *string            `json:"description"`
	Location      *string            `json:"location"`
	ContactInfo   *ContactInfoUpdate `json:"contact_info"`
	PhoneNumbers  []string           `json:"phone_numbers"`
}

type TranslateDTO struct {
	Content string `json:"content" binding:"required"`
}
type ChatRequest struct {
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

// ------------------------------------
type UserUpdateRequest struct {
	Name          *string `json:"name,omitempty"`
	Email         *string `json:"email,omitempty"`
	ProfilePicURL *string `json:"profilePicURL,omitempty"`

	// Mutually exclusive blocks
	UserDetail         *UserDetailUpdate         `json:"userDetail,omitempty"`
	OrganizationDetail *OrganizationDetailUpdate `json:"organizationDetail,omitempty"`
}

type UserDetailUpdate struct {
	Username *string `json:"username,omitempty"`
}

type OrganizationDetailUpdate struct {
	Description  *string            `json:"description,omitempty"`
	Location     *string            `json:"location,omitempty"`
	Type         *string            `json:"type,omitempty"` // "gov" | "private"
	ContactInfo  *ContactInfoUpdate `json:"contactInfo,omitempty"`
	PhoneNumbers *[]string          `json:"phoneNumbers,omitempty"`
}

type ContactInfoUpdate struct {
	Socials *map[string]string `json:"socials,omitempty"`
	Website *string            `json:"website,omitempty"`
}

func ToDomainAccountUpdate(req *UserUpdateRequest, existing *domain.Account) *domain.Account {
	account := *existing // copy existing account

	if req.Name != nil {
		account.Name = *req.Name
	}
	if req.ProfilePicURL != nil {
		account.ProfilePicURL = *req.ProfilePicURL
	}
	if req.Email != nil {
		account.Email = *req.Email
	}

	// User detail
	if req.UserDetail != nil {
		if account.OrganizationDetail != nil {
			// mutual exclusion will be validated in usecase
		}
		if account.UserDetail == nil {
			account.UserDetail = &domain.UserDetail{}
		}
		if req.UserDetail.Username != nil {
			account.UserDetail.Username = *req.UserDetail.Username
		}
	}

	// Organization detail
	if req.OrganizationDetail != nil {
		if account.UserDetail != nil {
			// mutual exclusion will be validated in usecase
		}
		if account.OrganizationDetail == nil {
			account.OrganizationDetail = &domain.OrganizationDetail{}
		}

		od := req.OrganizationDetail
		if od.Description != nil {
			account.OrganizationDetail.Description = *od.Description
		}
		if od.Location != nil {
			account.OrganizationDetail.Location = *od.Location
		}
		if od.Type != nil {
			account.OrganizationDetail.Type = domain.OrganizationType(*od.Type)
		}
		if od.ContactInfo != nil {
			if od.ContactInfo.Website != nil {
				account.OrganizationDetail.ContactInfo.Website = *od.ContactInfo.Website
			}
			if od.ContactInfo.Socials != nil {
				account.OrganizationDetail.ContactInfo.Socials = *od.ContactInfo.Socials
			}
		}
		if od.PhoneNumbers != nil {
			account.OrganizationDetail.PhoneNumbers = *od.PhoneNumbers
		}
	}

	return &account
}

// ------------------------------------

type ProcedureCreateRequest struct {
	Name           string `json:"name"`
	GroupID        string `json:"groupId"`
	OrganizationID string `json:"organizationId,omitempty"`

	// content
	Prerequisites []string       `json:"prerequisites"`
	Steps         map[int]string `json:"steps"`
	Result        string         `json:"result"`

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
	ID string `json:"id"`
	// OrganizationID string `json:"organization_id" binding:"required"`
	OrganizationID string `json:"organization_id"`
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
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Procedures []string `json:"procedures,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

type PaginatedPostsResponse struct {
	Posts []*domain.Post `json:"posts"`
	Total int64          `json:"total"`
	Page  int64          `json:"page"`
	Limit int64          `json:"limit"`
}

type UpdatePostDTO struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Procedures []string `json:"procedures,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

type ForgotDTO struct {
	Email string `json:"email"`
}

type ResetDTO struct {
	ResetToken  string `json:"resetToken"`
	NewPassword string `json:"new_password"`
}

type ActivateDTO struct {
	ActivateToken string `json:"activationToken"`
}

// //////////////////////////// procedure////
type ProcedureContentResponse struct {
	Prerequisites []string       `json:"prerequisites"`
	Steps         map[int]string `json:"steps"`
	Result        string         `json:"result"`
}

type ProcedureFeeResponse struct {
	Label    string  `json:"label"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type ProcessingTimeResponse struct {
	MinDays int `json:"minDays"`
	MaxDays int `json:"maxDays"`
}

type ProcedureResponse struct {
	ID             string                   `json:"id"`
	GroupID        *string                  `json:"groupId,omitempty"`
	OrganizationID string                   `json:"organizationId"`
	Name           string                   `json:"name"`
	Content        ProcedureContentResponse `json:"content"`
	Fees           ProcedureFeeResponse     `json:"fees"`
	ProcessingTime ProcessingTimeResponse   `json:"processingTime"`
	CreatedAt      time.Time                `json:"createdAt"`
	NoticeIDs      []string                 `json:"noticeIds"`
}

type Pagination struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type PaginatedProcedureResponse struct {
	Data       []ProcedureResponse `json:"data"`
	Pagination Pagination          `json:"pagination"`
}

// ====== Mappers ======

func toProcedureResponse(p *domain.Procedure) ProcedureResponse {
	return ProcedureResponse{
		ID:             p.ID,
		GroupID:        p.GroupID,
		OrganizationID: p.OrganizationID,
		Name:           p.Name,
		Content: ProcedureContentResponse{
			Prerequisites: p.Content.Prerequisites,
			Steps:         p.Content.Steps,
			Result:        p.Content.Result,
		},
		Fees: ProcedureFeeResponse{
			Label:    p.Fees.Label,
			Currency: p.Fees.Currency,
			Amount:   p.Fees.Amount,
		},
		ProcessingTime: ProcessingTimeResponse{
			MinDays: p.ProcessingTime.MinDays,
			MaxDays: p.ProcessingTime.MaxDays,
		},
		CreatedAt: p.CreatedAt,
		NoticeIDs: p.NoticeIDs,
	}
}

func toPaginatedProcedureResponse(procedures []*domain.Procedure, total, page, limit int64) PaginatedProcedureResponse {
	responses := make([]ProcedureResponse, len(procedures))
	for i, p := range procedures {
		responses[i] = toProcedureResponse(p)
	}

	return PaginatedProcedureResponse{
		Data: responses,
		Pagination: Pagination{
			Total: total,
			Page:  page,
			Limit: limit,
		},
	}
}

type PreferencesDTO struct {
	PreferredLang     string `json:"preferredLang"`
	PushNotification  bool   `json:"pushNotification"`
	EmailNotification bool   `json:"emailNotification"`
}

// notice related dto
type NoticeDTO struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	Tags           []string  `json:"tags"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (dto *NoticeDTO) ToDomainNotice() *domain.Notice {
	return &domain.Notice{
		ID:             dto.ID,
		OrganizationID: dto.OrganizationID,
		Title:          dto.Title,
		Content:        dto.Content,
		Tags:           dto.Tags,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
	}
}

func FromDomainNotice(n *domain.Notice) *NoticeDTO {
	return &NoticeDTO{
		ID:             n.ID,
		OrganizationID: n.OrganizationID,
		Title:          n.Title,
		Content:        n.Content,
		Tags:           n.Tags,
		CreatedAt:      n.CreatedAt,
		UpdatedAt:      n.UpdatedAt,
	}
}

// Paginated response for notices
type NoticeListResponse struct {
	Data  []NoticeDTO `json:"data"`
	Total int64       `json:"total"`
	Page  int64       `json:"page"`
	Limit int64       `json:"limit"`
}

type OrgCreateRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	OrgType string `json:"type"`
}

type ContactInfoDTO struct {
	Socials map[string]string `json:"socials,omitempty"`
	Website string            `json:"website,omitempty"`
}

type OrganizationDetailDTO struct {
	Description  string         `json:"description,omitempty"`
	Location     string         `json:"location,omitempty"`
	Type         string         `json:"type,omitempty"` // assuming OrganizationType is a custom enum
	ContactInfo  ContactInfoDTO `json:"contact_info,omitempty"`
	PhoneNumbers []string       `json:"phone_numbers,omitempty"`
}

type OrganizationResponseDTO struct {
	ID                 string                `json:"id"`
	Name               string                `json:"name,omitempty"`
	Email              string                `json:"email"`
	ProfilePicURL      string                `json:"profile_pic_url,omitempty"`
	Role               string                `json:"role"` // assuming domain.Role is a string or enum
	CreatedAt          time.Time             `json:"created_at"`
	OrganizationDetail OrganizationDetailDTO `json:"organization_detail,omitempty"`
}

func ToOrganizationDTO(account *domain.Account) OrganizationResponseDTO {
	return OrganizationResponseDTO{
		ID:            account.ID,
		Name:          account.Name,
		Email:         account.Email,
		ProfilePicURL: account.ProfilePicURL,
		Role:          string(account.Role),
		CreatedAt:     account.CreatedAt,
		OrganizationDetail: OrganizationDetailDTO{
			Description: account.OrganizationDetail.Description,
			Location:    account.OrganizationDetail.Location,
			Type:        string(account.OrganizationDetail.Type),
			ContactInfo: ContactInfoDTO{
				Socials: account.OrganizationDetail.ContactInfo.Socials,
				Website: account.OrganizationDetail.ContactInfo.Website,
			},
			PhoneNumbers: account.OrganizationDetail.PhoneNumbers,
		},
	}
}

type PaginatedOrgsResponse struct {
	Total    int64 `json:"total"`
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

func ToOrganizationDetailDTO(account *domain.Account) OrganizationDetailDTO {
	return OrganizationDetailDTO{
		Description: account.OrganizationDetail.Description,
		Location:    account.OrganizationDetail.Location,
		Type:        string(account.OrganizationDetail.Type),
		ContactInfo: ContactInfoDTO{
			Socials: account.OrganizationDetail.ContactInfo.Socials,
			Website: account.OrganizationDetail.ContactInfo.Website,
		},
		PhoneNumbers: account.OrganizationDetail.PhoneNumbers,
	}
}

// =================================================================================
// --- Search DTOs ---
// =================================================================================

// SearchResultResponse defines the response for a search query.
type SearchResultResponse struct {
	Procedures    []*domain.Procedure         `json:"procedures"`
	Organizations []*AccountOrgSearchResponse `json:"organizations"`
}

// AccountOrgSearchResponse defines the public-facing information for an organization in search results.
type AccountOrgSearchResponse struct {
	ID                 string                      `json:"id"`
	Name               string                      `json:"name"`
	Email              string                      `json:"email"`
	ProfilePicURL      string                      `json:"profile_pic_url,omitempty"`
	Role               domain.Role                 `json:"role"`
	OrganizationDetail *OrganizationDetailResponse `json:"organization_detail,omitempty"`
}

// OrganizationDetailResponse defines the public-facing organization details.
type OrganizationDetailResponse struct {
	Description  string                  `json:"description,omitempty"`
	Location     string                  `json:"location,omitempty"`
	Type         domain.OrganizationType `json:"type,omitempty"`
	ContactInfo  *ContactInfoResponse    `json:"contact_info,omitempty"`
	PhoneNumbers []string                `json:"phone_numbers,omitempty"`
}

// ContactInfoResponse defines the public-facing contact information.
type ContactInfoResponse struct {
	Socials map[string]string `json:"socials,omitempty"`
	Website string            `json:"website,omitempty"`
}

// =================================================================================
// --- Search Mapper Functions ---
// =================================================================================

// ToSearchJSON converts the domain search result to a JSON-friendly response.
func ToSearchJSON(sr *domain.SearchResult) *SearchResultResponse {
	orgs := make([]*AccountOrgSearchResponse, 0, len(sr.Organizations))
	for _, acc := range sr.Organizations {
		orgs = append(orgs, toAccountOrgSearchResponse(acc))
	}
	return &SearchResultResponse{
		Procedures:    sr.Procedures,
		Organizations: orgs,
	}
}

// toAccountOrgSearchResponse converts a domain search result for an org to a response DTO.
// (Renamed from ToAccountResponse for clarity)
func toAccountOrgSearchResponse(acc *domain.AccountOrgSearch) *AccountOrgSearchResponse {
	if acc == nil {
		return nil
	}
	return &AccountOrgSearchResponse{
		ID:                 acc.ID,
		Name:               acc.Name,
		Email:              acc.Email,
		ProfilePicURL:      acc.ProfilePicURL,
		Role:               acc.Role,
		OrganizationDetail: toOrganizationDetailResponse(acc.OrganizationDetail),
	}
}

// toOrganizationDetailResponse converts a domain org detail to a response DTO.
func toOrganizationDetailResponse(od *domain.OrganizationDetail) *OrganizationDetailResponse {
	if od == nil {
		return nil
	}
	return &OrganizationDetailResponse{
		Description:  od.Description,
		Location:     od.Location,
		Type:         od.Type,
		ContactInfo:  toContactInfoResponse(&od.ContactInfo),
		PhoneNumbers: od.PhoneNumbers,
	}
}

// toContactInfoResponse makes the mapping from domain to DTO explicit and safe.
func toContactInfoResponse(ci *domain.ContactInfo) *ContactInfoResponse {
	if ci == nil {
		return nil
	}
	return &ContactInfoResponse{
		Socials: ci.Socials,
		Website: ci.Website,
	}
}

type CreateChecklistRequest struct {
	ProcedureID string `json:"procedure_id"`
}

type UpdateChecklistRequest struct {
	ChecklistID string `json:"checklist_id"`
}

type GetChecklistByID struct {
	UserProcedureID string `json:"user_procedure_id"`
}

type UserProcedureResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProcedureID string    `json:"procedure_id"`
	Percent     int       `json:"percent"`
	Status      string    `json:"status"` // e.g., "Not Started", "In Progress", "Completed"
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToControllerUserProcedure(user *domain.UserProcedure) *UserProcedureResponse {
	return &UserProcedureResponse{
		ID:          user.ID,
		UserID:      user.UserID,
		ProcedureID: user.ProcedureID,
		Percent:     user.Percent,
		Status:      user.Status,
		UpdatedAt:   user.UpdatedAt,
	}
}

type ChecklistResponse struct {
	ID              string `json:"id"`
	UserProcedureID string `json:"user_procedure_id"`
	Type            string `json:"type"`       // "Requirement" or "Step"
	Content         string `json:"content"`    // The actual checklist item
	IsChecked       bool   `json:"is_checked"` // Whether the item is completed
}

func ToControllerChecklist(check *domain.Checklist) *ChecklistResponse {
	return &ChecklistResponse{
		ID:              check.ID,
		UserProcedureID: check.UserProcedureID,
		Type:            check.Type,
		Content:         check.Content,
		IsChecked:       check.IsChecked,
	}
}

type AIChatRequest struct {
	Query string `json:"query" binding:"required"`
}

type AIConversationResponse struct {
	ID                string                 `json:"id,omitempty"`
	UserID            string                 `json:"user_id"`
	Source            string                 `json:"source,omitempty"`
	Request           string                 `json:"request"`
	Response          string                 `json:"response"`
	Timestamp         time.Time              `json:"timestamp"`
	RelatedProcedures []*AIProcedureResponse `json:"procedures"`
}

type AIProcedureResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ToAIProcedureResponse(procedures *domain.AIProcedure) *AIProcedureResponse {
	return &AIProcedureResponse{
		ID:   procedures.Id,
		Name: procedures.Name,
	}
}

func ToAIChatResponse(answer *domain.AIChat) *AIConversationResponse {
	procedures := make([]*AIProcedureResponse, len(answer.RelatedProcedures))
	for i, proc := range answer.RelatedProcedures {
		procedures[i] = ToAIProcedureResponse(proc)
	}
	return &AIConversationResponse{
		ID:                answer.ID,
		UserID:            answer.UserID,
		Source:            answer.Source,
		Request:           answer.Request,
		Response:          answer.Response,
		Timestamp:         answer.Timestamp,
		RelatedProcedures: procedures,
	}
}

type PaginatedAIHisoryResponse struct {
	Data       []AIConversationResponse `json:"data"`
	Pagination Pagination               `json:"pagination"`
}

func toPaginatedAIHisory(aiHistory []*domain.AIChat, total, page, limit int64) *PaginatedAIHisoryResponse {
	conv := make([]AIConversationResponse, len(aiHistory))
	for i, conversation := range aiHistory {
		conv[i] = *ToAIChatResponse(conversation)
	}

	return &PaginatedAIHisoryResponse{
		Data: conv,
		Pagination: Pagination{
			Total: total,
			Page:  page,
			Limit: limit,
		},
	}
}

type OrgsListPaginated struct {
	Orgs     []OrganizationResponseDTO `json:"orgs"`
	Total    int64                     `json:"total"`
	Page     int64                     `json:"page"`
	PageSize int64                     `json:"pageSize"`
}

func toOrgsListPaginated(orgs []OrganizationResponseDTO, total, page, pageSize int64) *OrgsListPaginated {
	return &OrgsListPaginated{
		Orgs:     orgs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
