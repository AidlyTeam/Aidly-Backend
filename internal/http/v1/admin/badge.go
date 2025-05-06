package admin

import (
	"strings"

	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

func (h *AdminHandler) initBadgeRoutes(root fiber.Router) {
	badge := root.Group("/badge")

	badge.Get("/", h.GetBadges)
	badge.Get("/:badgeID", h.GetBadgeByID)
	badge.Post("/", h.CreateBadge)
	badge.Put("/:badgeID", h.UpdateBadge)
	badge.Delete("/:badgeID", h.DeleteBadge)

	// TODO: CREATE BADGE JSON AND SEND IT TO A USER ACCORDINGLY. NEED USERS DONATION COUNT FOR THIS BADGE SYSTEM
	// MAKE ORDER IN BADGE AND IN CODE YOU WILL SEND IT ALGORITHMYLICLY
}

// @Tags Badge
// @Summary Get all badges
// @Description Get list of all badges
// @Accept json
// @Produce json
// @Param id query string false "Badge ID"
// @Param page query string false "Page Number"
// @Param limit query string false "Limit Per Page"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/badge [get]
func (h *AdminHandler) GetBadges(c *fiber.Ctx) error {
	id := c.Query("id")
	page := c.Query("page")
	limit := c.Query("limit")

	badges, err := h.services.BadgeService().GetBadges(c.Context(), id, page, limit)
	if err != nil {
		return err
	}
	count, err := h.services.BadgeService().GetBadgeCount(c.Context(), id)
	if err != nil {
		return err
	}

	badgeViews := h.dtoManager.BadgeManager().ToBadgeViews(badges, count)

	return response.Response(200, "Badges fetched successfully", badgeViews)
}

// @Tags Badge
// @Summary Get badge by ID
// @Description Get badge details by badge ID
// @Accept json
// @Produce json
// @Param badgeID path string true "Badge ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/badge/{badgeID} [get]
func (h *AdminHandler) GetBadgeByID(c *fiber.Ctx) error {
	badgeID := c.Params("badgeID")
	badge, err := h.services.BadgeService().GetBadgeByID(c.Context(), badgeID)
	if err != nil {
		return err
	}
	badgeView := h.dtoManager.BadgeManager().ToBadgeView(badge)

	return response.Response(200, "Badge fetched successfully", badgeView)
}

// @Tags Badge
// @Summary Create a new badge
// @Description Creates a new badge
// @Accept json
// @Produce json
// @Param imageFile formData file true "Badge Image File"
// @Param badge formData dto.BadgeCreateDTO true "New badge data"
// @Success 201 {object} response.BaseResponse{}
// @Router /admin/badge [post]
func (h *AdminHandler) CreateBadge(c *fiber.Ctx) error {
	var badgeDTO dto.BadgeCreateDTO
	if err := c.BodyParser(&badgeDTO); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(badgeDTO); err != nil {
		return err
	}

	imageFile, err := c.FormFile("imageFile")
	if err != nil {
		return err
	}

	imagePath := h.services.UploadService().CreatePath(imageFile.Filename)

	id, err := h.services.BadgeService().CreateBadge(
		c.Context(),
		badgeDTO.Name,
		badgeDTO.Description,
		imagePath,
		badgeDTO.DonationThreshold,
	)
	if err != nil {
		return err
	}

	if err := h.services.UploadService().SaveImage(imageFile, imagePath); err != nil {
		return err
	}

	return response.Response(201, "Badge created successfully", id)
}

// @Tags Badge
// @Summary Update badge
// @Description Updates an existing badge
// @Accept json
// @Produce json
// @Param badgeID path string true "Badge ID"
// @Param imageFile formData file false "Badge Image File"
// @Param badge formData dto.BadgeUpdateDTO true "Updated badge data"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/badge/{badgeID} [put]
func (h *AdminHandler) UpdateBadge(c *fiber.Ctx) error {
	badgeID := c.Params("badgeID")

	var badgeDTO dto.BadgeUpdateDTO
	if err := c.BodyParser(&badgeDTO); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(badgeDTO); err != nil {
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

	if err := h.services.BadgeService().UpdateBadge(
		c.Context(),
		badgeID,
		badgeDTO.Name,
		badgeDTO.Description,
		imagePath,
		badgeDTO.DonationThreshold,
	); err != nil {
		return err
	}

	if imageFile != nil {
		if err := h.services.UploadService().SaveImage(imageFile, imagePath); err != nil {
			return err
		}
	}

	return response.Response(200, "Badge updated successfully", nil)
}

// @Tags Badge
// @Summary Delete badge
// @Description Deletes a badge by ID
// @Accept json
// @Produce json
// @Param badgeID path string true "Badge ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/badge/{badgeID} [delete]
func (h *AdminHandler) DeleteBadge(c *fiber.Ctx) error {
	badgeID := c.Params("badgeID")
	err := h.services.BadgeService().DeleteBadge(c.Context(), badgeID)
	if err != nil {
		return err
	}

	return response.Response(200, "Badge deleted successfully", nil)
}
