package services

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	hasherService "github.com/AidlyTeam/Aidly-Backend/pkg/hasher"
	"github.com/AidlyTeam/Aidly-Backend/pkg/paths"
	"github.com/google/uuid"
)

type UserService struct {
	db          *sql.DB
	queries     *repo.Queries
	utilService IUtilService
	userMutex   sync.Mutex
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

func (s *UserService) CivicLogin(ctx context.Context, fullName, email string, defaultRoleID uuid.UUID) (*repo.TUser, bool, error) {
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	users, err := s.queries.GetUsers(ctx, repo.GetUsersParams{
		Email: s.utilService.ParseString(email),
		Lim:   1,
		Off:   0,
	})
	if err != nil {
		return nil, false, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
	}
	if len(users) == 0 {
		name, surname := paths.SplitFullName(fullName)
		id, err := s.queries.CreateUser(ctx, repo.CreateUserParams{
			RoleID:  defaultRoleID,
			Email:   s.utilService.ParseString(email),
			Name:    s.utilService.ParseString(name),
			Surname: s.utilService.ParseString(surname),
		})
		if err != nil {
			return nil, false, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingUser, err)
		}

		user, err := s.queries.GetUserByID(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, false, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrUserNotFound)
			}
			return nil, false, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
		}

		return &user, true, nil
	}
	user := &users[0]

	return user, false, err
}

// Login & Register
func (s *UserService) AuthWallet(ctx context.Context, walletAddress, message, signature string, defaultRoleID uuid.UUID) (*repo.TUser, bool, error) {
	ok, err := hasherService.VerifySignature(walletAddress, message, signature)
	if err != nil {
		return nil, false, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrWalletVerificationError, err)
	}
	if !ok {
		return nil, false, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrInvalidWalletConnection)
	}

	users, err := s.queries.GetUsers(ctx, repo.GetUsersParams{
		WalletAddress: s.utilService.ParseString(walletAddress),
		Lim:           1,
		Off:           0,
	})
	if err != nil {
		return nil, false, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
	}
	if len(users) == 0 {
		id, err := s.queries.CreateUser(ctx, repo.CreateUserParams{
			RoleID:        defaultRoleID,
			WalletAddress: walletAddress,
		})
		if err != nil {
			return nil, false, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreatingUser, err)
		}

		user, err := s.queries.GetUserByID(ctx, id)
		if err != nil {
			if strings.Contains(err.Error(), "sql: no rows in result set") {
				return nil, false, serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrUserNotFound)
			}
			return nil, false, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
		}

		return &user, true, nil
	}
	user := &users[0]

	return user, false, err
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

func (s *UserService) Update(ctx context.Context, id, name, surname, email string) error {
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

	users, err := s.queries.GetUsers(ctx, repo.GetUsersParams{
		Email: s.utilService.ParseString(email),
		Lim:   1,
		Off:   0,
	})
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
	}
	if len(users) > 0 {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrEmailBeingUsed)
	}

	if err := s.queries.UpdateUser(ctx, repo.UpdateUserParams{
		UserID:  idUUID,
		Name:    s.utilService.ParseString(name),
		Surname: s.utilService.ParseString(surname),
		Email:   s.utilService.ParseString(email),
	}); err != nil {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrUpdatingUsers)
	}

	return nil
}

func (s *UserService) ConnectWallet(ctx context.Context, id, walletAddress, message, signature string) error {
	ok, err := hasherService.VerifySignature(walletAddress, message, signature)
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrWalletVerificationError, err)
	}
	if !ok {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrInvalidWalletConnection)
	}

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

	users, err := s.queries.GetUsers(ctx, repo.GetUsersParams{
		WalletAddress: s.utilService.ParseString(walletAddress),
		Lim:           1,
		Off:           0,
	})
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFilteringUsers, err)
	}
	if len(users) > 0 {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusBadRequest, serviceErrors.ErrWalletAddressBeingUsed)
	}

	if err := s.queries.ConnectWallet(ctx, repo.ConnectWalletParams{
		UserID:        idUUID,
		WalletAddress: walletAddress,
	}); err != nil {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrConnectingWallet)
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

func (s *UserService) Statistics(ctx context.Context) (*repo.StatisticsRow, error) {
	dbModel, err := s.queries.Statistics(ctx)
	if err != nil {
		return nil, serviceErrors.NewServiceErrorWithMessageAndError(serviceErrors.StatusInternalServerError, serviceErrors.ErrFetchingIstatistic, err)
	}

	return &dbModel, nil
}
