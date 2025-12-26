package services

import (
	"github.com/Rzero6/self-checkout-api/models"
	"github.com/Rzero6/self-checkout-api/repositories"
)

func CreateTransaction(transaction models.Transaction) (int64, error) {
	return repositories.CreateTransaction(transaction)
}

func PatchStatusTransaction(orderID string, status string) (*models.Transaction, error) {
	return repositories.PatchStatusTransaction(orderID, status)
}

func CheckExistingTransaction(cartID int64) (*models.Transaction, error) {
	return repositories.CheckExistingTransaction(cartID)
}
