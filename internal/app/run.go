package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AidlyTeam/Aidly-Backend/internal/config"
	"github.com/AidlyTeam/Aidly-Backend/internal/http"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/middlewares"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/server"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/AidlyTeam/Aidly-Backend/internal/services"
	validatorService "github.com/AidlyTeam/Aidly-Backend/pkg/validator_service"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
)

func Run(cfg *config.Config) {
	// Postgres Client
	connStr := fmt.Sprintf("user=%v password=%v dbname=%v port=%v sslmode=%v host=%v", cfg.DatabaseConfig.Managment.ManagmentUsername, cfg.DatabaseConfig.Managment.ManagmentPassword, cfg.DatabaseConfig.DBName, cfg.DatabaseConfig.Port, cfg.DatabaseConfig.SSLMode, cfg.DatabaseConfig.Host)
	conn, err := sql.Open(cfg.DatabaseConfig.Driver, connStr)
	if err != nil {
		return
	}
	if err := conn.Ping(); err != nil && err.Error() != "pq: database system is starting up" {
		panic(err)
	}
	if err := goose.Up(conn, cfg.Application.MigrationsPath); err != nil {
		panic(err)
	}
	// Repos
	queries := repo.New(conn)

	// Utilities Initialize
	validatorService := validatorService.NewValidatorService()

	// Service Initialize
	allServices := services.CreateNewServices(
		validatorService,
		queries,
		conn,
		cfg,
	)

	// First Run & Creating Default Admin
	firstRun(queries, allServices.RoleService(), allServices.UserService(), cfg)

	// Handler Initialize
	handlers := http.NewHandler(allServices, cfg)

	// Fiber İnitialize
	fiberServer := server.NewServer(cfg, response.ResponseHandler)

	// Captcha Initialize
	go func() {
		err := fiberServer.Run(handlers.Init(cfg.Application.DevMode, middlewares.InitMiddlewares(cfg)...))
		if err != nil {
			log.Fatalf("Error while running fiber server: %v", err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("Gracefully shutting down...")
	_ = fiberServer.Shutdown(context.Background())
	fmt.Println("Fiber was successful shutdown.")
}

func firstRun(repo *repo.Queries, roleService *services.RoleService, userService *services.UserService, cfg *config.Config) {
	ctx := context.Background()

	defaultRole, err := roleService.GetByName(ctx, "admin")
	if err != nil {
		log.Fatalf("Error default role not exists: %v", err)
	}

	ok, err := repo.IsDefaultUserExists(ctx)
	if err != nil {
		log.Fatalf("Error checking for admin user: %v", err)
	}
	if !ok {
		user, err := userService.AdminCreate(ctx, cfg.Application.Managment.WalletAddress, defaultRole.ID)
		if err != nil {
			log.Fatalf("Error creating for admin user: %v", err)
		}
		log.Println(user)
	}
}
