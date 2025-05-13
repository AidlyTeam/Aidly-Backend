package private

import (
	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/sessionStore"
	"github.com/gofiber/fiber/v2"
)

func (h *PrivateHandler) initUserRoutes(root fiber.Router) {
	user := root.Group("/user")
	user.Get("/profile", h.Profile)
	user.Post("/profile", h.UpdateProfile)
	user.Post("/connect", h.ConnectWallet)
}

// @Tags User
// @Summary Get User Profile
// @Description Retrieves users profile.
// @Accept json
// @Produce json
// @Success 200 {object} response.BaseResponse{}
// @Router /private/user/profile [get]
func (h *PrivateHandler) Profile(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)

	badges, err := h.services.BadgeService().GetUserBadges(c.Context(), userSession.UserID)
	if err != nil {
		return err
	}
	badgeCount, err := h.services.BadgeService().GetUserBadgeCount(c.Context(), userSession.UserID)
	if err != nil {
		return err
	}
	//count badge
	userProfileDTO := h.dtoManager.UserManager().ToUserProfile(*userSession, badges, badgeCount)

	return response.Response(200, "Status OK", userProfileDTO)
}

// @Tags User
// @Summary Update User Profile
// @Description Updates users profile.
// @Accept json
// @Produce json
// @Param newUserProfile body dto.UserProfileUpdateDTO true "New User Profile"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/user/profile [post]
func (h *PrivateHandler) UpdateProfile(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)

	var newUserProfile dto.UserProfileUpdateDTO
	if err := c.BodyParser(&newUserProfile); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(newUserProfile); err != nil {
		return err
	}

	// Update Profile
	if err := h.services.UserService().Update(c.Context(), userSession.UserID.String(), newUserProfile.Name, newUserProfile.Surname, newUserProfile.Email); err != nil {
		return err
	}

	// If a user with the 'first' role updates their profile, their role will be changed to 'user'.
	if userSession.Role == h.config.Defaults.Roles.FirstRole {
		defaultRole, err := h.services.RoleService().GetDefault(c.Context())
		if err != nil {
			return err
		}

		if err := h.services.UserService().ChangeUserRole(c.Context(), userSession.UserID, defaultRole.ID); err != nil {
			return err
		}

		userSession.SetRole(defaultRole.Name, defaultRole.ID)
	}

	sess, err := h.sess_store.Get(c)
	if err != nil {
		return err
	}
	userSession.SetNameSurname(newUserProfile.Name, newUserProfile.Surname)
	sess.Set("user", userSession)
	if err := sess.Save(); err != nil {
		return err
	}

	return response.Response(200, "Status OK", nil)
}

// @Tags User
// @Summary Connect Wallet
// @Description Connect Phantom Wallet.
// @Accept json
// @Produce json
// @Param wallet body dto.UserConnectWallectDTO true "Wallet"
// @Success 200 {object} response.BaseResponse{}
// @Router /private/user/connect [post]
func (h *PrivateHandler) ConnectWallet(c *fiber.Ctx) error {
	userSession := sessionStore.GetSessionData(c)

	var wallet dto.UserConnectWallectDTO
	if err := c.BodyParser(&wallet); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(wallet); err != nil {
		return err
	}

	if err := h.services.UserService().ConnectWallet(c.Context(), userSession.UserID.String(), wallet.WalletAddress, wallet.Message, wallet.Signature); err != nil {
		return err
	}

	sess, err := h.sess_store.Get(c)
	if err != nil {
		return err
	}
	userSession.SetWalletAddress(userSession.WalletAddress)
	sess.Set("user", userSession)
	if err := sess.Save(); err != nil {
		return err
	}

	return response.Response(200, "Status OK", nil)
}
