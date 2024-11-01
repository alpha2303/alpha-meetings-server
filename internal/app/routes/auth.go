package routes

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"

	"github.com/alpha2303/alpha-meetings/internal/app/auth"
	"github.com/alpha2303/alpha-meetings/internal/pkg/exceptions"
	"github.com/alpha2303/alpha-meetings/internal/pkg/helpers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRoutesInjector(auth *auth.Authenticator, rootGroup *gin.RouterGroup) {
	authGroup := rootGroup.Group("/auth")
	authGroup.GET("/login", loginHandler(auth))
	authGroup.GET("/callback", authCallbackHandler(auth))
	authGroup.GET("/logout", logoutHandler)
}

func loginHandler(auth *auth.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state, err := generateRandomState()

		if err != nil {
			helpers.SendResponse(
				ctx,
				http.StatusInternalServerError,
				err.Error(),
				nil,
			)
		}

		session := sessions.Default(ctx)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			helpers.SendResponse(
				ctx,
				http.StatusInternalServerError,
				err.Error(),
				nil,
			)
		}

		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
}

func generateRandomState() (string, error) {
	stateBuffer := make([]byte, 32)
	_, err := rand.Read(stateBuffer)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(stateBuffer)

	return state, nil
}

func authCallbackHandler(auth *auth.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			helpers.SendResponse(
				ctx,
				http.StatusBadRequest,
				exceptions.InvalidSessionStateException.Error(),
				nil,
			)
		}

		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			helpers.SendResponse(
				ctx,
				http.StatusUnauthorized,
				exceptions.AuthExchangeFailureException.Error(),
				nil,
			)
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			helpers.SendResponse(
				ctx,
				http.StatusInternalServerError,
				exceptions.IDTokenVerificationException.Error(),
				nil,
			)
		}

		var profile map[string]any
		if err := idToken.Claims(&profile); err != nil {
			helpers.SendResponse(
				ctx,
				http.StatusInternalServerError,
				err.Error(),
				nil,
			)
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)

		if err := session.Save(); err != nil {
			helpers.SendResponse(
				ctx,
				http.StatusInternalServerError,
				err.Error(),
				nil,
			)
		}

		// Check again later if redirect is required here
		// ctx.Redirect(http.StatusTemporaryRedirect, "/api/user/")
		helpers.SendResponse(
			ctx,
			http.StatusOK,
			"Logged in Successfully",
			profile,
		)
	}
}

func logoutHandler(ctx *gin.Context) {
	parsedLogoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		helpers.SendResponse(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
			nil,
		)
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host + "/api/")
	if err != nil {
		helpers.SendResponse(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
			nil,
		)
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	parsedLogoutUrl.RawQuery = parameters.Encode()

	ctx.SetCookie("auth-session", "", 0, "/", "localhost", false, false)

	ctx.Redirect(http.StatusTemporaryRedirect, parsedLogoutUrl.String())
}
