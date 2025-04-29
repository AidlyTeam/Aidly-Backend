package dto

import (
	"time"

	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
)

// CampaignDTOManager struct will manage DTO conversion for Campaign.
type CampaignDTOManager struct{}

// NewCampaignDTOManager creates and returns a new CampaignDTOManager.
func NewCampaignDTOManager() CampaignDTOManager {
	return CampaignDTOManager{}
}

// CampaignView struct will define the response format for campaign details.
type CampaignView struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userID"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ImagePath     string    `json:"imagePath"`
	TargetAmount  string    `json:"targetAmount"`
	RaisedAmount  string    `json:"raisedAmount"`
	IsVerified    bool      `json:"isVerified"`
	AcceptedToken string    `json:"acceptedToken"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
}

// ToCampaignView converts campaign data to a view format for response.
func (m *CampaignDTOManager) ToCampaignView(campaignData *repo.TCampaign) CampaignView {
	return CampaignView{
		ID:            campaignData.ID.String(),
		UserID:        campaignData.UserID.String(),
		Title:         campaignData.Title,
		Description:   campaignData.Description.String,
		ImagePath:     campaignData.ImagePath.String,
		TargetAmount:  campaignData.TargetAmount,
		RaisedAmount:  campaignData.RaisedAmount.String,
		IsVerified:    campaignData.IsVerified,
		AcceptedToken: campaignData.AcceptedTokenSymbol.String,
		StartDate:     campaignData.StartDate.Time,
		EndDate:       campaignData.EndDate.Time,
	}
}

// ToCampaignViews converts a slice of campaign data to an array of views.
func (m *CampaignDTOManager) ToCampaignViews(campaigns []repo.TCampaign) []CampaignView {
	var campaignViews []CampaignView
	for _, campaign := range campaigns {
		campaignViews = append(campaignViews, m.ToCampaignView(&campaign))
	}
	return campaignViews
}

// CampaignCreateDTO struct will define the data required to create a campaign.
type CampaignCreateDTO struct {
	Title         string `json:"title" validate:"required,max=100"`
	Description   string `json:"description" validate:"required,max=500"`
	WalletAddress string `json:"walletAddress" validate:"required,max=500"`
	TargetAmount  string `json:"targetAmount" validate:"required"`
	StartDate     string `json:"startDate" validate:"required"`
	EndDate       string `json:"endDate" validate:"required"`
}

// CampaignUpdateDTO struct will define the data required to update an existing campaign.
type CampaignUpdateDTO struct {
	Title         string `json:"title" validate:"omitempty,max=100"`
	Description   string `json:"description" validate:"omitempty,max=500"`
	WalletAddress string `json:"walletAddress" validate:"required,max=500"`
	TargetAmount  string `json:"targetAmount" validate:"omitempty"`
	StartDate     string `json:"startDate" validate:"omitempty"`
	EndDate       string `json:"endDate" validate:"omitempty"`
}

type CampaignChangeVerify struct {
	IsVerified bool `json:"isVerified"`
}
