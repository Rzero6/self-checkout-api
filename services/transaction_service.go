package services

import (
	"github.com/Rzero6/self-checkout-api/models"
	"github.com/Rzero6/self-checkout-api/repositories"
)

func CreateTransaction(transaction models.Transaction) (int64, error) {
	return repositories.CreateTransaction(transaction)
}

func GetTransactionByOrderID(orderID string, status string) (*models.Transaction, error) {

	transaction, tErr := repositories.GetTransactionByOrderID(orderID, status)
	if tErr != nil {
		return nil, tErr
	}
	details, dErr := GetTransactionDetailsByTransactionByOrderID(transaction.OrderID)
	if dErr != nil {
		return nil, dErr
	}
	transaction.Details = details

	return transaction, nil
}

func PatchStatusTransaction(orderID string, status string) (*models.Transaction, error) {
	return repositories.PatchStatusTransaction(orderID, status)
}

func CheckExistingTransaction(cartID int64) (*models.Transaction, error) {
	return repositories.CheckExistingTransaction(cartID)
}
