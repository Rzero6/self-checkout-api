package services

import (
	"github.com/Rzero6/self-checkout-api/internal/models"
	"github.com/Rzero6/self-checkout-api/internal/repositories"
)

func CheckCartExist(sessionID string) (int64, error) {
	return repositories.CheckCartExist(sessionID)
}
func CreateCart(cart models.Cart) (*models.Cart, error) {
	return repositories.CreateCart(cart)
}

func GetActiveCartBySession(sessionID string) (*models.Cart, error) {
	return repositories.GetActiveCartBySession(sessionID)
}

func PatchCartStatus(cartID int64, status string) (*models.Cart, error) {
	return repositories.PatchCartStatus(cartID, status)
}
