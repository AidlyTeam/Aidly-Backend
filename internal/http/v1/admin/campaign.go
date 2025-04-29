package admin

import (
	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

func (h *AdminHandler) initCampaignsRoutes(root fiber.Router) {
	campaign := root.Group("/campaign")

	campaign.Delete("/:campaignID", h.DeleteCampaign)
	campaign.Patch("/:campaignID/verify", h.ChangeCampaignVerified)
}

// @Tags Campaign
// @Summary Change Campaign Verification Status
// @Description Updates the verification status of a campaign.
// @Accept json
// @Produce json
// @Param campaignID path string true "Campaign ID"
// @Param verify body dto.CampaignChangeVerify true "Updated Campaign"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/campaign/{campaignID}/verify [patch]
func (h *AdminHandler) ChangeCampaignVerified(c *fiber.Ctx) error {
	campaignID := c.Params("campaignID")

	var verify dto.CampaignChangeVerify
	if err := c.BodyParser(&verify); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}
	if err := h.services.UtilService().Validator().ValidateStruct(verify); err != nil {
		return err
	}

	err := h.services.CampaignService().ChangeCampaignVerified(c.Context(), campaignID, verify.IsVerified)
	if err != nil {
		return err
	}

	return response.Response(200, "Campaign verification status updated", nil)
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
