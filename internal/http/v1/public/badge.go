package public

import (
	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

func (h *PublicHandler) initBadgeRoutes(root fiber.Router) {
	badge := root.Group("/badge")

	badge.Get("/:badgeID", h.GetBadgeByID)
}

// @Tags Badge
// @Summary Get Badge NFT
// @Description Get NFT Metadata
// @Accept json
// @Produce json
// @Param badgeID path string true "Badge ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /public/badge/{badgeID} [get]
func (h *PublicHandler) GetBadgeByID(c *fiber.Ctx) error {
	badgeID := c.Params("badgeID")
	badge, err := h.services.BadgeService().GetBadgeByID(c.Context(), badgeID)
	if err != nil {
		return err
	}
	badgeView := h.dtoManager.BadgeManager().ToMetadataView(badge)

	if badgeView == nil {
		return response.Response(404, serviceErrors.ErrBadgeNotFound, nil)
	}

	return response.Response(200, "Badge fetched successfully", badgeView)
}
