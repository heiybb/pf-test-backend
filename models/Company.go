package models

type Company struct {
	Id          int64   `bson:"company_id"`
	CompanyName string `bson:"company_name"`
}
