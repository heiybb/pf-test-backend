package models

type OrderVO struct {
	OrderName       string
	CustomerCompany string
	CustomerName    string
	OrderDate       string
	DeliveryAmount  float64
	TotalAmount     float64
	FormatDate      string
}
