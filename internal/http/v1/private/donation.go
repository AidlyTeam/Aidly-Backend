package private

import (
	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/sessionStore"
	"github.com/gofiber/fiber/v2"
)

func (h *PrivateHandler) initDonationsRoutes(root fiber.Router) {
	donation := root.Group("/donation")

	donation.Get("/", h.GetDonations)
	donation.Get("/:donationID", h.GetDonationByID)

	donation.Post("/", h.Donate)
}

// @Tags Donation
// @Summary Get Donations for User
// @Description Retrieves a list of donations based on given filters.
// @Accept json
// @Produce json
// @Param id query string false "Donation ID"
// @Param campaignID query string false "Campaign ID"
// @Param page query string false "Page Number"
// @Param limit query string false "Limit Per Page"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/donation [get]
func (h *PrivateHandler) GetDonations(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)

	id := c.Query("id")
	userID := userSession.UserID.String()
	campaignID := c.Query("campaignID")
	page := c.Query("page")
	limit := c.Query("limit")

	donations, err := h.services.DonationService().GetDonations(c.Context(), id, campaignID, userID, page, limit)
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
// @Summary Get Donation by ID for User
// @Description Retrieves a donation based on the provided donation ID. User can only view their own donation.
// @Accept json
// @Produce json
// @Param donationID path string true "Donation ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/donation/{donationID} [get]
func (h *PrivateHandler) GetDonationByID(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)
	donationID := c.Params("donationID")

	donation, err := h.services.DonationService().GetDonationByID(c.Context(), donationID)
	if err != nil {
		return err
	}

	if donation.UserID != userSession.UserID {
		return response.Response(403, "Forbidden", nil)
	}

	donationView := h.dtoManager.DonationManager().ToDonationView(donation)

	return response.Response(200, "Donation Retrieved Successfully", donationView)
}

// @Tags Donation
// @Summary Create a Donation for User
// @Description Creates a new donation for a user.
// @Accept json
// @Produce json
// @Param donation body dto.DonationCreateDTO true "New Donation"
// @Success 201 {object} response.BaseResponse{}
// @Router /private/donation [post]
func (h *PrivateHandler) Donate(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)

	var newDonation dto.DonationCreateDTO
	if err := c.BodyParser(&newDonation); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(newDonation); err != nil {
		return err
	}

	campaign, err := h.services.CampaignService().GetCampaignByID(c.Context(), newDonation.CampaignID)
	if err != nil {
		return nil
	}

	// TODO: CHECK IF THE DONATION IS SUCCEED WITH TRANSACTION ID. IF SUCCEED. THEN CREATE.
	// Make A web3 servis for this
	donationID, err := h.services.DonationService().CreateDonation(
		c.Context(),
		userSession.UserID,
		newDonation.Amount,
		newDonation.TransactionID,
		campaign,
	)
	if err != nil {
		return err
	}
	if err := h.services.CampaignService().UpdateCampaignValidity(c.Context(), newDonation.CampaignID); err != nil {
		return err
	}

	count, err := h.services.DonationService().CountDonations(c.Context(), "", userSession.UserID.String())
	if err != nil {
		return err
	}

	// Check the donation count if exist add.
	badgeID, err := h.services.BadgeService().CheckBadgeAndAdd(c.Context(), userSession.UserID, int32(count))
	if err != nil {
		return nil
	}

	resp := dto.DonationRequest{
		BadgeID:    *badgeID,
		DonationID: *donationID,
	}

	return response.Response(201, "Donation Created Successfully", resp)
}
