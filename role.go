package authapi

// AccessRole represents access role type
type AccessRole int

const (
	// SuperAdminRole has all permissions
	OwnerRole AccessRole = 500

	// AdminRole has admin specific permissions
	SuperUserRole AccessRole = 400

	// CompanyAdminRole can edit company specific things
	AdminRole AccessRole = 300

	// LocationAdminRole can edit location specific things
	SupervisorRole AccessRole = 200

	// UserRole is a standard user
	UserRole AccessRole = 100
)

// Role model
type Role struct {
	ID          AccessRole `json:"-" db:"id"`
	AccessLevel AccessRole `json:"-" db:"access_level"`
	Name        string     `json:"name"  db:"name"`
	Active      bool       `json:"-"  db:"active"`
}
