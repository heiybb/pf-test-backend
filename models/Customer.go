package models

type Customer struct {
	Id          string   `bson:"user_id"`
	Login       string   `bson:"login"`
	Password    string   `bson:"password"`
	Name        string   `bson:"name"`
	CompanyId   int64   `bson:"company_id"`
	CreditCards []string `bson:"credit_cards"`
}
