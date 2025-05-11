package private

import (
	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/sessionStore"
	"github.com/gofiber/fiber/v2"
)

func (h *PrivateHandler) initBadgeRoutes(root fiber.Router) {
	badge := root.Group("/badge")

	badge.Get("/user", h.GetUserBadges)
	badge.Get("/mint/:badgeID", h.MintNft)
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
	badgeViews := h.dtoManager.BadgeManager().ToUserBadgeViews(badges, count)

	return response.Response(200, "User badges fetched successfully", badgeViews)
}

// @Tags Badge
// @Summary Mint badge NFT
// @Description Mints an NFT for a specific badge owned by the user
// @Accept json
// @Produce json
// @Param badgeID path string true "Badge ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/badge/mint/{badgeID} [get]
func (h *PrivateHandler) MintNft(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)

	badgeID := c.Params("badgeID")

	userBadge, err := h.services.BadgeService().GetUserBadge(c.Context(), badgeID, userSession.UserID)
	if err != nil {
		return err
	}
	if !userBadge.IsNft {
		return response.Response(400, serviceErrors.ErrBadgeIsNotNFT, nil)
	}
	if userBadge.IsMinted {
		return response.Response(400, serviceErrors.ErrNFTAlreadyMinted, nil)
	}

	resp, err := h.services.BadgeService().MintNFT(
		c.Context(),
		userBadge.Name,
		userBadge.Symbol.String,
		userBadge.Uri.String,
		userBadge.SellerFee.Int32,
		userSession.WalletAddress,
	)
	if err != nil {
		return err
	}

	if err := h.services.BadgeService().ChangeIsMinted(c.Context(), userSession.UserID, userBadge.ID); err != nil {
		return nil
	}

	return response.Response(200, "NFT Minted", resp)
}
