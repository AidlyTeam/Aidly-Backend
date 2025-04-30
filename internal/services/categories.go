package services

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/AidlyTeam/Aidly-Backend/internal/domains"
	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/google/uuid"
)

type CategoryService struct {
	db          *sql.DB
	queries     *repo.Queries
	utilService IUtilService
}

func newCategoryService(
	db *sql.DB,
	queries *repo.Queries,
	utilService IUtilService,
) *CategoryService {
	return &CategoryService{
		db:          db,
		queries:     queries,
		utilService: utilService,
	}
}

// GetCategories returns a list of all categories.
func (s *CategoryService) GetCategories(ctx context.Context, page, limit string) ([]domains.Category, error) {
	pageNum, err := strconv.Atoi(page)
	if err != nil || page == "" {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limit == "" {
		limitNum = s.utilService.D().Limits.DefaultDonationLimit
	}

	categories, err := s.queries.GetCategories(ctx, repo.GetCategoriesParams{
		Lim: int32(limitNum),
		Off: (int32(pageNum) - 1) * int32(limitNum),
	})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCategories, err)
	}

	appModels := domains.ToCategories(categories)

	return appModels, nil
}

// GetCategoryByID returns a category by its ID.
func (s *CategoryService) GetCategoryByID(ctx context.Context, categoryID string) (*domains.Category, error) {
	id, err := s.utilService.NParseUUID(categoryID)
	if err != nil {
		return nil, err
	}

	category, err := s.queries.GetCategoryByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusNotFound, serviceErrors.ErrCategoryNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringCategories, err)
	}
	appCategory := domains.ToCategory(&category)

	return appCategory, nil
}

// CreateCategory creates a new category.
func (s *CategoryService) CreateCategory(ctx context.Context, name string) (*uuid.UUID, error) {
	id, err := s.queries.CreateCategory(ctx, name)
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingCategories, err)
	}
	return &id, nil
}

// UpdateCategory updates the name of a category.
func (s *CategoryService) UpdateCategory(ctx context.Context, idStr, name string) error {
	id, err := s.utilService.NParseUUID(idStr)
	if err != nil {
		return err
	}

	err = s.queries.UpdateCategory(ctx, repo.UpdateCategoryParams{
		ID:   id,
		Name: s.utilService.ParseString(name),
	})
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrUpdatingCategories, err)
	}

	return nil
}

// DeleteCategory deletes a category by ID.
func (s *CategoryService) DeleteCategory(ctx context.Context, categoryID string) error {
	id, err := s.utilService.NParseUUID(categoryID)
	if err != nil {
		return err
	}

	if _, err := s.GetCategoryByID(ctx, categoryID); err != nil {
		return err
	}

	err = s.queries.DeleteCategory(ctx, id)
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrDeletingCategories, err)
	}

	return nil
}

// CountDonations returns the total number of donations for a campaign or user.
func (s *CategoryService) CountCategories(ctx context.Context) (int64, error) {
	count, err := s.queries.CountCategory(ctx)
	if err != nil {
		return 0, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCountCategories, err)
	}

	return count, nil
}
