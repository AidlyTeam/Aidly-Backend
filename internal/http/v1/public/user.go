package public

import (
	dto "github.com/AidlyTeam/Aidly-Backend/internal/http/dtos"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/response"
	"github.com/AidlyTeam/Aidly-Backend/internal/http/sessionStore"

	"github.com/gofiber/fiber/v2"
)

func (h *PublicHandler) initUserRoutes(root fiber.Router) {
	root.Post("/login", h.Login)
	root.Post("/logout", h.Logout)
	root.Post("/auth", h.Auth)
}

// FIXME: THIS IS FOR DEVELOPMENT ONLY  - DELETE IT LATER
// @Tags Auth
// @Summary Login
// @Description Login
// @Accept json
// @Produce json
// @Success 200 {object} response.BaseResponse{}
// @Router /public/login [post]
func (h *PublicHandler) Login(c *fiber.Ctx) error {
	users, err := h.services.UserService().GetDefaultUser(c.Context())
	if err != nil {
		return err
	}

	role, err := h.services.RoleService().GetRoleByID(c.Context(), users.RoleID)
	if err != nil {
		return err
	}

	sess, err := h.sessionStore.Get(c)
	if err != nil {
		return err
	}
	sessionData := sessionStore.SessionData{}
	sessionData.ParseFromUser(users, role)
	sess.Set("user", sessionData)
	if err := sess.Save(); err != nil {
		return err
	}
	profileResponse := h.dtoManager.UserManager().ToUserProfile(sessionData, nil, 0)

	return response.Response(200, "Login successful", profileResponse)
}

// @Tags Auth
// @Summary Auth
// @Description Auth with Wallet
// @Accept json
// @Produce json
// @Param auth body dto.UserAuthWalletDTO true "Auth Information"
// @Success 200 {object} response.BaseResponse{}
// @Router /public/auth [post]
func (h *PublicHandler) Auth(c *fiber.Ctx) error {
	var auth dto.UserAuthWalletDTO
	if err := c.BodyParser(&auth); err != nil {
		return err
	}
	if err := h.services.UtilService().Validator().ValidateStruct(auth); err != nil {
		return err
	}

	firstRole, err := h.services.RoleService().GetByName(c.Context(), h.config.Defaults.Roles.FirstRole)
	if err != nil {
		return err
	}

	user, isRegister, err := h.services.UserService().AuthWallet(
		c.Context(),
		auth.WalletAddress,
		auth.Message,
		auth.Signature,
		firstRole.ID,
	)
	if err != nil {
		return err
	}

	sess, err := h.sessionStore.Get(c)
	if err != nil {
		return err
	}

	sessionData := sessionStore.SessionData{}
	if !isRegister {
		// If user is not registering. Which means he is login in. Then set the users role in session else firstrole.
		userRole, err := h.services.RoleService().GetRoleByID(c.Context(), user.RoleID)
		if err != nil {
			return err
		}
		sessionData.ParseFromUser(user, userRole)
	} else {
		sessionData.ParseFromUser(user, firstRole)
	}
	sess.Set("user", sessionData)
	if err := sess.Save(); err != nil {
		return err
	}
	profileResponse := h.dtoManager.UserManager().ToUserProfile(sessionData, nil, 0)

	return response.Response(200, "Authenticate Successful", profileResponse)
}

// @Tags Auth
// @Summary Logout
// @Description Logout
// @Accept json
// @Produce json
// @Success 200 {object} response.BaseResponse{}
// @Router /public/logout [post]
func (h *PublicHandler) Logout(c *fiber.Ctx) error {
	session, err := h.sessionStore.Get(c)
	if err != nil {
		return response.Response(500, "Failed to get session", err)
	}
	if err := session.Destroy(); err != nil {
		return response.Response(500, "Failed to destroy session", err)
	}

	return response.Response(200, "Logout successful", nil)
}
