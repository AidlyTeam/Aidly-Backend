package services

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/shopspring/decimal"

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

func (s *CampaignService) GetCampaigns(ctx context.Context, id, userID, IsVerified, page, limit string) ([]repo.TCampaign, error) {
	pageNum, err := strconv.Atoi(page)
	if err != nil || page == "" {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limit == "" {
		limitNum = s.utilService.D().Limits.DefaultCampaignLimit
	}

	var isVerifiedBool bool
	var isVerifiedSQL sql.NullBool

	if IsVerified != "" {
		isVerifiedBool, err = strconv.ParseBool(IsVerified)
		if err != nil {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, "Invalid isVerified value")
		}
		isVerifiedSQL = sql.NullBool{Bool: isVerifiedBool, Valid: true}
	}

	campaign, err := s.queries.GetCampaigns(ctx, repo.GetCampaignsParams{
		ID:         s.utilService.ParseNullUUID(id),
		UserID:     s.utilService.ParseNullUUID(userID),
		IsVerified: isVerifiedSQL,
		Lim:        int32(limitNum),
		Off:        (int32(pageNum) - 1) * int32(limitNum),
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
	title, description, walletAddress, imagePath, targetAmount, statusType string,
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

	validStatuses := map[string]bool{
		"normal":    true,
		"urgent":    true,
		"critical":  true,
		"featured":  true,
		"scheduled": true,
	}

	if _, ok := validStatuses[statusType]; !ok {
		return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrInvalidCampaignStatus)
	}

	campaignID, err := s.queries.CreateCampaign(ctx, repo.CreateCampaignParams{
		UserID:        userID,
		Title:         title,
		Description:   s.utilService.ParseString(description),
		WalletAddress: walletAddress,
		ImagePath:     s.utilService.ParseString(imagePath),
		TargetAmount:  targetAmount,
		StatusType:    statusType,
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
	id, title, description, walletAddress, imagePath, targetAmount, statusType string,
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

	validStatuses := map[string]bool{
		"normal":    true,
		"urgent":    true,
		"critical":  true,
		"featured":  true,
		"scheduled": true,
	}

	if _, ok := validStatuses[statusType]; !ok {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrInvalidCampaignStatus)
	}

	if err := s.queries.UpdateCampaign(ctx, repo.UpdateCampaignParams{
		CampaignID:    idUUID,
		WalletAddress: s.utilService.ParseString(walletAddress),
		Title:         s.utilService.ParseString(title),
		Description:   s.utilService.ParseString(description),
		ImagePath:     s.utilService.ParseString(imagePath),
		TargetAmount:  s.utilService.ParseString(targetAmount),
		StatusType:    s.utilService.ParseString(statusType),
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

func (s *CampaignService) ChangeCampaignVerified(
	ctx context.Context,
	id string, verify bool,
) error {
	idUUID, err := s.utilService.NParseUUID(id)
	if err != nil {
		return err
	}

	if err := s.queries.ChangeVerified(ctx, repo.ChangeVerifiedParams{
		CampaignID: idUUID,
		IsVerified: verify,
	}); err != nil {
		if err == sql.ErrNoRows {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrUpdatingCampaigns, err)
	}
	return nil
}

func (s *CampaignService) CheckCampaignValidity(ctx context.Context, campaignID string) (bool, error) {
	campaign, err := s.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return false, err
	}

	if campaign.EndDate.Valid && campaign.EndDate.Time.Before(time.Now()) {
		return false, nil
	}

	// Check if the target amount has been raised
	raisedAmount, err := decimal.NewFromString(campaign.RaisedAmount.String)
	if err != nil {
		return false, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrDecimalConvertionError, err)
	}
	targetAmount, err := decimal.NewFromString(campaign.TargetAmount)
	if err != nil {
		return false, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrDecimalConvertionError, err)
	}

	if targetAmount.LessThanOrEqual(raisedAmount) {
		return false, nil
	}

	return true, nil
}

// UpdateCampaignValidity checks if the campaign is valid and updates the validity status.
func (s *CampaignService) UpdateCampaignValidity(ctx context.Context, campaignID string) error {
	// Step 1: Check the current validity status of the campaign
	isValid, err := s.CheckCampaignValidity(ctx, campaignID)
	if err != nil {
		return err
	}

	// Step 2: Update the campaign's validity status in the database
	if err := s.queries.ChangeValid(ctx, repo.ChangeValidParams{
		CampaignID: uuid.MustParse(campaignID),
		IsValid:    isValid,
	}); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) AddCategory(ctx context.Context, campaignID string, categoryID uuid.UUID) (*uuid.UUID, error) {
	campaign, err := s.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	id, err := s.queries.CreateCampaignCategory(ctx, repo.CreateCampaignCategoryParams{
		CampaignID: campaign.ID,
		CategoryID: categoryID,
	})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreateCampaingCategory, err)
	}

	return &id, nil
}

func (s *CampaignService) RemoveCategory(ctx context.Context, campaignID string, categoryID uuid.UUID) error {
	campaign, err := s.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return err
	}

	if err := s.queries.DeleteCampaignCategory(ctx, repo.DeleteCampaignCategoryParams{
		CampaignID: campaign.ID,
		CategoryID: categoryID,
	}); err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrDeletingCampaignCategory, err)
	}

	return nil
}
