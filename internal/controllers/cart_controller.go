package controllers

import (
	"strconv"

	"github.com/Rzero6/self-checkout-api/internal/models"
	"github.com/Rzero6/self-checkout-api/internal/services"
	"github.com/Rzero6/self-checkout-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func StartCart(c *fiber.Ctx) error {
	sessionID := uuid.New()

	cart := models.Cart{
		SessionID: sessionID,
		Status:    string(utils.CartStatusActive),
	}

	result, err := services.CreateCart(cart)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataCreated("Cart", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataCreated("Cart", utils.MessageStatusSuccess),
		"data":    result,
	})
}

func GetCurrentCart(c *fiber.Ctx) error {
	sessionID := c.Get("X-Session-ID")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "session_id required",
		})
	}

	cart, err := services.GetActiveCartBySession(sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Cart", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	if cart == nil {
		return c.JSON(fiber.Map{
			"success": true,
			"message": utils.MessageDataRead("Cart", utils.MessageStatusSuccess),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataRead("Cart", utils.MessageStatusSuccess),
		"data":    cart,
	})
}

// Delete all product by deleting the cart :)
func DeleteAllProductFromCart(c *fiber.Ctx) error {
	sessionID := c.Get("X-Session-ID")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "session_id required",
		})
	}
	cartID, err := services.CheckCartExist(sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Cart", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	idStr := c.Params("cart_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid cart id",
			"error":   err.Error(),
		})
	}
	if cartID != id {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized access",
		})
	}

	if err := services.DeleteAllProductFromCart(cartID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataDeleted("Cart", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataDeleted("Cart", utils.MessageStatusSuccess),
	})
}
