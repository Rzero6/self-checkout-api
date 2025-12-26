package repositories

import (
	"errors"

	"github.com/Rzero6/self-checkout-api/config"
	"github.com/Rzero6/self-checkout-api/models"
	"github.com/Rzero6/self-checkout-api/utils"
	"github.com/jackc/pgx/v5"
)

func CheckExistingTransaction(cartID int64) (*models.Transaction, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		SELECT order_id, amount, status, payment_type, qris_link, expire_time
		FROM transactions
		WHERE cart_id = $1 AND status ILIKE $2
		LIMIT 1
	`
	var tx models.Transaction
	err := config.DB.QueryRow(
		ctx, query, cartID, utils.TransactionStatusPending,
	).Scan(
		&tx.OrderID, &tx.GrossAmount, &tx.Status, &tx.PaymentType, &tx.QRISLink, &tx.ExpireTime,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &tx, nil
}

func CreateTransaction(transaction models.Transaction) (int64, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		INSERT INTO transactions
			(cart_id, order_id, amount, status, payment_type, qris_link, expire_time)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var transactionID int64
	err := config.DB.QueryRow(
		ctx,
		query,
		transaction.CartID,
		transaction.OrderID,
		transaction.GrossAmount,
		transaction.Status,
		transaction.PaymentType,
		transaction.QRISLink,
		transaction.ExpireTime,
	).Scan(&transactionID)

	if err != nil {
		return 0, err
	}

	return transactionID, nil
}

func PatchStatusTransaction(orderID string, status string) (*models.Transaction, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		UPDATE transactions
		SET status = $1
		WHERE order_id = $2
		RETURNING cart_id, order_id, amount, status, payment_type, qris_link, expire_time
	`

	var tx models.Transaction

	err := config.DB.QueryRow(
		ctx,
		query,
		status,
		orderID,
	).Scan(
		&tx.CartID,
		&tx.OrderID,
		&tx.GrossAmount,
		&tx.Status,
		&tx.PaymentType,
		&tx.QRISLink,
		&tx.ExpireTime,
	)

	if err != nil {
		return nil, err
	}
	return &tx, nil
}
