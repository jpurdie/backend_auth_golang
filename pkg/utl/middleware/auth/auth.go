package auth

import (
	jwtUtil "github.com/jpurdie/authapi/pkg/utl/jwt"
	//"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// TokenParser represents JWT token parser
type TokenParser interface {
	ParseToken(string) (*jwt.Token, error)
}
type userInterface interface {
}

// Middleware makes JWT implement the Middleware interface.
//func Middleware(tokenParser TokenParser) echo.MiddlewareFunc {
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			handler := jwtUtil.New()
			err := handler.CheckJWT(c.Response(), c.Request())
			if err != nil {
				panic(err.Error())
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
