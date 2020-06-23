package authapi

type Company struct {
	Base
	Name        string        `json:"name"`
	Active      bool          `json:"active"`
	CompanyUser []CompanyUser `json:"-", pg:",many2many:company_users"`
}
