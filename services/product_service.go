package services

import (
	"github.com/Rzero6/self-checkout-api/models"
	"github.com/Rzero6/self-checkout-api/repositories"
)

func GetAllProducts(page, limit int) ([]models.Product, int, error) {
	products, err := repositories.GetAllProducts(page, limit)
	if err != nil {
		return nil, 0, err
	}

	totalItems, err := repositories.GetProductsCount()
	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func GetProductByBarcode(barcode string) (*models.Product, error) {
	return repositories.GetProductByBarcode(barcode)
}

func GetProductRandom() (*models.Product, error) {
	return repositories.GetProductRandom()
}

func CreateProduct(item models.Product) (*models.Product, error) {
	return repositories.CreateProduct(item)
}

func GetDonationsProduct() ([]models.Product, error) {
	return repositories.GetDonationsProduct()
}
