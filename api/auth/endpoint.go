package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"

	"golang.org/x/oauth2"

	google "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type authHandlerT func(
	ctx context.Context,
	client *http.Client,
	jwtService JWTService,
	gqlClient *graphql.Client,
) (string, string, error)

func generateStateOauthCookie(c *gin.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, 3600, "", "", false, false)

	return state
}

func makeCallbackHandler(config oauth2.Config, handler func(context.Context, *http.Client) (string, string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		oauthState, err := c.Cookie("oauthstate")
		if err != nil {
			log.
				WithFields(log.Fields{
					"details": err.Error(),
				}).
				Error("oauthstate cookie is missing")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}

		redirectPath, err := c.Cookie("redirect_path")
		if err != nil {
			log.
				WithFields(log.Fields{
					"details": err.Error(),
				}).
				Error("redirect_path cookie is missing")
			c.Redirect(http.StatusTemporaryRedirect, redirectPath)
			c.Abort()
			return
		}

		if c.Query("state") != oauthState {
			log.Error("oauthstate cookie and returned state are different")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}

		token, err := config.Exchange(c, c.Query("code"))
		if err != nil {
			log.
				WithFields(log.Fields{
					"details": err.Error(),
				}).
				Error("Failed to get authorization token")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}

		_, refreshToken, err := handler(c, config.Client(c, token))
		if err != nil {
			log.
				WithFields(log.Fields{
					"details": err.Error(),
				}).
				Error("Failed to handle auth callback")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}

		c.SetCookie("refresh_token", refreshToken, Config.JWT.RefreshTTL, "/", strings.Split(Config.RedirectURL.Host, ":")[0], false, true)
		c.Redirect(http.StatusTemporaryRedirect, Config.RedirectURL.String()+redirectPath)
	}
}

func makeLoginHandler(config oauth2.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		redirectPath := c.Query("redirectPath")
		if redirectPath == "" {
			redirectPath = "/"
		}

		oauthState := generateStateOauthCookie(c)
		u := config.AuthCodeURL(oauthState)

		c.SetCookie("redirect_path", redirectPath, 3600, "/", "", false, false)

		c.Redirect(http.StatusTemporaryRedirect, u)
	}
}

func makeAuthHandler(
	jwtService JWTService,
	gqlClient *graphql.Client,
	handler authHandlerT,
) func(ctx context.Context, client *http.Client) (string, string, error) {
	return func(ctx context.Context, client *http.Client) (string, string, error) {
		return handler(ctx, client, jwtService, gqlClient)
	}
}

func makeRefreshHandler(jwtService JWTService) func(c *gin.Context) {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh_token is missing"})
			return
		}

		authToken, refreshToken, err := jwtService.RefreshToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.SetCookie("refresh_token", refreshToken, Config.JWT.RefreshTTL, "/", strings.Split(Config.RedirectURL.Host, ":")[0], false, true)
		c.JSON(http.StatusOK, gin.H{"auth_token": authToken})
	}
}

func handleGoogleAuth(
	ctx context.Context,
	client *http.Client,
	jwtService JWTService,
	gqlClient *graphql.Client,
) (string, string, error) {
	svc, err := google.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return "", "", err
	}

	user, err := svc.Userinfo.Get().Do()

	existingUser, err := GetUserByEmail(ctx, gqlClient, user.Email)
	if err != nil && !IsNotFoundError(err) {
		return "", "", err
	}

	if IsNotFoundError(err) {
		existingUser, err = InsertUser(
			ctx,
			gqlClient,
			&UserModel{
				Email: user.Email,
				Name:  user.Name,
				Role:  StudentRole.String(),
			})

		if err != nil {
			return "", "", err
		}
	}

	return jwtService.CreateToken(existingUser.ID, existingUser.Role)
}

func configureGoogleAuthEndpoint(
	group *gin.RouterGroup,
	client *graphql.Client,
	jwtService JWTService,
) {
	oauthConfig := oauth2.Config{
		ClientID:     Config.Auth.Google.ID,
		ClientSecret: Config.Auth.Google.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		RedirectURL: Config.Auth.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	callbackHandler := makeAuthHandler(jwtService, client, handleGoogleAuth)

	group.GET("/google/callback", makeCallbackHandler(oauthConfig, callbackHandler))
	group.GET("/google/login", makeLoginHandler(oauthConfig))
}

// ConfigureAuthEndpoint sets up auth endpoint
func ConfigureAuthEndpoint(
	group *gin.RouterGroup,
	jwtService JWTService,
	client *graphql.Client,
) {
	configureGoogleAuthEndpoint(
		group,
		client,
		jwtService,
	)

	group.GET("/refresh", makeRefreshHandler(jwtService))
}
