package admin

import (
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

func (h *AdminHandler) initCampaignsRoutes(root fiber.Router) {
	campaign := root.Group("/campaign")

	campaign.Delete("/:campaignID", h.DeleteCampaign)
}

// @Tags Campaign
// @Summary Delete a Campaign
// @Description Deletes a campaign based on the provided campaign ID.
// @Accept json
// @Produce json
// @Param campaignID path string true "Campaign ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/campaign/{campaignID} [delete]
func (h *AdminHandler) DeleteCampaign(c *fiber.Ctx) error {
	campaignID := c.Params("campaignID")

	if err := h.services.CampaignService().DeleteCampaign(c.Context(), campaignID); err != nil {
		return err
	}

	return response.Response(200, "Campaign Deleted Successfully", nil)
}
