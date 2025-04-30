package domains

import (
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/google/uuid"
)

type Category struct {
	CategoryID uuid.UUID
	CampaignID uuid.UUID
	Name       string
}

func ToCategory(dbModel *repo.TCategory) *Category {
	return &Category{
		CategoryID: dbModel.ID,
		Name:       dbModel.Name,
	}
}

func ToCategories(dbModels []repo.TCategory) []Category {
	var categoryTmp []Category
	for _, category := range dbModels {
		categoryTmp = append(categoryTmp, Category{
			CategoryID: category.ID,
			Name:       category.Name,
		})
	}

	return categoryTmp
}

func ToCategoryCampaign(dbModel *repo.GetCampaignCategoriesRow) *Category {
	return &Category{
		CategoryID: dbModel.CategoryID,
		CampaignID: dbModel.CampaignID,
		Name:       dbModel.CategoryName,
	}
}

func ToCategoriesCampaign(dbModels []repo.GetCampaignCategoriesRow) []Category {
	var categoryTmp []Category
	for _, category := range dbModels {
		categoryTmp = append(categoryTmp, Category{
			CategoryID: category.CategoryID,
			CampaignID: category.CampaignID,
			Name:       category.CategoryName,
		})
	}

	return categoryTmp
}
