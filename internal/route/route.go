package route

import (
	"github.com/Vitaly-Baidin/auth-api/internal/controller"
	"github.com/Vitaly-Baidin/auth-api/internal/middleware"
	"github.com/Vitaly-Baidin/auth-api/internal/middleware/validators"
	"github.com/Vitaly-Baidin/auth-api/internal/service"
	"net/http"
)

const (
	basePath = "/v1"

	authPath = basePath + "/auth"

	authRegisterPath = authPath + "/register"
	authLoginPath    = authPath + "/login"
	authRefreshPost  = authPath + "/refresh"

	pingPath = basePath + "/ping"
)

func InitAuthRoutes(mux *http.ServeMux, ac controller.AuthController) {
	mux.Handle(authRegisterPath, validators.RegisterValidator(http.HandlerFunc(ac.Register)))
	mux.Handle(authLoginPath, validators.RegisterValidator(http.HandlerFunc(ac.Login)))
	mux.Handle(authRefreshPost, validators.RegisterValidator(http.HandlerFunc(ac.Refresh)))
}

func InitTestRoutes(mux *http.ServeMux, ts service.TokenService) {
	mux.Handle(pingPath, middleware.JWTMiddleware(http.HandlerFunc(controller.Ping), ts))
}
