package middlewares

import (
	"errors"
	"strings"

	"github.com/AidlyTeam/Aidly-Backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func InitMiddlewares(cfg *config.Config) (mws []func(*fiber.Ctx) error) {
	cors := cors.New(
		cors.Config{
			AllowOrigins:     strings.Join(cfg.HTTP.AllowedOrigins, ","),
			AllowMethods:     strings.Join(cfg.HTTP.AllowedMethods, ","),
			AllowHeaders:     strings.Join(cfg.HTTP.AllowedHeaders, ","),
			AllowCredentials: cfg.HTTP.AllowCredentials,
			ExposeHeaders:    strings.Join(cfg.HTTP.ExposedHeaders, ","),
		},
	)
	helmetMid := helmet.New(helmet.ConfigDefault)

	mws = append(mws, cors, helmetMid)

	if !cfg.Application.DevMode {
		myLimiter := limiter.New(limiter.Config{
			Max: 50,
			LimitReached: func(c *fiber.Ctx) error {
				return errors.New("TOO MANY")
			},
		})
		mws = append(mws, myLimiter)
	}

	return
}
