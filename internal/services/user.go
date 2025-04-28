package services

import (
	"context"
	"database/sql"
	"strings"

	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	hasherService "github.com/AidlyTeam/Aidly-Backend/pkg/hasher"
	"github.com/google/uuid"
)

type UserService struct {
	db          *sql.DB
	queries     *repo.Queries
	utilService IUtilService
}

func newUserService(
	db *sql.DB,
	queries *repo.Queries,
	utilService IUtilService,
) *UserService {
	return &UserService{
		db:          db,
		queries:     queries,
		utilService: utilService,
	}
}

func (s *UserService) Login(ctx context.Context, walletAddress string) (*repo.TUser, error) {
	user, err := s.queries.GetUserByWalletAddress(ctx, walletAddress)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrUserNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
	}
	return &user, nil
}

// Login & Register
func (s *UserService) AuthWallet(ctx context.Context, walletAddress, message, signature string, defaultRoleID uuid.UUID) (*repo.TUser, error) {
	ok, err := hasherService.VerifySignature(walletAddress, message, signature)
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrWalletVerificationError, err)
	}
	if !ok {
		return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrInvalidWalletConnection)
	}

	users, err := s.queries.GetUsers(ctx, repo.GetUsersParams{
		WalletAddress: s.utilService.ParseString(walletAddress),
		Lim:           1,
		Off:           0,
	})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
	}
	if len(users) == 0 {
		id, err := s.queries.CreateUser(ctx, repo.CreateUserParams{
			RoleID:        defaultRoleID,
			WalletAddress: walletAddress,
		})
		if err != nil {
			return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingUser, err)
		}

		user, err := s.queries.GetUserByID(ctx, id)
		if err != nil {
			if strings.Contains(err.Error(), "sql: no rows in result set") {
				return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrUserNotFound)
			}
			return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
		}

		return &user, nil
	}
	user := &users[0]

	return user, err
}

func (s *UserService) AdminCreate(ctx context.Context, walletAddress string, defaultRoleID uuid.UUID) (*repo.TUser, error) {
	id, err := s.queries.CreateUser(ctx, repo.CreateUserParams{
		RoleID:        defaultRoleID,
		WalletAddress: walletAddress,
		IsDefault:     true,
	})
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingUser, err)
	}

	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrUserNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
	}

	return &user, err
}

func (s *UserService) GetDefaultUser(ctx context.Context) (*repo.TUser, error) {
	user, err := s.queries.GetDefaultUser(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return nil, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrUserNotFound)
		}
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
	}

	return &user, nil
}

func (s *UserService) Update(ctx context.Context, id, name, surname string) error {
	idUUID, err := s.utilService.NParseUUID(id)
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrInvalidID)
	}

	if _, err := s.queries.GetUserByID(ctx, idUUID); err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrUserNotFound)
		}
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
	}

	if err := s.queries.UpdateUser(ctx, repo.UpdateUserParams{
		UserID:  idUUID,
		Name:    s.utilService.ParseString(surname),
		Surname: s.utilService.ParseString(surname),
	}); err != nil {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrUpdatingUsers)
	}

	return nil
}

func (s *UserService) ChangeUserRole(ctx context.Context, id, newRoleID uuid.UUID) error {
	if err := s.queries.ChangeUserRole(ctx, repo.ChangeUserRoleParams{
		UserID: id,
		RoleID: newRoleID,
	}); err != nil {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrUpdatingUserRole)
	}

	return nil
}
