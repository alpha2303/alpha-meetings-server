package routes

import (
	"github.com/alpha2303/alpha-meetings/internal/app/auth"
	"github.com/gin-gonic/gin"
)

type RouteInjector func(*auth.Authenticator, *gin.RouterGroup)

var (
	routeInjectors []RouteInjector = []RouteInjector{
		AuthRoutesInjector,
	}
)

func AddRoutes(auth *auth.Authenticator, rootGroup *gin.RouterGroup) {
	for _, fn := range routeInjectors {
		fn(auth, rootGroup)
	}
}
