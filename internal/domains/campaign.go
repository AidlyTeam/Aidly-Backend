package domains

import (
	"time"

	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/google/uuid"
)

type Campaign struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Title               string
	Description         string
	WalletAddress       string
	ImagePath           string
	TargetAmount        string
	RaisedAmount        string
	AcceptedTokenSymbol string
	IsVerified          bool
	IsValid             bool
	StatusType          string
	StartDate           time.Time
	EndDate             *time.Time
	CreatedAt           time.Time
	Categories          []Category
}

type Campaigns struct {
	Campaign   []Campaign
	TotalCount int
}

func ToCampaign(dbModel *repo.TCampaign, categoryModel []Category) *Campaign {
	appModel := Campaign{
		ID:                  dbModel.ID,
		UserID:              dbModel.UserID,
		Title:               dbModel.Title,
		Description:         dbModel.Description.String,
		WalletAddress:       dbModel.WalletAddress,
		ImagePath:           dbModel.ImagePath.String,
		TargetAmount:        dbModel.TargetAmount,
		RaisedAmount:        dbModel.RaisedAmount.String,
		AcceptedTokenSymbol: dbModel.AcceptedTokenSymbol,
		IsVerified:          dbModel.IsVerified,
		IsValid:             dbModel.IsValid,
		StatusType:          dbModel.StatusType,
		StartDate:           dbModel.StartDate.Time,
		EndDate:             &dbModel.EndDate.Time,
		CreatedAt:           dbModel.CreatedAt.Time,
		Categories:          categoryModel,
	}

	return &appModel
}

func MapCampaignsToDomain(campaigns []repo.TCampaign, categories []repo.GetCampaignCategoriesRow) []Campaign {
	var mappedCampaigns []Campaign

	for _, campaign := range campaigns {
		var relatedCategories []Category
		for _, category := range categories {
			if category.CampaignID == campaign.ID {
				relatedCategories = append(relatedCategories, Category{CategoryID: category.CampaignID, CampaignID: category.CampaignID, Name: category.CategoryName})
			}
		}

		// Map'lenen kampanyayı struct olarak ekliyoruz.
		mappedCampaigns = append(mappedCampaigns, Campaign{
			ID:                  campaign.ID,
			UserID:              campaign.UserID,
			Title:               campaign.Title,
			Description:         campaign.Description.String,
			WalletAddress:       campaign.WalletAddress,
			ImagePath:           campaign.ImagePath.String,
			TargetAmount:        campaign.TargetAmount,
			RaisedAmount:        campaign.RaisedAmount.String,
			AcceptedTokenSymbol: campaign.AcceptedTokenSymbol,
			IsVerified:          campaign.IsVerified,
			IsValid:             campaign.IsValid,
			StatusType:          campaign.StatusType,
			StartDate:           campaign.StartDate.Time,
			EndDate:             &campaign.EndDate.Time,
			CreatedAt:           campaign.CreatedAt.Time,
			Categories:          relatedCategories,
		})
	}

	return mappedCampaigns
}
