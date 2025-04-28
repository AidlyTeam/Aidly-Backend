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

func (s *CampaignService) GetCampaignByID(ctx context.Context, campaignID uuid.UUID) (*repo.TCampaign, error) {
	campaign, err := s.queries.GetCampaignByID(ctx, campaignID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCampaigns, err)
	}

	return &campaign, nil
}

func (s *CampaignService) CreateCampaign(
	ctx context.Context,
	userID uuid.UUID,
	title, description, imagePath, targetAmount string,
	startDate, endDate time.Time,
) (*uuid.UUID, error) {
	campaignID, err := s.queries.CreateCampaign(ctx, repo.CreateCampaignParams{
		UserID:       userID,
		Title:        title,
		Description:  s.utilService.ParseString(description),
		ImagePath:    s.utilService.ParseString(imagePath),
		TargetAmount: targetAmount,
		StartDate:    sql.NullTime{Time: startDate, Valid: !startDate.IsZero()},
		EndDate:      sql.NullTime{Time: endDate, Valid: !endDate.IsZero()},
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
	id, title, description, imagePath, targetAmount string,
	startDate, endDate time.Time,
) error {
	idUUID, err := s.utilService.NParseUUID(id)
	if err != nil {
		return err
	}

	if err := s.queries.UpdateCampaign(ctx, repo.UpdateCampaignParams{
		CampaignID:   idUUID,
		Title:        title,
		Description:  s.utilService.ParseString(description),
		ImagePath:    s.utilService.ParseString(imagePath),
		TargetAmount: targetAmount,
		StartDate:    sql.NullTime{Time: startDate, Valid: !startDate.IsZero()},
		EndDate:      sql.NullTime{Time: endDate, Valid: !endDate.IsZero()},
	}); err != nil {
		if err == sql.ErrNoRows {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrUpdatingCampaigns, err)
	}
	return nil
}
