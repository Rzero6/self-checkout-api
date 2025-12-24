package repositories

import (
	"errors"

	"github.com/Rzero6/self-checkout-api/config"
	"github.com/Rzero6/self-checkout-api/internal/models"
	"github.com/jackc/pgx/v5"
)

func GetTransactionDetailsByTransactionByOrderID(orderID string) ([]models.TransactionDetail, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		SELECT td.product_name, td.price, td.quantity, td.subtotal
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.order_id = $1
	`
	rows, err := config.DB.Query(ctx, query, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	details := []models.TransactionDetail{}
	for rows.Next() {
		var d models.TransactionDetail
		rows.Scan(&d.ProductName, &d.Price, &d.Quantity, &d.Subtotal)
		details = append(details, d)
	}
	return details, nil
}

func CreateTransactionDetail(transactionID int64, cartDetail models.CartDetail) (*models.TransactionDetail, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		INSERT INTO transaction_details (transaction_id, product_id, product_name, price, quantity, subtotal)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, transaction_id, product_id, product_name, price, quantity, subtotal
	`

	var created models.TransactionDetail
	err := config.DB.QueryRow(ctx, query,
		transactionID, cartDetail.ProductID, cartDetail.ProductName, cartDetail.Price, cartDetail.Quantity, cartDetail.Subtotal).Scan(
		&created.ID,
		&created.TransactionID,
		&created.ProductID,
		&created.ProductName,
		&created.Price,
		&created.Quantity,
		&created.Subtotal,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}
