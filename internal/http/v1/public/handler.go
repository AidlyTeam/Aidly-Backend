package public

import (
	"github.com/AidlyTeam/Aidly-Backend/internal/config"
	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/AidlyTeam/Aidly-Backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type PublicHandler struct {
	services     *services.Services
	sessionStore *session.Store
	dtoManager   dto.IDTOManager
	config       *config.Config
}

func NewPublicHandler(
	service *services.Services,
	sessionStore *session.Store,
	dtoManager dto.IDTOManager,
	config *config.Config,
) *PublicHandler {
	return &PublicHandler{
		services:     service,
		sessionStore: sessionStore,
		dtoManager:   dtoManager,
		config:       config,
	}
}

func (h *PublicHandler) Init(router fiber.Router) {
	root := router.Group("/public")

	root.Get("/", func(c *fiber.Ctx) error {
		return response.Response(200, "Welcome to Aidly API (Public Zone)", nil)
	})

	h.initUserRoutes(root)
}
