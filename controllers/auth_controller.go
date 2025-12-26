package controllers

import (
	"github.com/Rzero6/self-checkout-api/services"
	"github.com/Rzero6/self-checkout-api/utils"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil || req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "username and password are required",
		})
	}

	token, err := services.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageLogin(utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageLogin(utils.MessageStatusSuccess),
		"token":   token,
	})
}
