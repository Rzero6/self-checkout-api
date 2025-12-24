package main

import (
	"context"
	"errors"
	"log"

	"github.com/Rzero6/self-checkout-api/config"
	"github.com/Rzero6/self-checkout-api/internal/routes"
	"github.com/Rzero6/self-checkout-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app.Use(fiberRecover.New())   //recover middleware
	app.Use(cors.New(cors.Config{ //handle localhost request permission
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))
	routes.SetupRoutes(app, services.NewPaymentService())

	log.Fatal(app.Listen(":3000"))
}
