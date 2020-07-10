package models

type Order struct {
	Id           uint64 `json:"id" gorm:"primary_key"`
	CreateTime   string `json:"create_time"`
	OrderName    string `json:"order_name"`
	CustomerName string `json:"customer_name"`
}
