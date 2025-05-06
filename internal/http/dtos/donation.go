package dto

import (
	"fmt"
	"time"

	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/google/uuid"
)

type DonationDTOManager struct{}

// NewDonationDTOManager returns an instance of DonationDTOManager
func NewDonationDTOManager() DonationDTOManager {
	return DonationDTOManager{}
}

// DonationView represents the structure of a donation to be returned in the response
type DonationView struct {
	ID            string    `json:"id"`
	CampaignID    string    `json:"campaignID"`
	UserID        string    `json:"userID"`
	Amount        string    `json:"amount"`
	DonationDate  time.Time `json:"donationDate"`
	TransactionID string    `json:"transactionID"`
}

type DonationViews struct {
	Donations  []DonationView `json:"donations"`
	TotalCount int            `json:"totalCount"`
}

// ToDonationView transforms the raw donation data into a structured response view
func (DonationDTOManager) ToDonationView(donation *repo.TDonation) DonationView {
	return DonationView{
		ID:            donation.ID.String(),
		UserID:        donation.UserID.String(),
		CampaignID:    donation.CampaignID.String(),
		TransactionID: donation.TransactionID,
		Amount:        fmt.Sprintf("%v", donation.Amount),
		DonationDate:  donation.DonationDate.Time,
	}
}

func (m *DonationDTOManager) ToDonationViews(donation []repo.TDonation, count int64) *DonationViews {
	var donations []DonationView
	for _, model := range donation {
		donations = append(donations, m.ToDonationView(&model))
	}

	return &DonationViews{Donations: donations, TotalCount: int(count)}
}

// DonationCreateDTO represents the structure for creating a new donation
type DonationCreateDTO struct {
	CampaignID    string `json:"campaignID" validate:"required"`
	Amount        string `json:"amount" validate:"required,numeric"`
	TransactionID string `json:"transactionID" validate:"required"`
}

type DonationRequest struct {
	BadgeID    uuid.UUID `json:"badgeID"`
	DonationID uuid.UUID `json:"donationID"`
}
