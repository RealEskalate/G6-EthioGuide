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

type TranslateDTO struct {
	Content string `json:"content" binding:"required"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token,omitempty"`
}

type ForgotDTO struct {
	Email string `json:"email"`
}

type ResetDTO struct {
	ResetToken  string `json:"resetToken"`
	NewPassword string `json:"new_password"`
}

type ActivateDTO struct {
	ActivateToken string `json:"activatationToken"`
}

type ProcedureResponse struct {
	ID             int                    `json:"id"`
	GroupID        *int                   `json:"group_id,omitempty"`
	OrganizationID int                    `json:"organization_id" `
	Name           string                 `json:"name"`
	Content        map[string]interface{} `json:"content,omitempty"`
	Fees           map[string]interface{} `json:"fees,omitempty"`
	ProcessingTime map[string]interface{} `json:"processing_time,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
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
