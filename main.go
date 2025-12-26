package main

import (
	"context"
	"errors"
	"log"

	"github.com/Rzero6/self-checkout-api/config"
	"github.com/Rzero6/self-checkout-api/routes"
	"github.com/Rzero6/self-checkout-api/services"
	"github.com/gofiber/fiber/v2"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if errors.Is(err, context.Canceled) {
				return nil //handle nginx closed conn
			}
			return fiber.DefaultErrorHandler(c, err)
		},
	})
	app.Use(fiberRecover.New()) //recover middleware

	config.SetupCORS(app)

	routes.SetupRoutes(app, services.NewPaymentService())

	log.Fatal(app.Listen(":3000"))
}
