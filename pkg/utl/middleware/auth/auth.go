package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jpurdie/authapi"
	jwtUtil "github.com/jpurdie/authapi/pkg/utl/jwt"
	"github.com/jpurdie/authapi/pkg/utl/postgres"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strings"
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

func CheckAuthorization(requiredRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			op := "CheckAuthorization"

			//checking org ID is valid UUID
			orgIdReq := c.QueryParam("org_id")
			orgUUID, err := uuid.Parse(orgIdReq)
			if err != nil {
				return c.JSON(http.StatusUnprocessableEntity, "")
			}
			//made it here. is valid UUID

			log.Println("Received request with " + orgUUID.String())

			db, err := postgres.DBConn()
			log.Println(db.PoolStats())
			defer db.Close()

			if err != nil {
				log.Panicln(&authapi.Error{
					Op:   op,
					Code: authapi.EINTERNAL,
					Err:  err,
				})
				//return c.JSON(http.StatusInternalServerError, "")
			}

			roleName, orgID, userID := "", "", ""

			err = db.Model((*authapi.Role)(nil)).
				Column("role.name", "o.id", "u.id").
				Join("JOIN organization_users AS ou ON ou.role_id = role.id").
				Join("JOIN organizations AS o ON ou.role_id = role.id").
				Join("JOIN users AS u ON u.id = ou.user_id").
				Where("o.uuid = ?", orgUUID.String()).
				Where("u.external_id = ?", c.Get("sub").(string)).
				Select(&roleName, &orgID, &userID)

			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusUnauthorized, "")
			}

			for _, role := range requiredRoles {
				if strings.ToLower(role) == strings.ToLower(roleName) {

					c.Set("orgID", orgID)
					c.Set("roleName", roleName)
					c.Set("userID", userID)
					return next(c)
				}
			}
			return c.JSON(http.StatusUnauthorized, "")

		}
	}
}
