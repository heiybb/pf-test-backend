package models

type Delivery struct {
	Id          uint64 `json:"id" gorm:"primary_key"`
	OrderItemID uint64 `json:"order_item_id"`
	Quantity    int64 `json:"quantity"`
}
