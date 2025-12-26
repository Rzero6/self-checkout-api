package services

import (
	"github.com/Rzero6/self-checkout-api/models"
	"github.com/Rzero6/self-checkout-api/repositories"
)

func GetTransactionDetailsByTransactionByOrderID(orderID string) ([]models.TransactionDetail, error) {
	return repositories.GetTransactionDetailsByTransactionByOrderID(orderID)
}

func CreateTransactionDetails(transactionID int64, cartDetails []models.CartDetail) ([]models.TransactionDetail, error) {
	transactionDetails := make([]models.TransactionDetail, 0, len(cartDetails))
	for _, cartDetail := range cartDetails {
		detail, err := repositories.CreateTransactionDetail(transactionID, cartDetail)
		if err != nil {
			return nil, err
		}
		transactionDetails = append(transactionDetails, *detail)
	}
	return transactionDetails, nil
}
