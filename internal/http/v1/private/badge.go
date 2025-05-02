package private

import (
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/sessionStore"
	"github.com/gofiber/fiber/v2"
)

func (h *PrivateHandler) initBadgeRoutes(root fiber.Router) {
	badge := root.Group("/badge")

	badge.Get("/", h.GetUserBadges)
}

// @Tags Badge
// @Summary Get user's badges
// @Description Retrieves all badges owned by a specific user
// @Accept json
// @Produce json
// @Success 200 {object} response.BaseResponse{}
// @Router /private/badge/user/ [get]
func (h *PrivateHandler) GetUserBadges(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)

	badges, err := h.services.BadgeService().GetUserBadges(c.Context(), userSession.UserID)
	if err != nil {
		return err
	}
	count, err := h.services.BadgeService().GetBadgeCount(c.Context(), userSession.UserID.String())
	if err != nil {
		return err
	}

	badgeViews := h.dtoManager.BadgeManager().ToBadgeViews(badges, count)

	return response.Response(200, "User badges fetched successfully", badgeViews)
}
