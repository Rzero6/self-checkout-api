package models

type CartDetail struct {
	ID          uint64 `json:"id"`
	CartID      uint64 `json:"cart_id"`
	ProductID   uint64 `json:"product_id"`
	ProductName string `json:"product_name"`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
	Subtotal    int64  `json:"subtotal"`
}
