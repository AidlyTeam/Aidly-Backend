package http

import (
	"github.com/AidlyTeam/Aidly-Backend/docs"
	"github.com/AidlyTeam/Aidly-Backend/internal/config"
	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/sessionStore"
	v1 "github.com/AidlyTeam/Aidly-Backend/internal/http/v1"
	"github.com/AidlyTeam/Aidly-Backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Handler struct {
	services *services.Services
	config   *config.Config
}

func NewHandler(
	services *services.Services,
	config *config.Config,
) *Handler {
	return &Handler{
		services: services,
		config:   config,
	}
}

func (h *Handler) Init(devMode bool, middlewares ...func(*fiber.Ctx) error) *fiber.App {
	app := fiber.New()
	for i := range middlewares {
		app.Use(middlewares[i])
	}

	if devMode {
		docs.SwaggerInfo.Version = config.Version
		app.Get("/api/dev/*", swagger.New(swagger.Config{
			Title:                "Aidly Backend",
			TryItOutEnabled:      true,
			PersistAuthorization: true,
		}))
	}

	app.Static("/api/uploads", "./uploads")

	root := app.Group("/api")
	sessionStore := sessionStore.NewSessionStore()
	dtoManager := dto.CreateNewDTOManager()

	// init routes
	v1.NewV1Handler(h.services, dtoManager, h.config).Init(root, sessionStore)

	return app
}
