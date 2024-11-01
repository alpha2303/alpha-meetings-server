package app

import (
	"encoding/gob"

	"github.com/alpha2303/alpha-meetings/internal/app/auth"
	"github.com/alpha2303/alpha-meetings/internal/app/routes"
	"github.com/alpha2303/alpha-meetings/internal/pkg/helpers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func StartServer(addr string) error {
	serverApp := gin.Default()
	authenticator, err := auth.New()
	if err != nil {
		return err
	}

	gob.Register(map[string]any{})
	store := cookie.NewStore([]byte("secret"))
	serverApp.Use(sessions.Sessions("auth-session", store))

	rootRoute := serverApp.Group("/api")
	rootRoute.GET("/", func(ctx *gin.Context) {
		helpers.SendResponse(ctx, 200, "Health Check", map[string]string{"body": "Hello World!"})
	})

	routes.AddRoutes(authenticator, rootRoute)

	return serverApp.Run(addr)
}
