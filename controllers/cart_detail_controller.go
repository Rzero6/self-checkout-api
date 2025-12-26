package controllers

import (
	"strconv"

	"github.com/Rzero6/self-checkout-api/services"
	"github.com/Rzero6/self-checkout-api/utils"
	"github.com/gofiber/fiber/v2"
)

type AddCartDetailRequest struct {
	Barcode  string `json:"barcode"`
	Quantity int    `json:"quantity"`
}

func AddProductsToCart(c *fiber.Ctx) error {
	sessionID := c.Get(utils.SessionIDHeader)
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": utils.SessionIDMessage,
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
	if cartID == -1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized access",
		})
	}

	var req AddCartDetailRequest
	if err := c.BodyParser(&req); err != nil || req.Barcode == "" || req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Product barcode or quantity is invalid",
		})
	}

	cart, err := services.GetActiveCartBySession(sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Cart", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	if cart == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Cart", utils.MessageStatusFailed),
		})
	}

	item, err := services.GetProductByBarcode(req.Barcode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Product", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Product", utils.MessageStatusFailed),
		})
	}

	message, result, err := services.AddProductsToCart(*cart, *item, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataCreated("Detail", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    result,
	})
}
func GetCartDetailsBySessionID(c *fiber.Ctx) error {
	sessionID := c.Get(utils.SessionIDHeader)
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": utils.SessionIDMessage,
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
	if cartID == -1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized access",
		})
	}

	details, err := services.GetCartDetailsBySessionID(sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Details", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataRead("Details", utils.MessageStatusSuccess),
		"data":    details,
	})
}

type UpdateCartDetailRequest struct {
	ID       int64 `json:"id"`
	Quantity int   `json:"quantity"`
}

func UpdateDetailInCart(c *fiber.Ctx) error {
	sessionID := c.Get(utils.SessionIDHeader)
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": utils.SessionIDMessage,
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
	if cartID == -1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized update",
		})
	}

	var req UpdateCartDetailRequest
	if err := c.BodyParser(&req); err != nil || req.Quantity <= 0 || req.ID < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "detail id or quantity is invalid",
		})
	}

	detail, err := services.UpdateProductInCart(req.ID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataUpdated("Detail", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataUpdated("Detail", utils.MessageStatusSuccess),
		"data":    detail,
	})
}

func DeleteDetailFromCart(c *fiber.Ctx) error {
	sessionID := c.Get(utils.SessionIDHeader)
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": utils.SessionIDMessage,
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
	if cartID == -1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized delete",
		})
	}
	detailIDStr := c.Params("detail_id")
	detailID, err := strconv.ParseInt(detailIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid detail id",
		})
	}
	detail, err := services.DeleteProductFromCart(detailID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataDeleted("Detail", utils.MessageStatusSuccess),
		"data":    detail,
	})
}
