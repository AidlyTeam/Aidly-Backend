package services

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/AidlyTeam/Aidly-Backend/internal/domains"
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

func (s *CampaignService) GetCampaigns(ctx context.Context, id, userID, IsVerified, statusType, title, page, limit string) ([]domains.Campaign, error) {
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

	campaigns, err := s.queries.GetCampaigns(ctx, repo.GetCampaignsParams{
		ID:         s.utilService.ParseNullUUID(id),
		UserID:     s.utilService.ParseNullUUID(userID),
		Title:      s.utilService.ParseString(title),
		IsVerified: isVerifiedSQL,
		StatusType: s.utilService.ParseString(statusType),
		Lim:        int32(limitNum),
		Off:        (int32(pageNum) - 1) * int32(limitNum),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCampaigns, err)
	}

	// Domains Mapper
	var domainsCampaigns []domains.Campaign
	for _, campaign := range campaigns {
		categories, err := s.queries.GetCampaignCategories(ctx, repo.GetCampaignCategoriesParams{CampaignID: campaign.ID, Lim: 100, Off: 0})
		if err != nil {
			return nil, serviceErrors.NewServiceErrorWithMessageAndError(
				serviceErrors.StatusInternalServerError,
				serviceErrors.ErrFilteringCampaignCategories,
				err,
			)
		}

		appCategories := domains.ToCategoriesCampaign(categories)
		domainsCampaigns = append(domainsCampaigns, *domains.ToCampaign(&campaign, appCategories))
	}

	return domainsCampaigns, nil
}

func (s *CampaignService) GetCampaignByID(ctx context.Context, campaignID string) (*domains.Campaign, error) {
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

	categories, err := s.queries.GetCampaignCategories(ctx, repo.GetCampaignCategoriesParams{CampaignID: campaign.ID, Lim: 100, Off: 0})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCampaignCategories, err)
	}

	appCategories := domains.ToCategoriesCampaign(categories)
	domainsCampaigns := domains.ToCampaignByID(&campaign, appCategories)

	return domainsCampaigns, nil
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
	title, description, walletAddress, imagePath, targetAmount, statusType, acceptedTokenSymbol string,
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

	validTokenSymbol := map[string]bool{
		"SOL":  true,
		"ZBTC": true,
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

	if _, ok := validTokenSymbol[acceptedTokenSymbol]; !ok {
		return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrInvalidTokenSymbol)
	}

	campaignID, err := s.queries.CreateCampaign(ctx, repo.CreateCampaignParams{
		UserID:              userID,
		Title:               title,
		Description:         s.utilService.ParseString(description),
		WalletAddress:       walletAddress,
		ImagePath:           s.utilService.ParseString(imagePath),
		TargetAmount:        targetAmount,
		AcceptedTokenSymbol: acceptedTokenSymbol,
		StatusType:          statusType,
		StartDate:           sql.NullTime{Time: startDateTime, Valid: !startDateTime.IsZero()},
		EndDate:             sql.NullTime{Time: endDateTime, Valid: !endDateTime.IsZero()},
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

	if _, err := s.GetCampaignByID(ctx, campaignID); err != nil {
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
	id, title, description, walletAddress, imagePath, targetAmount, statusType, acceptedTokenSymbol string,
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

	validTokenSymbol := map[string]bool{
		"SOL":  true,
		"ZBTC": true,
	}

	if _, ok := validTokenSymbol[acceptedTokenSymbol]; !ok {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrInvalidTokenSymbol)
	}

	if _, ok := validStatuses[statusType]; !ok {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrInvalidCampaignStatus)
	}

	if err := s.queries.UpdateCampaign(ctx, repo.UpdateCampaignParams{
		CampaignID:          idUUID,
		WalletAddress:       s.utilService.ParseString(walletAddress),
		Title:               s.utilService.ParseString(title),
		Description:         s.utilService.ParseString(description),
		ImagePath:           s.utilService.ParseString(imagePath),
		TargetAmount:        s.utilService.ParseString(targetAmount),
		StatusType:          s.utilService.ParseString(statusType),
		AcceptedTokenSymbol: s.utilService.ParseString(acceptedTokenSymbol),
		StartDate:           sql.NullTime{Time: startDateTime, Valid: !startDateTime.IsZero()},
		EndDate:             sql.NullTime{Time: endDateTime, Valid: !endDateTime.IsZero()},
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

	if campaign.EndDate != nil && campaign.EndDate.Before(time.Now()) {
		return false, nil
	}

	// Check if the target amount has been raised
	raisedAmount, err := decimal.NewFromString(campaign.RaisedAmount)
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

	category, err := s.queries.GetCampaignCategoriesOne(ctx, repo.GetCampaignCategoriesOneParams{
		CampaignID: campaign.ID,
		CategoryID: categoryID,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCampaignCategories, err)
	}
	if err == nil && category.CategoryID != uuid.Nil {
		return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrCategoryAlreadyAdded)
	}

	id, err := s.queries.CreateCampaignCategory(ctx, repo.CreateCampaignCategoryParams{
		CampaignID: campaign.ID,
		CategoryID: categoryID,
	})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingCampaingCategory, err)
	}

	return &id, nil
}

func (s *CampaignService) RemoveCategory(ctx context.Context, campaignID string, categoryID uuid.UUID) error {
	campaign, err := s.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return err
	}

	if _, err := s.GetCampaignCategoryByIDs(ctx, campaignID, categoryID); err != nil {
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

func (s *CampaignService) GetCampaignCategoriesByID(ctx context.Context, campaignID, page, limit string) ([]repo.GetCampaignCategoriesRow, error) {
	campaign, err := s.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	pageNum, err := strconv.Atoi(page)
	if err != nil || page == "" {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limit == "" {
		limitNum = s.utilService.D().Limits.DefaultCampaignCategoryLimit
	}

	offset := (pageNum - 1) * limitNum
	categories, err := s.queries.GetCampaignCategories(ctx, repo.GetCampaignCategoriesParams{
		CampaignID: campaign.ID,
		Lim:        int32(limitNum),
		Off:        int32(offset),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCampaignCategories, err)
	}

	return categories, nil
}

func (s *CampaignService) GetCampaignCategoryByIDs(ctx context.Context, campaignID string, categoryID uuid.UUID) (*repo.TCampaignCategory, error) {
	campaign, err := s.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	categories, err := s.queries.GetCampaignCategoriesOne(ctx, repo.GetCampaignCategoriesOneParams{
		CampaignID: campaign.ID,
		CategoryID: categoryID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCampaignCategoryNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCampaignCategories, err)
	}

	return &categories, nil
}
