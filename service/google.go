package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	tokenEndPoint = "https://www.googleapis.com/oauth2/v4/token"
	OauthScope    = "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"
	googleAPI     = "https://people.googleapis.com/v1/people/me?personFields=emailAddresses"
	Redirect_uri  = "http://localhost:8080/api/oauth/callback"
)

var (
	PalamoClientID     = os.Getenv("PALAMO_CLIENT_ID")
	palamoClientSecret = os.Getenv("PALAMO_CLIENT_SECRET")
)

type Authority struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type GoogleUserSource struct {
	ResourceName   string         `json:"resourceName"`
	EmailAddresses []EmailAddress `json:"emailAddresses"`
}

type EmailAddress struct {
	Value string `json:"value"`
}

type GoogleUser struct {
	ID           string
	EmailAddress string
}

func RequestAccessToken(code, codeVerifier string) (Authority, error) {
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("client_id", PalamoClientID)
	values.Set("client_secret", palamoClientSecret)
	values.Set("code", code)
	values.Set("redirect_uri", Redirect_uri)
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

func FetchGoogleUser(token string) (*GoogleUser, error) {
	req, err := http.NewRequest("GET", googleAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch email address")
	}

	var userSource GoogleUserSource
	if err := json.NewDecoder(res.Body).Decode(&userSource); err != nil {
		return nil, err
	}

	userID := userSource.ResourceName[7:28]
	if len(userID) != 21 {
		return nil, errors.New("failed to get valid userID")
	}
	if len(userSource.EmailAddresses) == 0 || userSource.EmailAddresses[0].Value == "" {
		return nil, errors.New("failed to get email address")
	}

	user := GoogleUser{
		ID:           userID,
		EmailAddress: userSource.EmailAddresses[0].Value,
	}

	return &user, nil
}
