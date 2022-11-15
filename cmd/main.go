package main

import (
	"github.com/Vitaly-Baidin/auth-api/internal/controller"
	"github.com/Vitaly-Baidin/auth-api/internal/repository"
	"github.com/Vitaly-Baidin/auth-api/internal/route"
	"github.com/Vitaly-Baidin/auth-api/internal/service"
	"github.com/Vitaly-Baidin/auth-api/pkg/postgres"
	"github.com/Vitaly-Baidin/auth-api/pkg/redis"
	"log"
	"net/http"
)

func main() {
	// PostgreSQL + Redis
	pg, err := postgres.New("postgres://root:rootroot@localhost:5432/auth")
	if err != nil {
		log.Fatal(err)
	}

	r, err := redis.New("redis://localhost:6379/")
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	var ur repository.UserRepository = repository.NewUserRepoPG(pg)
	var tr repository.TokenRepository = repository.NewTokenRepoRedis(r)

	// Services
	var us service.UserService = service.NewUserService(ur)
	var ts service.TokenService = service.NewTokenService(tr)

	// Controllers
	var auth controller.AuthController = controller.NewAuthController(us, ts)

	mux := http.NewServeMux()

	route.InitAuthRoutes(mux, auth)
	route.InitTestRoutes(mux, ts)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
