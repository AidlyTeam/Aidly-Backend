package services

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"

	"github.com/google/uuid"
)

type CampaignService struct {
	db          *sql.DB
	queries     *repo.Queries
	utilService IUtilService
}

func newCampaignService(
	db *sql.DB,
	queries *repo.Queries,
	utilService IUtilService,
) *CampaignService {
	return &CampaignService{
		db:          db,
		queries:     queries,
		utilService: utilService,
	}
}

func (s *CampaignService) GetCampaigns(ctx context.Context, id, userID, page, limit string) ([]repo.TCampaign, error) {
	pageNum, err := strconv.Atoi(page)
	if err != nil || page == "" {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limit == "" {
		limitNum = s.utilService.D().Limits.DefaultCampaignLimit
	}

	campaign, err := s.queries.GetCampaigns(ctx, repo.GetCampaignsParams{
		ID:     s.utilService.ParseNullUUID(id),
		UserID: s.utilService.ParseNullUUID(userID),
		Lim:    int32(limitNum),
		Off:    (int32(pageNum) - 1) * int32(limitNum),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCampaigns, err)
	}

	return campaign, nil
}

func (s *CampaignService) GetCampaignByID(ctx context.Context, campaignID string) (*repo.TCampaign, error) {
	id, err := s.utilService.NParseUUID(campaignID)
	if err != nil {
		return nil, err
	}

	campaign, err := s.queries.GetCampaignByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCampaigns, err)
	}

	return &campaign, nil
}

func (s *CampaignService) CheckTheOwnerOfCampaign(ctx context.Context, id string, userID uuid.UUID) error {
	idUUID, err := s.utilService.NParseUUID(id)
	if err != nil {
		return err
	}

	if _, err := s.queries.GetUserCampaign(ctx, repo.GetUserCampaignParams{
		CampaignID: idUUID,
		UserID:     userID,
	}); err != nil {
		if err == sql.ErrNoRows {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrDeletingCampaigns, err)
	}

	return nil
}

func (s *CampaignService) CreateCampaign(
	ctx context.Context,
	userID uuid.UUID,
	title, description, walletAddress, imagePath, targetAmount string,
	startDate, endDate string,
) (*uuid.UUID, error) {
	startDateTime, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, err
	}
	endDateTime, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, err
	}

	campaignID, err := s.queries.CreateCampaign(ctx, repo.CreateCampaignParams{
		UserID:        userID,
		Title:         title,
		Description:   s.utilService.ParseString(description),
		WalletAddress: walletAddress,
		ImagePath:     s.utilService.ParseString(imagePath),
		TargetAmount:  targetAmount,
		StartDate:     sql.NullTime{Time: startDateTime, Valid: !startDateTime.IsZero()},
		EndDate:       sql.NullTime{Time: endDateTime, Valid: !endDateTime.IsZero()},
	})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingCampaigns, err)
	}

	return &campaignID, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, campaignID string) error {
	id, err := s.utilService.NParseUUID(campaignID)
	if err != nil {
		return err
	}

	if err = s.queries.DeleteCampaign(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrDeletingCampaigns, err)
	}
	return nil
}

func (s *CampaignService) UpdateCampaign(
	ctx context.Context,
	userID uuid.UUID,
	id, title, description, walletAddress, imagePath, targetAmount string,
	startDate, endDate string,
) error {
	idUUID, err := s.utilService.NParseUUID(id)
	if err != nil {
		return err
	}

	// IDOR Safe
	if err := s.CheckTheOwnerOfCampaign(ctx, id, userID); err != nil {
		return err
	}

	var startDateTime, endDateTime time.Time
	if startDate != "" {
		startDateTime, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			return err
		}
	}
	if endDate != "" {
		endDateTime, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			return err
		}
	}

	if err := s.queries.UpdateCampaign(ctx, repo.UpdateCampaignParams{
		CampaignID:    idUUID,
		WalletAddress: walletAddress,
		Title:         title,
		Description:   s.utilService.ParseString(description),
		ImagePath:     s.utilService.ParseString(imagePath),
		TargetAmount:  s.utilService.ParseString(targetAmount),
		StartDate:     sql.NullTime{Time: startDateTime, Valid: !startDateTime.IsZero()},
		EndDate:       sql.NullTime{Time: endDateTime, Valid: !endDateTime.IsZero()},
	}); err != nil {
		if err == sql.ErrNoRows {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrUpdatingCampaigns, err)
	}
	return nil
}
