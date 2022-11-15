package main

import (
	"fmt"
	"github.com/Vitaly-Baidin/auth-api/config"
	"github.com/Vitaly-Baidin/auth-api/internal/controller"
	"github.com/Vitaly-Baidin/auth-api/internal/repository"
	"github.com/Vitaly-Baidin/auth-api/internal/route"
	"github.com/Vitaly-Baidin/auth-api/internal/service"
	"github.com/Vitaly-Baidin/auth-api/pkg/httpserver"
	"github.com/Vitaly-Baidin/auth-api/pkg/postgres"
	"github.com/Vitaly-Baidin/auth-api/pkg/redis"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	fmt.Println(cfg.PG.URL)
	fmt.Println(cfg.Redis.Addresses)

	// PostgreSQL + Redis
	pg, err := postgres.New(
		cfg.PG.URL,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.ConnAttempts(cfg.PG.ConnAttempts),
		postgres.ConnTimeout(time.Duration(cfg.PG.ConnTimeoutInSec)*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	r, err := redis.New(redis.Addresses(cfg.Redis.Addresses))
	if err != nil {
		log.Fatal(err)
	}

	// Migration

	migration(cfg.PG.URL)

	// Repositories
	var ur repository.UserRepository = repository.NewUserRepoPG(pg)
	var tr repository.TokenRepository = repository.NewTokenRepoRedis(r)

	// Services
	var us service.UserService = service.NewUserService(ur)
	var ts service.TokenService = service.NewTokenService(cfg, tr)

	// Controllers
	var auth controller.AuthController = controller.NewAuthController(us, ts)

	mux := http.NewServeMux()

	route.InitAuthRoutes(mux, auth)
	route.InitTestRoutes(mux, ts)

	httpServer := httpserver.New(mux,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(time.Duration(cfg.HTTP.ReadTimeoutInSec)*time.Second),
		httpserver.WriteTimeout(time.Duration(cfg.HTTP.WriteTimeoutInSec)*time.Second),
		httpserver.ShutdownTimeout(time.Duration(cfg.HTTP.ShutdownTimeoutInSec)*time.Second),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdowns
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = r.Close()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - redis.Close: %w", err))
	}

	pg.Close()
}
