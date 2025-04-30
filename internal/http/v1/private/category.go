package private

import (
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

func (h *PrivateHandler) initCategoryRoutes(root fiber.Router) {
	category := root.Group("/category")

	category.Get("/", h.GetCategories)
}

// @Tags Category
// @Summary Get Categories
// @Description Retrieves a list of categories with optional pagination.
// @Accept json
// @Produce json
// @Param page query string false "Page Number"
// @Param limit query string false "Limit Per Page"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/category [get]
func (h *PrivateHandler) GetCategories(c *fiber.Ctx) error {
	page := c.Query("page")
	limit := c.Query("limit")

	categories, err := h.services.CategoryService().GetCategories(c.Context(), page, limit)
	if err != nil {
		return err
	}
	count, err := h.services.CategoryService().CountCategories(c.Context())
	if err != nil {
		return err
	}
	view := h.dtoManager.CategoryManager().ToCategoryViews(categories, count)

	return response.Response(200, "Categories Retrieved Successfully", view)
}
