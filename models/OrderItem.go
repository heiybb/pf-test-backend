package models

type OrderItem struct {
	Id           uint64  `json:"id" gorm:"primary_key"`
	OrderID      uint64  `json:"order_id"`
	PricePerUnit string `json:"price_per_unit"`
	Quantity     int64  `json:"quantity"`
	Product      string  `json:"product"`
}
