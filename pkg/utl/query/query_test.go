package query_test

import (
	"testing"

	"github.com/labstack/echo"

	"github.com/jpurdie/authapi"

	"github.com/stretchr/testify/assert"

	"github.com/jpurdie/authapi/pkg/utl/query"
)

func TestList(t *testing.T) {
	type args struct {
		user authapi.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *authapi.ListQuery
		wantErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: authapi.AuthUser{
				Role: authapi.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: authapi.AuthUser{
				Role:      authapi.CompanyAdminRole,
				CompanyID: 1,
			}},
			wantData: &authapi.ListQuery{
				Query: "company_id = ?",
				ID:    1},
		},
		{
			name: "Location admin user",
			args: args{user: authapi.AuthUser{
				Role:       authapi.LocationAdminRole,
				CompanyID:  1,
				LocationID: 2,
			}},
			wantData: &authapi.ListQuery{
				Query: "location_id = ?",
				ID:    2},
		},
		{
			name: "Normal user",
			args: args{user: authapi.AuthUser{
				Role: authapi.UserRole,
			}},
			wantErr: echo.ErrForbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.List(tt.args.user)
			assert.Equal(t, tt.wantData, q)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
