package services

import (
	"context"
	"database/sql"
	"strconv"

	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/google/uuid"
)

type DonationService struct {
	db          *sql.DB
	queries     *repo.Queries
	utilService IUtilService
}

func newDonationService(
	db *sql.DB,
	queries *repo.Queries,
	utilService IUtilService,
) *DonationService {
	return &DonationService{
		db:          db,
		queries:     queries,
		utilService: utilService,
	}
}

// GetDonations retrieves donations based on the provided parameters.
func (s *DonationService) GetDonations(ctx context.Context, id, campaignID, userID, page, limit string) ([]repo.TDonation, error) {
	pageNum, err := strconv.Atoi(page)
	if err != nil || page == "" {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limit == "" {
		limitNum = s.utilService.D().Limits.DefaultDonationLimit
	}

	donations, err := s.queries.GetDonations(ctx, repo.GetDonationsParams{
		ID:         s.utilService.ParseNullUUID(id),
		UserID:     s.utilService.ParseNullUUID(userID),
		CampaignID: s.utilService.ParseNullUUID(campaignID),
		Lim:        int32(limitNum),
		Off:        (int32(pageNum) - 1) * int32(limitNum),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrDonationNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringDonation, err)
	}

	return donations, nil
}

// GetDonationByID retrieves a donation by its ID.
func (s *DonationService) GetDonationByID(ctx context.Context, donationID string) (*repo.TDonation, error) {
	id, err := s.utilService.NParseUUID(donationID)
	if err != nil {
		return nil, err
	}

	donation, err := s.queries.GetDonationByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrDonationNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringDonation, err)
	}

	return &donation, nil
}

// CreateDonation creates a new donation record.
func (s *DonationService) CreateDonation(ctx context.Context, userID uuid.UUID, campaignID uuid.UUID, amount, transactionID string) (*uuid.UUID, error) {
	donationID, err := s.queries.CreateDonation(ctx, repo.CreateDonationParams{
		CampaignID:    campaignID,
		UserID:        userID,
		Amount:        amount,
		TransactionID: transactionID,
	})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingDontaions, err)
	}

	return &donationID, nil
}

// DeleteDonation deletes a donation by its ID.
func (s *DonationService) DeleteDonation(ctx context.Context, donationID string) error {
	id, err := s.utilService.NParseUUID(donationID)
	if err != nil {
		return err
	}

	if _, err := s.GetDonationByID(ctx, donationID); err != nil {
		return err
	}

	if err := s.queries.DeleteDonation(ctx, id); err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrDeletingDonations, err)
	}

	return nil
}

// CountDonations returns the total number of donations for a campaign or user.
func (s *DonationService) CountDonations(ctx context.Context, campaignID, userID string) (int64, error) {
	count, err := s.queries.CountDonations(ctx, repo.CountDonationsParams{
		CampaignID: s.utilService.ParseNullUUID(campaignID),
		UserID:     s.utilService.ParseNullUUID(userID),
	})
	if err != nil {
		return 0, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCountDonations, err)
	}

	return count, nil
}

// CheckIfUserHasDonated checks if a user has already donated to a campaign.
func (s *DonationService) CheckIfUserHasDonated(ctx context.Context, userID uuid.UUID, campaignID uuid.UUID) (bool, error) {
	count, err := s.CountDonations(ctx, campaignID.String(), userID.String())
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
