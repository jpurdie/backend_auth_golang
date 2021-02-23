package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	jwtUtil "github.com/jpurdie/authapi/pkg/utl/jwt"
	"github.com/labstack/echo/v4"
)

// TokenParser represents JWT token parser
type TokenParser interface {
	ParseToken(string) (*jwt.Token, error)
}
type userInterface interface {
}

// Middleware makes JWT implement the Middleware interface.
//func Middleware(tokenParser TokenParser) echo.MiddlewareFunc {
func Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			handler := jwtUtil.New()
			err := handler.CheckJWT(c.Response(), c.Request())
			if err != nil {
				return c.String(http.StatusUnauthorized, "")
			}

			userContext := c.Request().Context().Value("user").(*jwt.Token)
			claims := userContext.Claims.(jwt.MapClaims)
			sub := claims["sub"].(string)
			iss := claims["iss"].(string)
			//	aud := claims["aud"].(string)
			iat := int(claims["iat"].(float64))
			exp := int(claims["exp"].(float64))
			azp := claims["azp"].(string)
			scope := claims["scope"].(string)

			c.Set("foo", "bar")
			c.Set("sub", sub)
			c.Set("iss", iss)
			//	c.Set("aud", aud)
			c.Set("iat", iat)
			c.Set("exp", exp)
			c.Set("azp", azp)
			c.Set("scope", scope)
			return next(c)
		}
	}
}

type ProfileStruct struct {
	RoleID   int    `db:"roleID"`
	RoleName string `db:"roleName"`
	OrgID    int    `db:"orgID"`
	UserID   int    `db:"userID"`
	ProfID   int    `db:"profileID"`
}

func CheckAuthorization(db sqlx.DB, requiredRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			op := "CheckAuthorization"
			//checking org ID is valid UUID
			orgIdReq := c.Request().Header.Get("Org-ID")
			orgUUID, err := uuid.Parse(orgIdReq)
			if err != nil {
				return c.NoContent(http.StatusUnprocessableEntity)
			}
			//made it here. is valid UUID
			log.Println("Received request with " + orgUUID.String())

			if err != nil {
				log.Panicln(&authapi.Error{
					Op:   op,
					Code: authapi.EINTERNAL,
					Err:  err,
				})
			}
			profStruct := ProfileStruct{}

			query := "SELECT " +
				"role.name as \"roleName\", " +
				"role.id as \"roleID\", " +
				"o.id as \"orgID\", " +
				"u.id as \"userID\", " +
				"p.id as \"profileID\" " +
				"FROM roles AS role " +
				"JOIN profiles AS p ON p.role_id = role.id " +
				"JOIN organizations AS o ON p.role_id = role.id " +
				"JOIN users AS u ON u.id = p.user_id " +
				"WHERE o.uuid = $1 " +
				"AND u.external_id = $2 " +
				"AND o.active = true " +
				"AND p.active = true;"

			row := db.QueryRowx(query, orgUUID.String(), c.Get("sub").(string)).StructScan(&profStruct)
			if row != nil {
				log.Println(err)
				return c.JSON(http.StatusUnauthorized, "")
			}

			for _, role := range requiredRoles {
				if strings.ToLower(role) == strings.ToLower(profStruct.RoleName) {
					c.Set("orgID", profStruct.OrgID)
					c.Set("roleName", profStruct.RoleName)
					c.Set("roleLevel", profStruct.RoleID)
					c.Set("userID", profStruct.UserID)
					c.Set("profileID", profStruct.ProfID)
					return next(c)
				}
			}
			return c.JSON(http.StatusUnauthorized, "")
		}
	}
}
