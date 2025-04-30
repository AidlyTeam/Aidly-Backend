package admin

import (
	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/gofiber/fiber/v2"
)

func (h *AdminHandler) initCategoryRoutes(root fiber.Router) {
	category := root.Group("/category")

	category.Get("/:categoryID", h.GetCategoryByID)
	category.Post("/", h.CreateCategory)
	category.Put("/:categoryID", h.UpdateCategory)
	category.Delete("/:categoryID", h.DeleteCategory)
}

// @Tags Category
// @Summary Get Category by ID
// @Description Retrieves a category by its ID.
// @Accept json
// @Produce json
// @Param categoryID path string true "Category ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/category/{categoryID} [get]
func (h *AdminHandler) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("categoryID")

	category, err := h.services.CategoryService().GetCategoryByID(c.Context(), id)
	if err != nil {
		return err
	}
	view := h.dtoManager.CategoryManager().ToCategoryView(category)

	return response.Response(200, "Category Retrieved Successfully", view)
}

// @Tags Category
// @Summary Create Category
// @Description Creates a new category.
// @Accept json
// @Produce json
// @Param req body dto.CategoryCreateDTO true "Create Category Request"
// @Success 201 {object} response.BaseResponse{}
// @Router /admin/category [post]
func (h *AdminHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CategoryCreateDTO
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(req); err != nil {
		return err
	}

	id, err := h.services.CategoryService().CreateCategory(c.Context(), req.Name)
	if err != nil {
		return err
	}

	return response.Response(201, "Category Created Successfully", id)
}

// @Tags Category
// @Summary Update Category
// @Description Updates an existing category.
// @Accept json
// @Produce json
// @Param categoryID path string true "Category ID"
// @Param req body dto.CategoryUpdateDTO true "Update Category Request"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/category/{categoryID} [put]
func (h *AdminHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("categoryID")

	var req dto.CategoryUpdateDTO
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(req); err != nil {
		return err
	}

	if err := h.services.CategoryService().UpdateCategory(c.Context(), id, req.Name); err != nil {
		return err
	}

	return response.Response(200, "Category Updated Successfully", nil)
}

// @Tags Category
// @Summary Delete Category
// @Description Deletes a category by its ID.
// @Accept json
// @Produce json
// @Param categoryID path string true "Category ID"
// @Success 200 {object} response.BaseResponse{}
// @Router /admin/category/{categoryID} [delete]
func (h *AdminHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("categoryID")

	if err := h.services.CategoryService().DeleteCategory(c.Context(), id); err != nil {
		return err
	}

	return response.Response(200, "Category Deleted Successfully", nil)
}
