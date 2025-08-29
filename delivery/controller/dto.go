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
	User         domain.Account `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token,omitempty"`
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
