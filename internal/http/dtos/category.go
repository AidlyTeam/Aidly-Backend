package dto

import (
	"github.com/AidlyTeam/Aidly-Backend/internal/domains"
)

type CategoryDTOManager struct{}

// NewCategoryDTOManager returns an instance of CategoryDTOManager
func NewCategoryDTOManager() CategoryDTOManager {
	return CategoryDTOManager{}
}

// CategoryView represents the structure of a category to be returned in the response
type CategoryView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryViews struct {
	Categories []CategoryView `json:"categories"`
	TotalCount int            `json:"totalCount,omitempty"`
}

// ToCategoryView transforms the raw category data into a structured response view
func (CategoryDTOManager) ToCategoryView(category *domains.Category) CategoryView {
	return CategoryView{
		ID:   category.CategoryID.String(),
		Name: category.Name,
	}
}

func (m *CategoryDTOManager) ToCategoryViews(categories []domains.Category, count int64) *CategoryViews {
	var categoryViews []CategoryView
	for _, model := range categories {
		categoryViews = append(categoryViews, m.ToCategoryView(&model))
	}

	return &CategoryViews{Categories: categoryViews, TotalCount: int(count)}
}

// CategoryCreateDTO represents the structure for creating a new category
type CategoryCreateDTO struct {
	Name string `json:"name" validate:"required,min=3,max=8"`
}

// CategoryUpdateDTO represents the structure for updating an existing category
type CategoryUpdateDTO struct {
	Name string `json:"name" validate:"omitempty,min=3,max=8"`
}
