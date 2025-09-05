package domain

import (
	"context"
)

type IAccountRepository interface {
	Create(ctx context.Context, user *Account) error
	GetById(ctx context.Context, id string) (*Account, error)
	GetByEmail(ctx context.Context, email string) (*Account, error)
	GetByUsername(ctx context.Context, username string) (*Account, error)
	GetOrgs(ctx context.Context, filter GetOrgsFilter) ([]*Account, int64, error)
	// GetByPhoneNumber(ctx context.Context, phone string) (*Account, error)
	UpdatePassword(ctx context.Context, accountID, newPassword string) error
	UpdateProfile(ctx context.Context, account Account) error

	// -----
	ExistsByEmail(ctx context.Context, email, excludeID string) (bool, error)
	ExistsByUsername(ctx context.Context, username, excludeID string) (bool, error)
	UpdateUserFields(ctx context.Context, userIDstr string, update map[string]interface{}) error
}

type ITokenRepository interface {
	CreateToken(ctx context.Context, token *Token) (*Token, error)
	GetToken(ctx context.Context, tokentype, token string) (string, error)
	DeleteToken(ctx context.Context, tokentype, token string) error
}

type ICategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetCategories(ctx context.Context, options *CategorySearchAndFilter) ([]*Category, int64, error)
}

type IProcedureRepository interface {
	Create(ctx context.Context, procedure *Procedure) error
	GetByID(ctx context.Context, id string) (*Procedure, error)
	Update(ctx context.Context, id string, procedure *Procedure) error
	Delete(ctx context.Context, id string) error
	SearchAndFilter(ctx context.Context, opttions ProcedureSearchFilterOptions) ([]*Procedure, int64, error)
	SearchByEmbedding(ctx context.Context, queryVec []float64, limit int) ([]*Procedure, error)
}

type IFeedbackRepository interface {
	SubmitFeedback(ctx context.Context, feedback *Feedback) error
	GetFeedbackByID(ctx context.Context, id string) (*Feedback, error)
	GetAllFeedbacksForProcedure(ctx context.Context, procedureID string, filter *FeedbackFilter) ([]*Feedback, int64, error)
	UpdateFeedbackStatus(ctx context.Context, feedbackID string, newFeedback *Feedback) error
	GetAllFeedbacks(ctx context.Context, filter *FeedbackFilter) ([]*Feedback, int64, error)
}

type IPostRepository interface {
	CreatePost(ctx context.Context, Post *Post) (*Post, error)
	GetPosts(ctx context.Context, opts PostFilters) ([]*Post, int64, error)
	GetPostByID(ctx context.Context, id string) (*Post, error)
	UpdatePost(ctx context.Context, Post *Post) (*Post, error)
	DeletePost(ctx context.Context, id, userID, role string) error
}
type IPreferencesRepository interface {
	Create(ctx context.Context, preferences *Preferences) error
	GetByUserID(ctx context.Context, userID string) (*Preferences, error)
	UpdateByUserID(ctx context.Context, userID string, preferences *Preferences) error
}

type INoticeRepository interface {
	Create(ctx context.Context, notice *Notice) error
	GetByFilter(ctx context.Context, filter *NoticeFilter) ([]*Notice, error)
	CountByFilter(ctx context.Context, filter *NoticeFilter) (int64, error)
	Update(ctx context.Context, id string, notice *Notice) error
	Delete(ctx context.Context, id string) error
}

type ISearchRepository interface {
	Search(ctx context.Context, filter SearchFilterRequest) (*SearchResult, error)
	FindProcedures(ctx context.Context, filter SearchFilterRequest) ([]*Procedure, error)
	FindOrganizations(ctx context.Context, filter SearchFilterRequest) ([]*AccountOrgSearch, error)
}

type IChecklistRepository interface {
	CreateChecklist(ctx context.Context, userid, procdureID string) (*UserProcedure, error)
	GetProcedures(ctx context.Context, userid string) ([]*UserProcedure, error)
	GetChecklistByUserProcedureID(ctx context.Context, userprocedureID string) ([]*Checklist, error)
	ToggleCheck(ctx context.Context, checklistID string) error
	FindCheck(ctx context.Context, checklistID string) (*Checklist, error)
	CountDocumentsChecklist(ctx context.Context, filter interface{}) (int64, error)
	UpdateUserProcedure(ctx context.Context, filter interface{}, update map[string]interface{}) error
}

type IAIChatRepository interface {
	Save(ctx context.Context, chat *AIChat) error
	GetByUser(ctx context.Context, userID string, limit int) ([]*AIChat, error)
	DeleteByUser(ctx context.Context, userID string) error
}
