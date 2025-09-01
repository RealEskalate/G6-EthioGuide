package domain

type SearchResult struct {
	Procedures    []Procedure
	Organizations []AccountOrgSearch
}

type AccountOrgSearch struct {
	ID            string
	Name          string
	Email         string
	ProfilePicURL string
	Role          Role

	OrganizationDetail *OrganizationDetail
}

func ToSearch(orgs *Account) *AccountOrgSearch {
	return &AccountOrgSearch{
		ID:                 orgs.ID,
		Name:               orgs.Name,
		Email:              orgs.Email,
		ProfilePicURL:      orgs.ProfilePicURL,
		Role:               orgs.Role,
		OrganizationDetail: orgs.OrganizationDetail,
	}
}

type SearchFilterRequest struct {
	Query string 
	Page  string 
	Limit string
}