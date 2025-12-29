package models

import "time"

type Transaction struct {
	ID          uint64              `json:"id"`
	CartID      uint64              `json:"cart_id"`
	OrderID     string              `json:"order_id"`
	GrossAmount int64               `json:"amount"`
	Status      string              `json:"status"`
	PaymentType string              `json:"payment_type"`
	QRISLink    *string             `json:"qris_link"`
	ExpireTime  *time.Time          `json:"expire_time"`
	CreatedAt   *time.Time          `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}
