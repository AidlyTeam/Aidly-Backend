package dto

import (
	"github.com/AidlyTeam/Aidly-Backend/internal/http/sessionStore"
	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
)

type UserDTOManager struct{}

func NewUserDTOManager() UserDTOManager {
	return UserDTOManager{}
}

type UserNameSurnameView struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func (UserDTOManager) ToUserNameSurname(name, surname string) *UserNameSurnameView {
	return &UserNameSurnameView{
		Name:    name,
		Surname: surname,
	}
}

type UserProfileView struct {
	ID            string      `json:"id"`
	RoleID        string      `json:"roleID"`
	WalletAddress string      `json:"walletAddress"`
	RoleName      string      `json:"role"`
	Name          string      `json:"name"`
	Surname       string      `json:"surname"`
	Badges        *BadgeViews `json:"badges"`
}

func (UserDTOManager) ToUserProfile(userData sessionStore.SessionData, badges []repo.GetUserBadgesRow, badgeCount int64) UserProfileView {
	badgeManager := new(BadgeDTOManager)

	return UserProfileView{
		ID:            userData.UserID.String(),
		RoleID:        userData.RoleID.String(),
		RoleName:      userData.Role,
		Name:          userData.Name,
		Surname:       userData.Surname,
		WalletAddress: userData.WalletAddress,
		Badges:        badgeManager.ToUserBadgeViews(badges, badgeCount),
	}
}

type UserAuthWalletDTO struct {
	WalletAddress string `json:"walletAddress" validate:"required"`
	Message       string `json:"message" validate:"required"`
	Signature     string `json:"signatureBase58" validate:"required"`
}

type UserCivicAuthDTO struct {
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type UserProfileUpdateDTO struct {
	Name    string `json:"name" validate:"omitempty,max=30"`
	Surname string `json:"surname" validate:"omitempty,max=30"`
	Email   string `json:"email" validate:"omitempty,email"`
}
