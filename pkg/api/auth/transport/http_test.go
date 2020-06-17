package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/auth"
	"github.com/jpurdie/authapi/pkg/api/auth/transport"
	"github.com/jpurdie/authapi/pkg/utl/jwt"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/jpurdie/authapi/pkg/utl/mock"
	"github.com/jpurdie/authapi/pkg/utl/mock/mockdb"
	"github.com/jpurdie/authapi/pkg/utl/server"

	"github.com/go-pg/pg/v9/orm"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *authapi.AuthToken
		udb        *mockdb.User
		jwt        *mock.JWT
		sec        *mock.Secure
	}{
		{
			name:       "Invalid request",
			req:        `{"username":"juzernejm"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fail on FindByUsername",
			req:        `{"username":"juzernejm","password":"hunter123"}`,
			wantStatus: http.StatusInternalServerError,
			udb: &mockdb.User{
				FindByUsernameFn: func(orm.DB, string) (authapi.User, error) {
					return authapi.User{}, authapi.ErrGeneric
				},
			},
		},
		{
			name:       "Success",
			req:        `{"username":"juzernejm","password":"hunter123"}`,
			wantStatus: http.StatusOK,
			udb: &mockdb.User{
				FindByUsernameFn: func(orm.DB, string) (authapi.User, error) {
					return authapi.User{
						Password: "hunter123",
						Active:   true,
					}, nil
				},
				UpdateFn: func(db orm.DB, u authapi.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(authapi.User) (string, error) {
					return "jwttokenstring", nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				TokenFn: func(string) string {
					return "refreshtoken"
				},
			},
			wantResp: &authapi.AuthToken{Token: "jwttokenstring", RefreshToken: "refreshtoken"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, tt.jwt, tt.sec, nil), r, nil)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/login"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(authapi.AuthToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				tt.wantResp.RefreshToken = response.RefreshToken
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestRefresh(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *authapi.RefreshToken
		udb        *mockdb.User
		jwt        *mock.JWT
	}{
		{
			name:       "Fail on FindByToken",
			req:        "refreshtoken",
			wantStatus: http.StatusInternalServerError,
			udb: &mockdb.User{
				FindByTokenFn: func(orm.DB, string) (authapi.User, error) {
					return authapi.User{}, authapi.ErrGeneric
				},
			},
		},
		{
			name:       "Success",
			req:        "refreshtoken",
			wantStatus: http.StatusOK,
			udb: &mockdb.User{
				FindByTokenFn: func(orm.DB, string) (authapi.User, error) {
					return authapi.User{
						Username: "johndoe",
						Active:   true,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(authapi.User) (string, error) {
					return "jwttokenstring", nil
				},
			},
			wantResp: &authapi.RefreshToken{Token: "jwttokenstring"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, tt.jwt, nil, nil), r, nil)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/refresh/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(authapi.RefreshToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestMe(t *testing.T) {
	cases := []struct {
		name       string
		wantStatus int
		wantResp   authapi.User
		header     string
		udb        *mockdb.User
		rbac       *mock.RBAC
	}{
		{
			name:       "Fail on user view",
			wantStatus: http.StatusInternalServerError,
			udb: &mockdb.User{
				ViewFn: func(orm.DB, int) (authapi.User, error) {
					return authapi.User{}, authapi.ErrGeneric
				},
			},
			rbac: &mock.RBAC{
				UserFn: func(echo.Context) authapi.AuthUser {
					return authapi.AuthUser{ID: 1}
				},
			},
			header: mock.HeaderValid(),
		},
		{
			name:       "Success",
			wantStatus: http.StatusOK,
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, i int) (authapi.User, error) {
					return authapi.User{
						Base: authapi.Base{
							ID: i,
						},
						CompanyID:  2,
						LocationID: 3,
						Email:      "john@mail.com",
						FirstName:  "John",
						LastName:   "Doe",
					}, nil
				},
			},
			rbac: &mock.RBAC{
				UserFn: func(echo.Context) authapi.AuthUser {
					return authapi.AuthUser{ID: 1}
				},
			},
			header: mock.HeaderValid(),
			wantResp: authapi.User{
				Base: authapi.Base{
					ID: 1,
				},
				CompanyID:  2,
				LocationID: 3,
				Email:      "john@mail.com",
				FirstName:  "John",
				LastName:   "Doe",
			},
		},
	}

	client := &http.Client{}
	jwt, err := jwt.New("HS256", "jwtsecret123", 60, 4)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, nil, nil, tt.rbac), r, authMw.Middleware(jwt))
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/me"
			req, err := http.NewRequest("GET", path, nil)
			req.Header.Set("Authorization", tt.header)
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp.ID != 0 {
				var response authapi.User
				if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
