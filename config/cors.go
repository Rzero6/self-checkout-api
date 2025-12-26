package config

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupCORS(app *fiber.App) {
	env := os.Getenv("APP_ENV")
	frontUrl := os.Getenv("FRONT_APP_URL")

	if env == "production" && frontUrl != "" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     frontUrl,
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Session-ID",
			AllowCredentials: true,
		}))
	} else {
		app.Use(cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				if origin == "" {
					return true
				}
				return origin == "http://localhost:5173" ||
					strings.HasSuffix(origin, ".ngrok-free.app")
			},
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Session-ID, ngrok-skip-browser-warning",
			AllowCredentials: true,
		}))
	}
}
