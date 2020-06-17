package Auth0

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/encoding/json"
	"net/http"
	"os"
	"strings"
)

var ctx = context.Background()

func buildRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PW"),
		DB:       0, // use default DB
	})
}

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
}

func FetchAccessToken() (string, error) {
	rdb := buildRedisClient()
	accessToken, err := rdb.Get(ctx, "auth0_access_token").Result()
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(err)
	if accessToken != "" {
		fmt.Println("Access Token is present.")
		return accessToken, nil
	}
	fmt.Println("Access Token is not present. Going out to Auth0")

	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	url := "https://" + domain + "/oauth/token"
	audience := "https://" + domain + "/api/v2/"
	payload := strings.NewReader("{\"client_id\":\"" + clientId + "\",\"client_secret\": \"" + clientSecret + "\",\"audience\":\"" + audience + "\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	//	body, _ := ioutil.ReadAll(res.Body)

	//	fmt.Println(string(body))

	var atr accessTokenResp
	json.NewDecoder(res.Body).Decode(&atr)

	fmt.Println(atr)

	if res.Body != nil {

		err := rdb.Set(ctx, "auth0_access_token", atr.AccessToken, 30000).Err()
		if err != nil {
			return "", err
		}
	}
	return atr.AccessToken, nil

}

func CreateUser() error {
	accessToken, err := FetchAccessToken()
	if err != nil {
		return err
	}
	print(accessToken)
	return nil
}
