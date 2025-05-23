package services

import (
	"database/sql"

	"github.com/AidlyTeam/Aidly-Backend/internal/config"
	"github.com/AidlyTeam/Aidly-Backend/internal/config/models"
	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/google/uuid"
)

type IService interface {
	UtilService() IUtilService
	UploadService() *uploadService
	UserService() *UserService
	RoleService() *RoleService
	CampaignService() *CampaignService
	DonationService() *DonationService
	CategoryService() *CategoryService
	BadgeService() *BadgeService
}

type Services struct {
	utilService     IUtilService
	uploadService   *uploadService
	userService     *UserService
	roleService     *RoleService
	campaignService *CampaignService
	donationService *DonationService
	categoryService *CategoryService
	badgeService    *BadgeService
}

func CreateNewServices(
	validatorService IValidatorService,
	queries *repo.Queries,
	db *sql.DB,
	cfg *config.Config,
) *Services {
	utilService := newUtilService(validatorService, &cfg.Defaults)
	uploadService := newUploadService(utilService)
	userService := newUserService(db, queries, utilService)
	roleService := newRoleService(db, queries, utilService)
	campaignService := newCampaignService(db, queries, utilService)
	donationService := newDonationService(db, queries, utilService)
	categoryService := newCategoryService(db, queries, utilService)
	badgeService := newBadgeService(db, queries, utilService)

	return &Services{
		utilService:     utilService,
		userService:     userService,
		roleService:     roleService,
		campaignService: campaignService,
		uploadService:   uploadService,
		donationService: donationService,
		categoryService: categoryService,
		badgeService:    badgeService,
	}
}

func (s *Services) UtilService() IUtilService {
	return s.utilService
}

func (s *Services) UploadService() *uploadService {
	return s.uploadService
}

func (s *Services) UserService() *UserService {
	return s.userService
}

func (s *Services) RoleService() *RoleService {
	return s.roleService
}

func (s *Services) CampaignService() *CampaignService {
	return s.campaignService
}

func (s *Services) DonationService() *DonationService {
	return s.donationService
}

func (s *Services) CategoryService() *CategoryService {
	return s.categoryService
}

func (s *Services) BadgeService() *BadgeService {
	return s.badgeService
}

// ------------------------------------------------------

type IValidatorService interface {
	ValidateStruct(s any) error
}

type IUtilService interface {
	Validator() IValidatorService
	D() *models.Defaults
	ParseUUID(id string) (uuid.UUID, error)  // ID can be null
	NParseUUID(id string) (uuid.UUID, error) // ID cannot be null
	ParseString(str string) sql.NullString
	ParseNullUUID(str string) uuid.NullUUID
}

// -------------------

type utilService struct {
	validatorService IValidatorService
	defaults         *models.Defaults
}

func newUtilService(
	validatorService IValidatorService,
	defaults *models.Defaults,
) IUtilService {
	return &utilService{
		validatorService: validatorService,
		defaults:         defaults,
	}
}

func (s *utilService) Validator() IValidatorService {
	return s.validatorService
}

func (s *utilService) D() *models.Defaults {
	return s.defaults
}

func (s *utilService) ParseUUID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.UUID{}, nil
	}
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, serviceErrors.NewServiceErrorWithMessageAndError(
			serviceErrors.StatusBadRequest,
			serviceErrors.ErrInvalidID,
			err,
		)
	}
	return parsedUUID, nil
}

func (s *utilService) NParseUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, serviceErrors.NewServiceErrorWithMessageAndError(
			serviceErrors.StatusBadRequest,
			serviceErrors.ErrInvalidID,
			err,
		)
	}
	return parsedUUID, nil
}

func (s *utilService) ParseString(str string) sql.NullString {
	var value string
	var valid bool

	if str == "" {
		value = ""
		valid = false
	} else {
		value = str
		valid = true
	}

	return sql.NullString{String: value, Valid: valid}
}

func (s *utilService) ParseNullUUID(str string) uuid.NullUUID {
	var value uuid.UUID
	var valid bool

	if str == "" {
		valid = false
	} else {
		parsedUUID, err := uuid.Parse(str)
		if err != nil {
			valid = false
		} else {
			value = parsedUUID
			valid = true
		}
	}

	return uuid.NullUUID{UUID: value, Valid: valid}
}
