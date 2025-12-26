package controllers

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Rzero6/self-checkout-api/models"
	"github.com/Rzero6/self-checkout-api/services"
	"github.com/Rzero6/self-checkout-api/utils"
	"github.com/gofiber/fiber/v2"
)

func GetAllProducts(c *fiber.Ctx) error {
	page, pageErr := strconv.Atoi(c.Query("page", "1"))
	limit, limitErr := strconv.Atoi(c.Query("limit", "10"))
	if pageErr != nil || limitErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid parameter page or limit",
		})
	}

	items, totalItems, err := services.GetAllProducts(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Product", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	var next, prev interface{} = nil, nil

	if page < totalPages {
		next = fmt.Sprintf("/api/products?page=%d&limit=%d", page+1, limit)
	}
	if page > 1 {
		prev = fmt.Sprintf("/api/products?page=%d&limit=%d", page-1, limit)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataRead("Products", utils.MessageStatusSuccess),
		"data":    items,
		"meta": fiber.Map{
			"page":        page,
			"per_page":    limit,
			"total_items": totalItems,
			"total_pages": totalPages,
		},
		"links": fiber.Map{
			"next": next,
			"prev": prev,
		},
	})
}

func GetProductByBarcode(c *fiber.Ctx) error {
	barcode := c.Query("barcode")
	if barcode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Parameter barcode is required",
		})
	}

	item, err := services.GetProductByBarcode(barcode)
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
			"error":   "Product not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataRead("Product", utils.MessageStatusSuccess),
		"data":    item,
	})
}

func GetProductRandom(c *fiber.Ctx) error {
	item, err := services.GetProductRandom()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Product", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataRead("Product", utils.MessageStatusSuccess),
		"data":    item,
	})
}

func CreateProduct(c *fiber.Ctx) error {
	item := new(models.Product)

	if err := c.BodyParser(item); err != nil || item.Name == "" || item.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Product name or price is invalid",
		})
	}

	barcode := utils.Generate13DigitID()
	item.Barcode = barcode

	result, err := services.CreateProduct(*item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataCreated("Product", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataCreated("Product", utils.MessageStatusSuccess),
		"data":    result,
	})
}
func GetDonationsProduct(c *fiber.Ctx) error {
	items, err := services.GetDonationsProduct()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Products", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataRead("Products", utils.MessageStatusSuccess),
		"data":    items,
	})
}
