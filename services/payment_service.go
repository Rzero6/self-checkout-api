package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Rzero6/self-checkout-api/config"
	"github.com/Rzero6/self-checkout-api/models"
	"github.com/Rzero6/self-checkout-api/repositories"
	"github.com/Rzero6/self-checkout-api/utils"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type PaymentService struct {
	client    coreapi.Client
	serverKey string
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		client:    config.NewMidtransClient(),
		serverKey: config.NewMidtransClient().ServerKey,
	}
}

func (s *PaymentService) CreateTransaction(
	tx models.Transaction,
	paymentType string,
) (*coreapi.ChargeResponse, *models.Transaction, error) {

	var itemDetails []midtrans.ItemDetails
	for _, item := range tx.Details {
		itemDetails = append(itemDetails, midtrans.ItemDetails{
			ID:    strconv.FormatUint(item.ProductID, 10),
			Name:  item.ProductName,
			Price: item.Price,
			Qty:   int32(item.Quantity),
		})
	}

	req := &coreapi.ChargeReq{
		PaymentType: coreapi.CoreapiPaymentType(paymentType),
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  tx.OrderID,
			GrossAmt: tx.GrossAmount,
		},
		Items: &itemDetails,
	}

	resp, mtErr := s.client.ChargeTransaction(req)
	if mtErr != nil {
		return nil, nil, fmt.Errorf("midtrans charge failed: %w", mtErr)
	}
	tx.PaymentType = resp.PaymentType
	tx.Status = strings.ToUpper(resp.TransactionStatus)
	txPtr, err := HandlePaymentType(tx, *resp)
	if err != nil {
		return nil, nil, err
	}
	tx = *txPtr
	transactionID, err := CreateTransaction(tx)
	if err != nil {
		return nil, nil, err
	}
	tx.ID = uint64(transactionID)

	return resp, &tx, nil
}

func (s *PaymentService) MidtransNotification(
	notif coreapi.TransactionStatusResponse,
) error {

	if !utils.VerifyMidtransSignature(notif, s.serverKey) {
		return fmt.Errorf("invalid signature, not legit from Midtrans")
	}

	_, mtErr := s.client.CheckTransaction(notif.OrderID)
	if mtErr != nil {
		return mtErr
	}

	return nil
}

func (s *PaymentService) CancelTransaction(orderID string) (*coreapi.CancelResponse, error) {
	// Midtrans cancel
	resp, mtErr := s.client.CancelTransaction(orderID)
	if mtErr != nil {
		return nil, fmt.Errorf("midtrans cancel failed: %s", mtErr.Message)
	}
	return resp, nil
}

func (s *PaymentService) GetTransactionStatus(orderID string) (*models.Transaction, error) {
	resp, mtErr := s.client.CheckTransaction(orderID)
	if mtErr != nil {
		return nil, fmt.Errorf("midtrans check error: %s", mtErr.Message)
	}
	paymentStatus := utils.TranslateMidtransPaymentStatus(*resp)
	//Patch transaction status
	var result models.Transaction
	transaction, err := repositories.PatchStatusTransaction(orderID, paymentStatus)
	if err != nil {
		grossAmt, _ := strconv.ParseFloat(resp.GrossAmount, 64)
		result = models.Transaction{
			OrderID:     resp.OrderID,
			Status:      strings.ToUpper(resp.TransactionStatus),
			PaymentType: resp.PaymentType,
			GrossAmount: int64(grossAmt),
		}
		return &result, fmt.Errorf("failed to update transactions db: %s", err.Error())
	}
	//Get transaction details
	details, err := GetTransactionDetailsByTransactionByOrderID(orderID)
	if err != nil {
		result.Details = details
		return &result, fmt.Errorf("failed to read transactions db: %s", err.Error())
	}
	//Update Cart Status
	if transaction.Status == utils.TransactionStatusSuccess {
		_, err := PatchCartStatus(int64(transaction.CartID), string(utils.CartStatusCheckedOut))
		if err != nil {
			return &result, fmt.Errorf("failed to read transactions db: %s", err.Error())
		}
	}

	transaction.Details = details
	return transaction, nil
}

func HandlePaymentType(tx models.Transaction, chargeResp coreapi.ChargeResponse) (*models.Transaction, error) {
	switch tx.PaymentType {
	case string(coreapi.PaymentTypeQris):
		for _, action := range chargeResp.Actions {
			if action.Name == "generate-qr-code-v2" {
				tx.QRISLink = &action.URL
				break
			}
		}
		expireTime, err := utils.ParseJakartaToUTC(chargeResp.ExpiryTime)
		if err != nil {
			return nil, err
		}
		tx.ExpireTime = expireTime
		return &tx, err
	default:
		return nil, fmt.Errorf("Other payment method not available")
	}
}
