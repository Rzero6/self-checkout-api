package controllers

import (
	"strconv"
	"strings"

	"github.com/Rzero6/self-checkout-api/models"
	"github.com/Rzero6/self-checkout-api/services"
	"github.com/Rzero6/self-checkout-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go/coreapi"
)

type PaymentController struct {
	service *services.PaymentService
}

func NewPaymentController(service *services.PaymentService) *PaymentController {
	return &PaymentController{service: service}
}

type CreatePayment struct {
	PaymentType string `json:"payment_type"`
}

func (p *PaymentController) CreateTransaction(ctx *fiber.Ctx) error {
	sessionID := ctx.Get(utils.SessionIDHeader)
	if sessionID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": utils.SessionIDMessage,
		})
	}
	cartID, cErr := services.CheckCartExist(sessionID)
	if cErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Cart", utils.MessageStatusFailed),
			"error":   cErr.Error(),
		})
	}
	if cartID == -1 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Cart", utils.MessageStatusFailed),
		})
	}

	var req CreatePayment
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid input body",
			"error":   err.Error(),
		})
	}
	// default to QRIS
	req.PaymentType = "qris"
	//
	if req.PaymentType == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "payment_type are required",
		})
	}

	cartItem, ciErr := services.GetCartDetailsBySessionID(sessionID)
	if ciErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Details", utils.MessageStatusFailed),
			"error":   ciErr.Error(),
		})
	}
	if len(cartItem) <= 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "No product in the cart",
		})
	}
	//Check existing transaction
	existingTransaction, tErr := services.CheckExistingTransaction(int64(cartItem[0].CartID))
	if tErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Transaction", utils.MessageStatusFailed),
			"error":   tErr.Error(),
		})
	}
	if existingTransaction != nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "Payment Gateway already exist",
			"data":    existingTransaction,
		})
	}

	var grossAmt int64
	for _, item := range cartItem {
		grossAmt += item.Subtotal
	}

	orderID := utils.GenerateOrderID()
	tx := models.Transaction{
		CartID:      cartItem[0].CartID,
		OrderID:     orderID,
		GrossAmount: int64(grossAmt),
	}
	//Create Transaction
	_, transaction, err := p.service.CreateTransaction(tx, req.PaymentType)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataCreated("Transaction", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	//Create DetailsTransaction
	transactionDetails, tdErr := services.CreateTransactionDetails(int64(transaction.ID), cartItem)
	if tdErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataCreated("Transaction's Details", utils.MessageStatusFailed),
			"error":   tdErr.Error(),
		})
	}
	transaction.Details = transactionDetails
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataCreated("Transaction", utils.MessageStatusSuccess),
		// "payment": resp,
		"data": transaction,
	})
}

func (p *PaymentController) GetTransactionStatus(ctx *fiber.Ctx) error {
	sessionID := ctx.Get(utils.SessionIDHeader)
	if sessionID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": utils.SessionIDMessage,
		})
	}
	if _, err := uuid.Parse(sessionID); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid session_id",
		})
	}
	orderID := ctx.Params("order_id")
	if orderID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "order_id required",
		})
	}

	transaction, err := p.service.GetTransactionStatus(orderID)
	//Error from midtrans
	if err != nil && transaction == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Transaction", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	//Error from db
	if err != nil && transaction != nil {
		return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"success": true,
			"message": utils.MessageDataUpdated("Transaction", utils.MessageStatusFailed),
			"error":   err.Error(),
			"data":    transaction,
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataRead("Transaction", utils.MessageStatusSuccess),
		"data":    transaction,
	})
}
func (p *PaymentController) MidtransNotification(ctx *fiber.Ctx) error {
	var notif coreapi.TransactionStatusResponse

	if err := ctx.BodyParser(&notif); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := p.service.MidtransNotification(notif); err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	// DB update
	paymentStatus := utils.TranslateMidtransPaymentStatus(notif)
	transaction, err := services.PatchStatusTransaction(notif.OrderID, paymentStatus)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if paymentStatus == utils.TransactionStatusSuccess {
		_, err := services.PatchCartStatus(int64(transaction.CartID), string(utils.CartStatusCheckedOut))
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (p *PaymentController) CancelTransaction(ctx *fiber.Ctx) error {
	sessionID := ctx.Get(utils.SessionIDHeader)
	if sessionID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": utils.SessionIDMessage,
		})
	}
	cartID, err := services.CheckCartExist(sessionID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Cart", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	if cartID == -1 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Cart", utils.MessageStatusFailed),
		})
	}

	orderID := ctx.Params("order_id")
	if orderID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "order_id required",
		})
	}

	resp, err := p.service.CancelTransaction(orderID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataUpdated("Transaction", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}
	var paymentStatus string
	var message string
	if resp.TransactionStatus == "cancel" {
		paymentStatus = "CANCELLED"
		message = "Transaction successfully cancelled"
	} else {
		paymentStatus = strings.ToUpper(resp.TransactionStatus)
		message = "Transaction is pending to cancel"
	}
	// DB update
	transaction, err := services.PatchStatusTransaction(orderID, paymentStatus)
	if err != nil {
		grossamt, _ := strconv.ParseFloat(resp.GrossAmount, 64)
		secTransaction := models.Transaction{
			OrderID:     resp.OrderID,
			GrossAmount: int64(grossamt),
			Status:      resp.TransactionStatus,
			PaymentType: resp.PaymentType,
		}
		return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"success": true,
			"message": message + ", pending update to database",
			"data":    secTransaction,
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Cancelation success",
		"data":    transaction,
	})
}

func GetTransactionDetails(ctx *fiber.Ctx) error {
	sessionID := ctx.Get(utils.SessionIDHeader)
	if sessionID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": utils.SessionIDMessage,
		})
	}
	if _, err := uuid.Parse(sessionID); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid session_id",
		})
	}

	orderID := ctx.Params("order_id")
	if orderID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "order_id required",
		})
	}

	transactionDetails, err := services.GetTransactionDetailsByTransactionByOrderID(orderID)
	if err != nil {
		return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"success": false,
			"message": utils.MessageDataRead("Transaction Details", utils.MessageStatusFailed),
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": utils.MessageDataRead("Transaction Details", utils.MessageStatusSuccess),
		"data":    transactionDetails,
	})
}
