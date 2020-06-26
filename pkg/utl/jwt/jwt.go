package jwt

import (
	"encoding/json"
	"errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"

	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jpurdie/authapi"

	"github.com/dgrijalva/jwt-go"
)

var minSecretLen = 128

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Secret key used for signing.
	key []byte

	// Duration for which the jwt token is valid.
	ttl time.Duration

	// JWT signing algorithm
	algo jwt.SigningMethod
}

// ParseToken parses token from Authorization header
//func (s Service) ParseToken(authHeader string) (*jwt.Token, error) {
//	parts := strings.SplitN(authHeader, " ", 2)
//	if !(len(parts) == 2 && parts[0] == "Bearer") {
//		return nil, authapi.ErrGeneric
//	}
//
//	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
//		if s.algo != token.Method {
//			return nil, authapi.ErrGeneric
//		}
//		return s.key, nil
//	})
//
//}

// GenerateToken generates new JWT token and populates it with user data
func (s Service) GenerateToken(u authapi.User) (string, error) {
	return jwt.NewWithClaims(s.algo, jwt.MapClaims{
		"id": u.Base.ID,
		"e":  u.Email,
		//	"r":   u.Role.AccessLevel,
		"c": u.OrganizationID,
		//"l":   u.LocationID,
		"exp": time.Now().Add(s.ttl).Unix(),
	}).SignedString(s.key)

}

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func New() *jwtmiddleware.JWTMiddleware {

	myFunc := func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		aud := os.Getenv("AUTH0_AUDIENCE")
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAud {
			return token, errors.New("Invalid audience.")
		}
		// Verify 'iss' claim
		iss := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}

		cert, err := getPemCert(token)
		if err != nil {
			panic(err.Error())
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	}
	options := jwtmiddleware.Options{
		ValidationKeyGetter: myFunc,
		SigningMethod:       jwt.SigningMethodRS256,
	}

	return jwtmiddleware.New(options)

}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func checkScope(scope string, tokenString string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims, ok := token.Claims.(*CustomClaims)

	hasScope := false
	if ok && token.Valid {
		result := strings.Split(claims.Scope, " ")
		for i := range result {
			if result[i] == scope {
				hasScope = true
			}
		}
	}

	return hasScope
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func responseJSON(message string, w http.ResponseWriter, statusCode int) {
	response := Response{message}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
