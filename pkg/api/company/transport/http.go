package transport

import (
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/company"
	AuthUtil "github.com/jpurdie/authapi/pkg/utl/Auth"
	"github.com/jpurdie/authapi/pkg/utl/Auth0"
	"github.com/labstack/echo"
	"net/http"
)

// HTTP represents company http service
type HTTP struct {
	svc company.Service
}

// NewHTTP creates new company http service
func NewHTTP(svc company.Service, r *echo.Group) {
	h := HTTP{svc}
	ur := r.Group("/companies")
	// swagger:route POST /v1/users users userCreate
	// Creates new user account.
	// responses:
	//  200: userResp
	//  400: errMsg
	//  401: err
	//  403: errMsg
	//  500: err
	ur.POST("", h.create)

	// swagger:operation GET /v1/users users listUsers
	// ---
	// summary: Returns list of users.
	// description: Returns list of users. Depending on the user role requesting it, it may return all users for SuperAdmin/Admin users, all company/location users for Company/Location admins, and an error for non-admin users.
	// parameters:
	// - name: limit
	//   in: query
	//   description: number of results
	//   type: int
	//   required: false
	// - name: page
	//   in: query
	//   description: page number
	//   type: int
	//   required: false
	// responses:
	//   "200":
	//     "$ref": "#/responses/userListResp"
	//   "400":
	//     "$ref": "#/responses/errMsg"
	//   "401":
	//     "$ref": "#/responses/err"
	//   "403":
	//     "$ref": "#/responses/err"
	//   "500":
	//     "$ref": "#/responses/err"
	//ur.GET("", h.list)
	//
	//// swagger:operation GET /v1/users/{id} users getUser
	//// ---
	//// summary: Returns a single user.
	//// description: Returns a single user by its ID.
	//// parameters:
	//// - name: id
	////   in: path
	////   description: id of user
	////   type: int
	////   required: true
	//// responses:
	////   "200":
	////     "$ref": "#/responses/userResp"
	////   "400":
	////     "$ref": "#/responses/err"
	////   "401":
	////     "$ref": "#/responses/err"
	////   "403":
	////     "$ref": "#/responses/err"
	////   "404":
	////     "$ref": "#/responses/err"
	////   "500":
	////     "$ref": "#/responses/err"
	//ur.GET("/:id", h.view)
	//
	//// swagger:operation PATCH /v1/users/{id} users userUpdate
	//// ---
	//// summary: Updates user's contact information
	//// description: Updates user's contact information -> first name, last name, mobile, phone, address.
	//// parameters:
	//// - name: id
	////   in: path
	////   description: id of user
	////   type: int
	////   required: true
	//// - name: request
	////   in: body
	////   description: Request body
	////   required: true
	////   schema:
	////     "$ref": "#/definitions/userUpdate"
	//// responses:
	////   "200":
	////     "$ref": "#/responses/userResp"
	////   "400":
	////     "$ref": "#/responses/errMsg"
	////   "401":
	////     "$ref": "#/responses/err"
	////   "403":
	////     "$ref": "#/responses/err"
	////   "500":
	////     "$ref": "#/responses/err"
	//ur.PATCH("/:id", h.update)
	//
	//// swagger:operation DELETE /v1/users/{id} users userDelete
	//// ---
	//// summary: Deletes a user
	//// description: Deletes a user with requested ID.
	//// parameters:
	//// - name: id
	////   in: path
	////   description: id of user
	////   type: int
	////   required: true
	//// responses:
	////   "200":
	////     "$ref": "#/responses/ok"
	////   "400":
	////     "$ref": "#/responses/err"
	////   "401":
	////     "$ref": "#/responses/err"
	////   "403":
	////     "$ref": "#/responses/err"
	////   "500":
	////     "$ref": "#/responses/err"
	//ur.DELETE("/:id", h.delete)
}

// Custom errors
var (
	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
	ErrPasswordNotValid    = echo.NewHTTPError(http.StatusBadRequest, "passwords are not in the valid format")
)

// User create request
// swagger:model userCreate
type createOrgUserReq struct {
	CompanyName     string `json:"orgName" validate:"required,min=4"`
	FirstName       string `json:"firstName" validate:"required,min=2"`
	LastName        string `json:"lastName" validate:"required,min=2"`
	Password        string `json:"password" validate:"required`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email"`
}

func (h HTTP) create(c echo.Context) error {
	r := new(createOrgUserReq)

	if err := c.Bind(r); err != nil {
		return err
	}

	if r.Password != r.PasswordConfirm {
		return ErrPasswordsNotMaching
	}

	if !AuthUtil.VerifyPassword(r.Password) {
		return ErrPasswordNotValid
	}

	company := authapi.Company{Name: r.CompanyName, Active: true}

	u := authapi.User{
		Password:   r.Password,
		Email:      r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		ExternalID: "",
		Active:     true,
	}
	x := uuid.New()
	cu := authapi.CompanyUser{Company: &company, User: &u, UUID: x}

	externalID, err := Auth0.CreateUser(u)
	if err != nil {
		return err
	}

	u.ExternalID = externalID

	companyUser, err := h.svc.Create(c, cu)
	if err != nil {
		//delete auth0 user
		return err
	}

	err = Auth0.SendVerificationEmail(u)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, companyUser)
}

//type listResponse struct {
//	Companies []authapi.Company `json:"companies"`
//	Page      int             `json:"page"`
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
