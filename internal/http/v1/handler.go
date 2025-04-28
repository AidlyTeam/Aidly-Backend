package v1

import (
	"github.com/AidlyTeam/Aidly-Backend/internal/config"
	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/v1/admin"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/v1/private"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/v1/public"
	"github.com/AidlyTeam/Aidly-Backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type V1Handler struct {
	services   *services.Services
	dtoManager dto.IDTOManager
	config     *config.Config
}

func NewV1Handler(
	services *services.Services,
	dtoManager dto.IDTOManager,
	config *config.Config,
) *V1Handler {
	return &V1Handler{
		services:   services,
		dtoManager: dtoManager,
		config:     config,
	}
}

func (h *V1Handler) Init(router fiber.Router, sessionStore *session.Store) {
	root := router.Group("/v1")
	root.Get("/", func(c *fiber.Ctx) error {
		return response.Response(200, "Welcome to Aidly API (Root Zone)", nil)
	})

	// Init Fiber Session Store
	admin.NewAdminHandler(
		h.services,
		sessionStore,
		h.dtoManager,
		h.config,
	).Init(root)
	public.NewPublicHandler(h.services,
		sessionStore,
		h.dtoManager,
		h.config,
	).Init(root)
	private.NewPrivateHandler(h.services,
		sessionStore,
		h.dtoManager,
		h.config,
	).Init(root)
}
