package private

import (
	"strings"

	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/sessionStore"
	solana_service "github.com/AidlyTeam/Aidly-Backend/pkg/web3"
	"github.com/gofiber/fiber/v2"
)

func (h *PrivateHandler) initCampaignsRoutes(root fiber.Router) {
	campaign := root.Group("/campaign")

	campaign.Get("/", h.GetCampaigns)
	campaign.Get("/:campaignID", h.GetCampaignByID)

	campaign.Post("/", h.CreateCampaign)
	campaign.Put("/:campaignID", h.UpdateCampaign)
	campaign.Delete("/:campaignID", h.DeleteCampaign)

	campaign.Get("/:campaignID/isValid", h.CheckCampaignValidity)

	campaign.Post("/:campaignID/category", h.AddCategoryToCampaign)
	campaign.Delete("/:campaignID/category", h.RemoveCategoryFromCampaign)
}

// @Tags Campaign
// @Summary Get Campaigns
// @Description Retrieves a list of campaigns based on given filters.
// @Accept json
// @Produce json
// @Param id query string false "Campaign ID"
// @Param userID query string false "User ID"
// @Param categoryIDList query string false "Category ID List"
// @Param status query string false "Status"
// @Param isVerified query string false "Campaign Verifiy"
// @Param page query string false "Page Number"
// @Param limit query string false "Limit Per Page"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/campaign [get]
func (h *PrivateHandler) GetCampaigns(c *fiber.Ctx) error {
	id := c.Query("id")
	userID := c.Query("userID")
	categoryIDList := c.Query("categoryIDList")
	status := c.Query("status")
	isVerified := c.Query("isVerified")
	page := c.Query("page")
	limit := c.Query("limit")

	campaigns, err := h.services.CampaignService().GetCampaigns(c.Context(), id, userID, isVerified, status, page, limit)
	if err != nil {
		return err
	}
	campaignsView := h.dtoManager.CampaignManager().ToCampaignViews(campaigns, categoryIDList)

	return response.Response(200, "Campaigns Retrieved Successfully", campaignsView)
}

// @Tags Campaign
// @Summary Get Campaign by ID
// @Description Retrieves a campaign based on the provided campaign ID.
// @Accept json
// @Produce json
// @Param campaignID path string true "Campaign ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/campaign/{campaignID} [get]
func (h *PrivateHandler) GetCampaignByID(c *fiber.Ctx) error {
	campaignID := c.Params("campaignID")

	campaign, err := h.services.CampaignService().GetCampaignByID(c.Context(), campaignID)
	if err != nil {
		return err
	}
	campaignView := h.dtoManager.CampaignManager().ToCampaignView(campaign)

	return response.Response(200, "Campaign Retrieved Successfully", campaignView)
}

// @Tags Campaign
// @Summary Create a Campaign for User
// @Description Creates a new campaign.
// @Accept json
// @Produce json
// @Param imageFile formData file true "Campaign Image File"
// @Param campaign formData dto.CampaignCreateDTO true "New Campaign"
// @Success 201 {object} response.BaseResponse{}
// @Router /private/campaign [post]
func (h *PrivateHandler) CreateCampaign(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)

	var newCampaign dto.CampaignCreateDTO
	if err := c.BodyParser(&newCampaign); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(newCampaign); err != nil {
		return err
	}

	imageFile, err := c.FormFile("imageFile")
	if err != nil {
		return err
	}

	if ok := solana_service.IsValidSolanaAddress(newCampaign.WalletAddress); !ok {
		return response.Response(400, serviceErrors.ErrNotValidWalletAddress, nil)
	}

	imagePath := h.services.UploadService().CreatePath(imageFile.Filename)

	// TODO: CHECK IF THE WALLET ADDRESS IS REALY SOLANA WALLET ADRESS. CREATE WEB3 SERVICE
	campaignID, err := h.services.CampaignService().CreateCampaign(
		c.Context(),
		userSession.UserID,
		newCampaign.Title,
		newCampaign.Description,
		newCampaign.WalletAddress,
		imagePath,
		newCampaign.TargetAmount,
		newCampaign.StatusType,
		newCampaign.AcceptedTokenSymbol,
		newCampaign.StartDate,
		newCampaign.EndDate,
	)
	if err != nil {
		return err
	}

	if err := h.services.UploadService().SaveImage(imageFile, imagePath); err != nil {
		return err
	}

	return response.Response(201, "Campaign Created Successfully", campaignID)
}

// @Tags Campaign
// @Summary Update a Campaign for User
// @Description Updates the details of an existing campaign.
// @Accept json
// @Produce json
// @Param campaignID path string true "Campaign ID"
// @Param imageFile formData file false "Campaign Image File"
// @Param campaign formData dto.CampaignUpdateDTO true "Updated Campaign"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/campaign/{campaignID} [put]
func (h *PrivateHandler) UpdateCampaign(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)
	campaignID := c.Params("campaignID")

	var updateCampaign dto.CampaignUpdateDTO
	if err := c.BodyParser(&updateCampaign); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(updateCampaign); err != nil {
		return err
	}

	imageFile, err := c.FormFile("imageFile")
	if err != nil && !strings.Contains(err.Error(), "there is no uploaded file associated with the given key") {
		return err
	}

	var imagePath string
	if imageFile != nil {
		imagePath = h.services.UploadService().CreatePath(imageFile.Filename)
	}

	if updateCampaign.WalletAddress != "" {
		if ok := solana_service.IsValidSolanaAddress(updateCampaign.WalletAddress); !ok {
			return response.Response(400, serviceErrors.ErrNotValidWalletAddress, nil)
		}
	}

	if err := h.services.CampaignService().UpdateCampaign(
		c.Context(),
		userSession.UserID,
		campaignID,
		updateCampaign.Title,
		updateCampaign.Description,
		updateCampaign.WalletAddress,
		imagePath,
		updateCampaign.TargetAmount,
		updateCampaign.StatusType,
		updateCampaign.AcceptedTokenSymbol,
		updateCampaign.StartDate,
		updateCampaign.EndDate,
	); err != nil {
		return err
	}

	if imageFile != nil {
		if err := h.services.UploadService().SaveImage(imageFile, imagePath); err != nil {
			return err
		}
	}

	return response.Response(200, "Campaign Updated Successfully", nil)
}

// @Tags Campaign
// @Summary Delete a Campaign for User
// @Description Deletes a campaign based on the provided campaign ID.
// @Accept json
// @Produce json
// @Param campaignID path string true "Campaign ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/campaign/{campaignID} [delete]
func (h *PrivateHandler) DeleteCampaign(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)
	campaignID := c.Params("campaignID")

	if err := h.services.CampaignService().CheckTheOwnerOfCampaign(c.Context(), campaignID, userSession.UserID); err != nil {
		return err
	}

	if err := h.services.CampaignService().DeleteCampaign(c.Context(), campaignID); err != nil {
		return err
	}

	return response.Response(200, "Campaign Deleted Successfully", nil)
}

// @Tags Campaign
// @Summary Check if the Campaign is Valid
// @Description Checks if the campaign is valid based on its end date, target amount raised, etc. Request This Before Donation!!!!!!!!
// @Accept json
// @Produce json
// @Param campaignID path string true "Campaign ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/campaign/{campaignID}/isValid [Get]
func (h *PrivateHandler) CheckCampaignValidity(c *fiber.Ctx) error {
	campaignID := c.Params("campaignID")

	isValid, err := h.services.CampaignService().CheckCampaignValidity(c.Context(), campaignID)
	if err != nil {
		return err
	}

	if err := h.services.CampaignService().UpdateCampaignValidity(c.Context(), campaignID); err != nil {
		return err
	}

	return response.Response(200, "Successful", isValid)
}

// @Tags Campaign
// @Summary Add Category to Campaign
// @Description Adds a category to the campaign.
// @Accept json
// @Produce json
// @Param campaignID path string true "Campaign ID"
// @Param req body dto.CampaignCategoryAddDelete true "Category ID"
// @Success 201 {object} response.BaseResponse{}
// @Router /private/campaign/{campaignID}/category [post]
func (h *PrivateHandler) AddCategoryToCampaign(c *fiber.Ctx) error {
	campaignID := c.Params("campaignID")

	var req dto.CampaignCategoryAddDelete
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(req); err != nil {
		return err
	}

	category, err := h.services.CategoryService().GetCategoryByID(c.Context(), req.CategoryID)
	if err != nil {
		return err
	}

	campaignCategoryID, err := h.services.CampaignService().AddCategory(c.Context(), campaignID, category.CategoryID)
	if err != nil {
		return err
	}

	return response.Response(201, "Category Added to Campaign Successfully", campaignCategoryID)
}

// @Tags Campaign
// @Summary Remove Category from Campaign
// @Description Removes a category from the campaign.
// @Accept json
// @Produce json
// @Param campaignID path string true "Campaign ID"
// @Param req body dto.CampaignCategoryAddDelete true "Category ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/campaign/{campaignID}/category [delete]
func (h *PrivateHandler) RemoveCategoryFromCampaign(c *fiber.Ctx) error {
	campaignID := c.Params("campaignID")

	var req dto.CampaignCategoryAddDelete
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(req); err != nil {
		return err
	}

	category, err := h.services.CategoryService().GetCategoryByID(c.Context(), req.CategoryID)
	if err != nil {
		return err
	}

	if err := h.services.CampaignService().RemoveCategory(c.Context(), campaignID, category.CategoryID); err != nil {
		return err
	}

	return response.Response(200, "Category Removed from Campaign Successfully", nil)
}
