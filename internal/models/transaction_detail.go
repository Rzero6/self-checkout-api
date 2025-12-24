package models

type TransactionDetail struct {
	ID            uint64 `json:"id"`
	TransactionID uint64 `json:"transaction_id"`
	ProductID     uint64 `json:"product_id"`
	ProductName   string `json:"product_name"`
	Price         int64  `json:"price"`
	Quantity      int    `json:"quantity"`
	Subtotal      int64  `json:"subtotal"`
}
