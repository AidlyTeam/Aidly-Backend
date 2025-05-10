package services

import (
	"context"
	"database/sql"
	"strconv"

	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/google/uuid"
)

type BadgeService struct {
	db          *sql.DB
	queries     *repo.Queries
	utilService IUtilService
}

func newBadgeService(
	db *sql.DB,
	queries *repo.Queries,
	utilService IUtilService,
) *BadgeService {
	return &BadgeService{
		db:          db,
		queries:     queries,
		utilService: utilService,
	}
}

// GetBadges retrieves badges with optional ID filtering and pagination.
func (s *BadgeService) GetBadges(ctx context.Context, id, isNft, page, limit string) ([]repo.TBadge, error) {
	pageNum, err := strconv.Atoi(page)
	if err != nil || page == "" {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limit == "" {
		limitNum = s.utilService.D().Limits.DefaultBadgeLimit
	}

	var isNftBool bool
	var isNftSQL sql.NullBool

	if isNft != "" {
		isNftBool, err = strconv.ParseBool(isNft)
		if err != nil {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, "Invalid isVerified value")
		}
		isNftSQL = sql.NullBool{Bool: isNftBool, Valid: true}
	}

	badges, err := s.queries.GetBadges(ctx, repo.GetBadgesParams{
		ID:    s.utilService.ParseNullUUID(id),
		IsNft: isNftSQL,
		Off:   (int32(pageNum) - 1) * int32(limitNum),
		Lim:   int32(limitNum),
	})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringBadge, err)
	}

	return badges, nil
}

// GetBadgeByID fetches a badge by its ID.
func (s *BadgeService) GetBadgeByID(ctx context.Context, badgeID string) (*repo.TBadge, error) {
	id, err := s.utilService.NParseUUID(badgeID)
	if err != nil {
		return nil, err
	}

	badge, err := s.queries.GetBadgeByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrBadgeNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringBadge, err)
	}

	return &badge, nil
}

// CreateBadge inserts a new badge into the database.
func (s *BadgeService) CreateBadge(ctx context.Context, symbol, name, description, iconPath string, sellerFee int32, isNft bool, donationThreshold int32) (*uuid.UUID, error) {
	// Check if the donation threshold is already being used
	ok, err := s.queries.ExistsBadgeByThreshold(ctx, donationThreshold)
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringBadge, err)
	}
	if ok {
		return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrThresholdAlreadyBeingUsed)
	}

	// Create the new badge in the database
	badgeID, err := s.queries.CreateBadge(ctx, repo.CreateBadgeParams{
		Symbol:            s.utilService.ParseString(symbol),
		Name:              name,
		Description:       s.utilService.ParseString(description),
		IconPath:          s.utilService.ParseString(iconPath),
		SellerFee:         sql.NullInt32{Int32: sellerFee, Valid: sellerFee != 0},
		IsNft:             isNft,
		DonationThreshold: donationThreshold,
	})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingBadge, err)
	}

	return &badgeID, nil
}

// UpdateBadge updates an existing badge.
func (s *BadgeService) UpdateBadge(ctx context.Context, badgeID, symbol, name, description, iconPath, uri string, sellerFee int32, donationThreshold int32) error {
	// Parse the badgeID from string to UUID
	id, err := s.utilService.NParseUUID(badgeID)
	if err != nil {
		return err
	}

	// Check if the donation threshold is already being used
	ok, err := s.queries.ExistsBadgeByThreshold(ctx, donationThreshold)
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringBadge, err)
	}
	if ok {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrThresholdAlreadyBeingUsed)
	}

	// Update the badge in the database
	err = s.queries.UpdateBadge(ctx, repo.UpdateBadgeParams{
		BadgeID:           id,
		Symbol:            s.utilService.ParseString(symbol),
		Name:              s.utilService.ParseString(name),
		Description:       s.utilService.ParseString(description),
		IconPath:          s.utilService.ParseString(iconPath),
		Uri:               s.utilService.ParseString(uri),
		SellerFee:         sql.NullInt32{Int32: sellerFee, Valid: sellerFee != 0},
		DonationThreshold: sql.NullInt32{Int32: donationThreshold, Valid: donationThreshold != 0},
	})
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrUpdatingBadge, err)
	}

	return nil
}

// DeleteBadge deletes a badge by its ID.
func (s *BadgeService) DeleteBadge(ctx context.Context, badgeID string) error {
	id, err := s.utilService.NParseUUID(badgeID)
	if err != nil {
		return err
	}

	if _, err := s.GetBadgeByID(ctx, badgeID); err != nil {
		return err
	}

	if err := s.queries.DeleteBadge(ctx, id); err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrDeletingBadge, err)
	}

	return nil
}

// Add badge to a user
func (s *BadgeService) AddUserBadge(ctx context.Context, userID, badgeID uuid.UUID) (uuid.UUID, error) {
	exists, err := s.queries.GetUserBadgeExists(ctx, repo.GetUserBadgeExistsParams{
		UserID:  userID,
		BadgeID: badgeID,
	})
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrUserBadgeAlreadyExists)
	}

	badgeID, err = s.queries.AddUserBadge(ctx, repo.AddUserBadgeParams{
		UserID:  userID,
		BadgeID: badgeID,
	})
	if err != nil {
		return uuid.Nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingUserBadge, err)
	}

	return badgeID, nil
}

func (s *BadgeService) GetBadgeCount(ctx context.Context, id string) (int64, error) {
	count, err := s.queries.CountBadge(ctx, s.utilService.ParseNullUUID(id))
	if err != nil {
		return 0, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCountBadge, err)
	}

	return count, nil
}

func (s *BadgeService) GetUserBadgeCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := s.queries.CountUserBadge(ctx, userID)
	if err != nil {
		return 0, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCountBadge, err)
	}

	return count, nil
}

// Check if a user already has a badge
func (s *BadgeService) UserBadgeExists(ctx context.Context, userID, badgeID uuid.UUID) (bool, error) {
	return s.queries.GetUserBadgeExists(ctx, repo.GetUserBadgeExistsParams{
		UserID:  userID,
		BadgeID: badgeID,
	})
}

// Get all badges a user has
func (s *BadgeService) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]repo.TBadge, error) {
	return s.queries.GetUserBadges(ctx, userID)
}

// Remove a badge from a user
func (s *BadgeService) RemoveUserBadge(ctx context.Context, userID, badgeID uuid.UUID) error {
	return s.queries.RemoveUserBadge(ctx, repo.RemoveUserBadgeParams{
		UserID:  userID,
		BadgeID: badgeID,
	})
}

func (s *BadgeService) CheckBadgeAndAdd(ctx context.Context, userID uuid.UUID, count int32) (*uuid.UUID, error) {
	var id uuid.UUID

	ok, err := s.queries.ExistsBadgeByThreshold(ctx, count)
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringBadge, err)
	}
	if ok {
		badge, err := s.queries.GetBadgeByDonationCount(ctx, count)
		if err != nil {
			return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringBadge, err)
		}

		id, err = s.AddUserBadge(ctx, userID, badge.ID)
		if err != nil {
			return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingUserBadge, err)
		}
	}

	return &id, nil
}
