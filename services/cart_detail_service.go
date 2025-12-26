package services

import (
	"strconv"

	"github.com/Rzero6/self-checkout-api/models"
	"github.com/Rzero6/self-checkout-api/repositories"
)

func AddProductsToCart(cart models.Cart, product models.Product, quantity int) (string, *models.CartDetail, error) {
	cartDetail, err := repositories.CheckExistingCartDetail(int64(cart.ID), int64(product.ID))
	if err != nil {
		return err.Error(), nil, err
	}
	var message string
	if cartDetail == nil {
		//Create
		result, err := repositories.CreateCartDetail(int64(cart.ID), product, quantity)
		if err != nil {
			return err.Error(), nil, err
		}
		message = strconv.Itoa(result.Quantity) + "x " + result.ProductName + " has been added to the Cart"
		return message, result, nil
	} else {
		//Update
		result, err := repositories.UpdateProductInCart(int64(cartDetail.ID), cartDetail.Quantity+quantity)
		if err != nil {
			return err.Error(), nil, err
		}
		message = result.ProductName + " quantity changed to " + strconv.Itoa(result.Quantity)
		return message, result, nil
	}
}

func GetCartDetailsBySessionID(sessionID string) ([]models.CartDetail, error) {
	return repositories.GetCartDetailsBySessionID(sessionID)
}

func DeleteProductFromCart(id int64) (*models.CartDetail, error) {
	return repositories.DeleteProductFromCart(id)
}
func UpdateProductInCart(id int64, quantity int) (*models.CartDetail, error) {
	return repositories.UpdateProductInCart(id, quantity)
}

func DeleteAllProductFromCart(cartID int64) error {
	return repositories.DeleteAllProductFromCart(cartID)
}
