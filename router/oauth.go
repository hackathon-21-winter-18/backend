package router

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/hackathon-21-winter-18/backend/service"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	sessionCodeVerifierKey = "code_verifier"
	sessionUserKey         = "user"
	authEndPoint           = "https://accounts.google.com/o/oauth2/v2/auth?"
	codeChallengeMethod    = "S256"
)

// 旧ログインシステム
// type LoginRequestBody struct {
// 	Name     string `json:"name,omitempty"`
// 	Password string `json:"password,omitempty"`
// }

// type LoginResponse struct {
// 	ID   uuid.UUID `json:"id,omitempty"`
// 	Name string    `json:"name,omitempty"`
// }

type Me struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	UnreadNotices int    `json:"unreadNotices"`
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
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sess.Options.SameSite = http.SameSiteNoneMode
	sess.Options.Secure = true

	codeVerifier := model.RandAlphabetAndNumberString(43)
	// log.Print(codeVerifier + "aaaaaaaaaaaaaaaa")
	sess.Values[sessionCodeVerifierKey] = codeVerifier

	codeVerifierHash := sha256.Sum256([]byte(codeVerifier))
	encoder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_").WithPadding(base64.NoPadding)

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	values := url.Values{}
	values.Add("response_type", "code")
	values.Add("client_id", service.PalamoClientID)
	values.Add("scope", service.OauthScope)
	values.Add("redirect_uri", service.Redirect_uri)
	values.Add("access_type", "offline")
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
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	codeVerifier, ok := sess.Values[sessionCodeVerifierKey].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get code_verifier")
	}

	res, err := service.RequestAccessToken(code, codeVerifier)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// log.Print(strconv.Itoa(res.ExpiresIn) + "fffffffffffffffffffffffffffffffffffffffffffffffffff")
	// log.Print(res.RefreshToken + "dddddddddddddddddddddddddddd")
	googleUser, err := service.FetchGoogleUser(res.AccessToken)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx := c.Request().Context()
	userID, err := model.GetUserIDByGoogleID(ctx, googleUser.ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if userID == nil {
		userID, err = model.CreateUser(ctx, googleUser.ID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	sess.Options.SameSite = http.SameSiteNoneMode
	sess.Options.Secure = true
	sess.Values["userID"] = &userID
	sess.Values["email"] = googleUser.EmailAddress
	sess.Save(c.Request(), c.Response())

	return echo.NewHTTPError(http.StatusOK, googleUser)
	// return c.Redirect(http.StatusSeeOther, "/")
}

// func postSignUp(c echo.Context) error {
// 	var req LoginRequestBody
// 	c.Bind(&req)

// 	if req.Password == "" || req.Name == "" {
// 		return c.String(http.StatusBadRequest, "invalid request")
// 	}

// 	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.Logger().Error(err)
// 		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
// 	}

// 	userID, err := model.PostSignUp(c, req.Name, hashedPass)
// 	if err != nil {
// 		c.Logger().Error(err)
// 		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
// 	}

// 	res := LoginResponse{
// 		ID:   *userID,
// 		Name: req.Name,
// 	}

// 	return echo.NewHTTPError(http.StatusOK, res)
// }

// func postLogin(c echo.Context) error {
// 	var req LoginRequestBody
// 	c.Bind(&req)

// 	if req.Password == "" || req.Name == "" {
// 		return c.String(http.StatusBadRequest, "invalid request")
// 	}

// 	userID, err := model.PostLogin(c, req.Name, req.Password)
// 	if err != nil {
// 		c.Logger().Error(err)
// 		return generateEchoError(err)
// 	}

// 	res := LoginResponse{
// 		ID:   *userID,
// 		Name: req.Name,
// 	}

// 	return echo.NewHTTPError(http.StatusOK, res)
// }