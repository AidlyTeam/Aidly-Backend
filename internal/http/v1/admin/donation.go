package admin

import (
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

func (h *AdminHandler) initDonationsRoutes(root fiber.Router) {
	donation := root.Group("/donation")

	donation.Get("/", h.GetDonations)
	donation.Get("/:donationID", h.GetDonationByID)

	donation.Delete("/:donationID", h.DeleteDonation)
}

// @Tags Donation
// @Summary Get Donations
// @Description Retrieves a list of donations based on given filters.
// @Accept json
// @Produce json
// @Param id query string false "Donation ID"
// @Param userID query string false "User ID"
// @Param campaignID query string false "Campaign ID"
// @Param page query string false "Page Number"
// @Param limit query string false "Limit Per Page"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/donation [get]
func (h *AdminHandler) GetDonations(c *fiber.Ctx) error {
	id := c.Query("id")
	userID := c.Query("userID")
	campaignID := c.Query("campaignID")
	page := c.Query("page")
	limit := c.Query("limit")

	donations, err := h.services.DonationService().GetDonations(c.Context(), id, userID, campaignID, page, limit)
	if err != nil {
		return err
	}
	count, err := h.services.DonationService().CountDonations(c.Context(), campaignID, userID)
	if err != nil {
		return err
	}
	donationsView := h.dtoManager.DonationManager().ToDonationViews(donations, count)

	return response.Response(200, "Donations Retrieved Successfully", donationsView)
}

// @Tags Donation
// @Summary Get Donation by ID
// @Description Retrieves a donation based on the provided donation ID. User can only view their own donation.
// @Accept json
// @Produce json
// @Param donationID path string true "Donation ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/donation/{donationID} [get]
func (h *AdminHandler) GetDonationByID(c *fiber.Ctx) error {
	donationID := c.Params("donationID")

	donation, err := h.services.DonationService().GetDonationByID(c.Context(), donationID)
	if err != nil {
		return err
	}
	donationView := h.dtoManager.DonationManager().ToDonationView(nil, donation)

	return response.Response(200, "Donation Retrieved Successfully", donationView)
}

// @Tags Donation
// @Summary Delete a Donation
// @Description Deletes a donation based on the provided donation ID.
// @Accept json
// @Produce json
// @Param donationID path string true "Donation ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/donation/{donationID} [delete]
func (h *AdminHandler) DeleteDonation(c *fiber.Ctx) error {
	donationID := c.Params("donationID")

	if err := h.services.DonationService().DeleteDonation(c.Context(), donationID); err != nil {
		return err
	}

	return response.Response(200, "Donation Deleted Successfully", nil)
}
