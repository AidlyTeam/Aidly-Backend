package sessionStore

import (
	"encoding/gob"

	repo "github.com/AidlyTeam/Aidly-Backend/internal/repos/out"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type SessionData struct {
	UserID        uuid.UUID
	RoleID        uuid.UUID
	Role          string
	Name          string
	Surname       string
	Email         string
	WalletAddress string
}

func (s *SessionData) ParseFromUser(user *repo.TUser, role *repo.TRole) {
	s.UserID = user.ID
	s.RoleID = user.RoleID
	s.Role = role.Name
	s.Name = user.Name.String
	s.Email = user.Email
	s.Surname = user.Surname.String
	s.WalletAddress = user.WalletAddress
}

func (s *SessionData) SetWalletAddress(walletAddress string) {
	s.WalletAddress = walletAddress
}

func (s *SessionData) SetNameSurname(name, surname string) {
	if name != "" {
		s.Name = name
	}
	if surname != "" {
		s.Surname = surname
	}
}

func (s *SessionData) SetRole(roleName string, roleID uuid.UUID) {
	s.Role = roleName
	s.RoleID = roleID
}

func GetSessionData(c *fiber.Ctx) *SessionData {
	user := c.Locals("user")
	if user == nil {
		return nil
	}
	sessionData, ok := user.(SessionData)
	if !ok {
		return nil
	}
	return &sessionData
}

func NewSessionStore(storage ...fiber.Storage) *session.Store {
	if len(storage) <= 0 {
		storage = append(storage, session.ConfigDefault.Storage)
	}
	gob.Register(SessionData{})
	return session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: true,
		Storage:        storage[0],
	})
}
