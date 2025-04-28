package admin

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *AdminHandler) initUserRoutes(root fiber.Router) {
	user := root.Group("/user")
	fmt.Println("User Admin Routes", user)
}
