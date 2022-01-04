package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	tokenEndPoint       = "https://www.googleapis.com/oauth2/v4/token"
	OauthScope          = "https://www.googleapis.com/auth/userinfo.email"
)

var (
	PalamoClientID     = "868575110926-203ufh1rh90bv5vugit11buforu4q5p9.apps.googleusercontent.com"
	palamoClientSecret = os.Getenv("PALAMO_CLIENT_SECRET")
)

type Authority struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func RequestAccessToken(code, codeVerifier string) (Authority, error) {
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("client_id", PalamoClientID)
	values.Set("code", code)
	values.Set("code_verifier", codeVerifier)

	reqBody := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", tokenEndPoint, reqBody)
	if err != nil {
		return Authority{}, err
	}
	req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return Authority{}, err
	} else if res.StatusCode != http.StatusOK {
		return Authority{}, fmt.Errorf("failed to acquire access token")
	}

	var authRes Authority
	err = json.NewDecoder(res.Body).Decode(&authRes)
	if err != nil {
		return Authority{}, err
	}

	return authRes, nil
}

func FetchGoogleEmailAddress(token string) (string, error) {
	req, err := http.NewRequest("GET", OauthScope, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	} else if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch email address")
	}

} 
