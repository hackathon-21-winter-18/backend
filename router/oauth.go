package router

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/hackathon-21-winter-18/backend/service"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionCodeVerifierKey = "code_verifier"
	// sessionUserKey         = "user"
	redirect_url        = "http://localhost:8080/api/oauth/callback"
	authEndPoint        = "https://accounts.google.com/o/oauth2/v2/auth?"
	codeChallengeMethod = "S256"
)

// 旧ログインシステム
type LoginRequestBody struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}

type Me struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	UnreadNotices int    `json:"unreadNotices"`
}

func postSignUp(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	if req.Password == "" || req.Name == "" {
		return c.String(http.StatusBadRequest, "invalid request")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	userID, err := model.PostSignUp(c, req.Name, hashedPass)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	res := LoginResponse{
		ID:   *userID,
		Name: req.Name,
	}

	return echo.NewHTTPError(http.StatusOK, res)
}

func postLogin(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	if req.Password == "" || req.Name == "" {
		return c.String(http.StatusBadRequest, "invalid request")
	}

	userID, err := model.PostLogin(c, req.Name, req.Password)
	if err != nil {
		c.Logger().Error(err)
		return generateEchoError(err)
	}

	res := LoginResponse{
		ID:   *userID,
		Name: req.Name,
	}

	return echo.NewHTTPError(http.StatusOK, res)
}

func postLogout(c echo.Context) error {
	err := s.RevokeSession(c)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return echo.NewHTTPError(http.StatusOK)
}

func getWhoamI(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return errSessionNotFound(err)
	}

	userID := sess.Values["userID"].(string)
	ctx := c.Request().Context()
	name, err := model.GetMe(ctx, userID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	userIDinUUID, err := uuid.Parse(userID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	unreadNotices, err := model.GetCountOfUnreadNotices(ctx, userIDinUUID)
	if err != nil {
		c.Logger().Error(err)
		return generateEchoError(err)
	}

	return echo.NewHTTPError(http.StatusOK, Me{
		ID:            userID,
		Name:          name,
		UnreadNotices: unreadNotices,
	})
}

func generatePKCE(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sess.Options.SameSite = http.SameSiteNoneMode
	sess.Options.Secure = true

	codeVerifier := randAlphabetAndNumberString(43)
	sess.Values[sessionCodeVerifierKey] = codeVerifier

	codeVerifierHash := sha256.Sum256([]byte(codeVerifier))
	encoder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_").WithPadding(base64.NoPadding)

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	values := url.Values{}
	values.Add("response_type", "code")
	values.Add("client_id", service.PalamoClientID)
	values.Add("scope", service.OauthScope)
	values.Add("redirect_uri", redirect_url)
	values.Add("code_challenge_method", codeChallengeMethod)
	values.Add("code_challenge", encoder.EncodeToString(codeVerifierHash[:]))

	return c.Redirect(http.StatusFound, authEndPoint+values.Encode())
}

func authCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if len(code) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "code is required")
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	codeVerifier, ok := sess.Values[sessionCodeVerifierKey].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	res, err := service.RequestAccessToken(code, codeVerifier)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	email, err := service.FetchGoogleEmailAddress(res.AccessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}



	return echo.NewHTTPError(http.StatusOK)
}

var randSrcPool = sync.Pool{
	New: func() interface{} {
		return rand.NewSource(time.Now().UnixNano())
	},
}

const (
	rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

func randAlphabetAndNumberString(n int) string {
	b := make([]byte, n)
	randSrc := randSrcPool.Get().(rand.Source)
	cache, remain := randSrc.Int63(), rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), rs6LetterIdxMax
		}
		idx := int(cache & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i--
		}
		cache >>= rs6LetterIdxBits
		remain--
	}
	randSrcPool.Put(randSrc)
	return string(b)
}
