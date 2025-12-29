package controllers

import (
	"log"
	"net/mail"
	"os"
	"strconv"

	"github.com/Rzero6/self-checkout-api/services"
	"github.com/Rzero6/self-checkout-api/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceRequest struct {
	OrderID string `json:"order_id"`
	Email   string `json:"email"`
}

func CreateInvoice(ctx *fiber.Ctx) error {
	var req InvoiceRequest
	if err := ctx.BodyParser(&req); err != nil || req.OrderID == "" || req.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Order ID or email is invalid",
		})
	}
	_, mErr := mail.ParseAddress(req.Email)
	if mErr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email is invalid",
			"error":   mErr.Error(),
		})
	}

	transaction, err := services.GetTransactionByOrderID(req.OrderID, utils.TransactionStatusSuccess)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Transaction", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	var invoiceItems []utils.InvoiceItem
	for _, detail := range transaction.Details {
		item := utils.InvoiceItem{
			Name:     detail.ProductName,
			Price:    utils.FormatPrice(detail.Price),
			Qty:      strconv.FormatInt(int64(detail.Quantity), 10),
			Subtotal: utils.FormatPrice(detail.Subtotal),
		}
		invoiceItems = append(invoiceItems, item)
	}

	invoice := utils.InvoiceData{
		OrderID:     transaction.OrderID,
		Status:      transaction.Status,
		PaymentType: transaction.PaymentType,
		Total:       utils.FormatPrice(transaction.GrossAmount),
		Date:        transaction.CreatedAt.Format("02 Jan 2006"),
		Items:       invoiceItems,
	}

	pdf, err := utils.GeneratePDF(invoice)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"success": false,
				"message": "Failed to create PDF",
				"error":   err.Error(),
			},
		)
	}

	defer func() {
		if err := os.Remove(pdf); err != nil {
			log.Println("failed to delete pdf:", err)
		}
	}()

	attachments := []string{
		pdf,
	}
	if err := utils.SendEmail(
		req.Email,
		"Your Invoice "+transaction.OrderID,
		"Thank you for your purchase!",
		attachments); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to send invoice to " + req.Email,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Invoice has been sent to " + req.Email,
		"data":    true,
	})
}
