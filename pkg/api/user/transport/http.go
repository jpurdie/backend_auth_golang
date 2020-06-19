package transport

import (
	"github.com/jpurdie/authapi/pkg/api/user"
	"net/http"

	"github.com/labstack/echo"
)

// HTTP represents user http service
type HTTP struct {
	svc user.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc user.Service, r *echo.Group) {
	h := HTTP{svc}
	ur := r.Group("/users")
	//ur.POST("", h.create)
	//ur.GET("", h.list)
	//ur.GET("/:id", h.view)
	//ur.PATCH("/:id", h.update)
	//ur.DELETE("/:id", h.delete)

	uc := ur.Group("/companies")
	uc.GET("/", h.listCompanies)

}
func (h HTTP) listCompanies(c echo.Context) error {
	return c.JSON(http.StatusOK, "list orgs")
}

// Custom errors
//var (
//	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
//)
//
//// User create request
//// swagger:model userCreate
//type createReq struct {
//	FirstName       string `json:"first_name" validate:"required"`
//	LastName        string `json:"last_name" validate:"required"`
//	Username        string `json:"username" validate:"required,min=3,alphanum"`
//	Password        string `json:"password" validate:"required,min=8"`
//	PasswordConfirm string `json:"password_confirm" validate:"required"`
//	Email           string `json:"email" validate:"required,email"`
//
//	CompanyID  int                `json:"company_id" validate:"required"`
//	LocationID int                `json:"location_id" validate:"required"`
//	RoleID     authapi.AccessRole `json:"role_id" validate:"required"`
//}
//

//func (h HTTP) create(c echo.Context) error {
//	r := new(createReq)
//
//	if err := c.Bind(r); err != nil {
//
//		return err
//	}
//
//	if r.Password != r.PasswordConfirm {
//		return ErrPasswordsNotMaching
//	}
//
//	if r.RoleID < authapi.SuperAdminRole || r.RoleID > authapi.UserRole {
//		return authapi.ErrBadRequest
//	}
//
//	usr, err := h.svc.Create(c, authapi.User{
//		Username:   r.Username,
//		Password:   r.Password,
//		Email:      r.Email,
//		FirstName:  r.FirstName,
//		LastName:   r.LastName,
//		CompanyID:  r.CompanyID,
//		LocationID: r.LocationID,
//		RoleID:     r.RoleID,
//	})
//
//	if err != nil {
//		return err
//	}
//
//	return c.JSON(http.StatusOK, usr)
//}
//
//type listResponse struct {
//	Users []authapi.User `json:"users"`
//	Page  int            `json:"page"`
//}
//
//func (h HTTP) list(c echo.Context) error {
//	var req authapi.PaginationReq
//	if err := c.Bind(&req); err != nil {
//		return err
//	}
//
//	result, err := h.svc.List(c, req.Transform())
//
//	if err != nil {
//		return err
//	}
//
//	return c.JSON(http.StatusOK, listResponse{result, req.Page})
//}
//
//func (h HTTP) view(c echo.Context) error {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		return authapi.ErrBadRequest
//	}
//
//	result, err := h.svc.View(c, id)
//	if err != nil {
//		return err
//	}
//
//	return c.JSON(http.StatusOK, result)
//}
//
//// User update request
//// swagger:model userUpdate
//type updateReq struct {
//	ID        int    `json:"-"`
//	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=2"`
//	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=2"`
//	Mobile    string `json:"mobile,omitempty"`
//	Phone     string `json:"phone,omitempty"`
//	Address   string `json:"address,omitempty"`
//}
//
//func (h HTTP) update(c echo.Context) error {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		return authapi.ErrBadRequest
//	}
//
//	req := new(updateReq)
//	if err := c.Bind(req); err != nil {
//		return err
//	}
//
//	usr, err := h.svc.Update(c, user.Update{
//		ID:        id,
//		FirstName: req.FirstName,
//		LastName:  req.LastName,
//		Mobile:    req.Mobile,
//		Phone:     req.Phone,
//		Address:   req.Address,
//	})
//
//	if err != nil {
//		return err
//	}
//
//	return c.JSON(http.StatusOK, usr)
//}
//
//func (h HTTP) delete(c echo.Context) error {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		return authapi.ErrBadRequest
//	}
//
//	if err := h.svc.Delete(c, id); err != nil {
//		return err
//	}
//
//	return c.NoContent(http.StatusOK)
//}
