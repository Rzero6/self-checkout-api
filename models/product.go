package models

type Product struct {
	ID      uint64 `json:"id"`
	Barcode string `json:"barcode"`
	Name    string `json:"name"`
	Price   int64  `json:"price"`
}
