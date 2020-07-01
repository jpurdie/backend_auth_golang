package auth0

import (
	"bytes"
	"context"
	"errors"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/utl/redis"
	"github.com/segmentio/encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var ctx = context.Background()

var (
	ErrUnableToReachAuth0 = errors.New("unable to reach authentication service")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUnableToCreateUser = errors.New("unable to create user")
)

type accessTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	Expires      string `json:"expires_in"`
}

func FetchAccessToken() (string, error) {
	rdb := redis.BuildRedisClient()
	accessToken, _ := rdb.Get(ctx, "auth0_access_token").Result()

	if accessToken != "" {
		log.Println("Access Token is present.")
		return accessToken, nil
	}
	log.Println("Access Token is not present. Going out to Auth0")

	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	url := "https://" + domain + "/oauth/token"
	audience := "https://" + domain + "/api/v2/"
	payload := strings.NewReader("{\"client_id\":\"" + clientId + "\",\"client_secret\": \"" + clientSecret + "\",\"audience\":\"" + audience + "\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	log.Println("HTTP Response Status:", res.StatusCode, http.StatusText(res.StatusCode))
	if res.StatusCode != 201 && res.StatusCode != 200 {
		return "", errors.New("Unable to get access token")
	}
	defer res.Body.Close()

	var atr accessTokenResp
	json.NewDecoder(res.Body).Decode(&atr)

	log.Println("Access token length " + string(len(atr.AccessToken)))

	if res.Body != nil {

		//set the duration time to the expires in. The expires in integer from Auth0 is in seconds
		err := rdb.Set(ctx, "auth0_access_token", atr.AccessToken, time.Duration(30)*time.Second).Err()
		err = rdb.Set(ctx, "auth0_refresh_token", atr.RefreshToken, time.Duration(30)*time.Second).Err()
		err = rdb.Set(ctx, "auth0_id_token", atr.IDToken, time.Duration(30)*time.Second).Err()
		err = rdb.Set(ctx, "auth0_access_token_expires_in", atr.Expires, time.Duration(30)*time.Second).Err()
		if err != nil {
			return "", err
		}
	}
	return atr.AccessToken, nil

}

type appMetaData struct {
}
type createUserReq struct {
	Email         string      `json:"email"`
	Blocked       bool        `json:"blocked"`
	EmailVerified bool        `json:"email_verified"`
	AppMetaData   appMetaData `json:"app_metadata"`
	GivenName     string      `json:"given_name"`
	FamilyName    string      `json:"family_name"`
	Name          string      `json:"name"`
	Nickname      string      `json:"nickname"`
	Connection    string      `json:"connection"`
	Password      string      `json:"password"`
	VerifyEmail   bool        `json:"verify_email"`
}
type createUserResp struct {
	UserId string `json:"user_id"`
}

func CreateUser(u authapi.User) (string, error) {
	log.Println("Inside CreateUser()")
	const op = "Auth0.CreateUser"

	accessToken, err := FetchAccessToken()
	if err != nil {
		return "", ErrUnableToReachAuth0
	}
	a := appMetaData{}
	userReq := createUserReq{
		Email:         u.Email,
		Blocked:       false,
		EmailVerified: false,
		AppMetaData:   a,
		GivenName:     u.FirstName,
		FamilyName:    u.LastName,
		Name:          u.FirstName + " " + u.LastName,
		Nickname:      u.FirstName,
		Connection:    os.Getenv("AUTH0_CONNECTION"),
		Password:      u.Password,
		VerifyEmail:   false,
	}

	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	log.Println("Inside CreateUser()")

	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/api/v2/users"
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(userReq)
	log.Println("Inside CreateUser()")

	req, err := http.NewRequest("POST", url, b)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode == 409 {
		return "", &authapi.Error{
			Op:   op,
			Code: authapi.ECONFLICT,
			Err:  ErrUserAlreadyExists,
		}
	} else if res.StatusCode != 201 {
		return "", &authapi.Error{
			Op:   op,
			Code: authapi.EINTERNAL,
			Err:  ErrUnableToCreateUser,
		}
	}

	var cur createUserResp
	err = json.NewDecoder(res.Body).Decode(&cur)
	if err != nil {
		log.Fatal(err)
	}

	return cur.UserId, nil
}

type verEmailReq struct {
	ExternalID string `json:"user_id"`
	ClientID   string `json:"client_id"`
}
type verEmailResp struct {
	Status  string `json:"status"`
	Type    string `json:"type"`
	Created string `json:"created_at"`
	ID      string `json:"id"`
}

func SendVerificationEmail(u authapi.User) error {
	log.Println("Inside Send Verification")

	accessToken, err := FetchAccessToken()
	if err != nil {
		log.Println(err)

		return ErrUnableToReachAuth0
	}
	verEmailReq := verEmailReq{
		ExternalID: u.ExternalID,
		ClientID:   os.Getenv("AUTH0_CLIENT_ID"),
	}

	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/api/v2/jobs/verification-email"
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(verEmailReq)
	req, _ := http.NewRequest("POST", url, b)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	if res.StatusCode != 201 {
		log.Println(res.StatusCode)
		return errors.New("Unable to send verification email")
	}

	var vResp verEmailResp
	json.NewDecoder(res.Body).Decode(&vResp)
	log.Println("Inside Send Verification")

	return nil

}

func DeleteUser(u authapi.User) error {
	log.Println("Deleting Auth0 user " + u.ExternalID)
	accessToken, err := FetchAccessToken()
	if err != nil {
		return ErrUnableToReachAuth0
	}
	url := "https://" + os.Getenv("AUTH0_DOMAIN") + "/api/v2/users/" + u.ExternalID
	b := new(bytes.Buffer)
	req, _ := http.NewRequest("DELETE", url, b)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	if res.StatusCode != 204 {
		return errors.New("Unable to delete user from auth0 " + u.ExternalID)
	}

	var vResp verEmailResp
	json.NewDecoder(res.Body).Decode(&vResp)

	return nil

}
