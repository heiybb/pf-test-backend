package models

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var PostgresDB *gorm.DB
var MongoDB *mongo.Database
var mongoSession *mongo.Session

func DataInit() {
	//Init the postgres database
	postgresInit()
	//Insert the orders to postgres database
	insertOrders()
	//Insert the order_items to postgres database
	insertOrderItems()
	//Insert the deliveries to postgres database
	insertDeliveries()

	//Init the mongo database
	mongoInit()
	//Insert the customers
	insertCustomers()
	//Insert the companies
	insertCustomerCompanies()
}

func postgresInit() {
	postgres, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=test dbname=test password=K5ec4f1Cc31d6817 sslmode=disable")

	if err != nil {
		panic("Failed to connect to postgres database!")
	}

	fmt.Println("Drop the existing table in the postgres database")
	postgres.Exec("DROP TABLE IF EXISTS orders")
	postgres.Exec("DROP TABLE IF EXISTS deliveries")
	postgres.Exec("DROP TABLE IF EXISTS order_items")

	postgres.AutoMigrate(&Delivery{}, &Order{}, &OrderItem{})
	PostgresDB = postgres
}

func mongoInit() {

	clientOptions := options.Client().ApplyURI("mongodb://test:K5ec4f1Cc31d6817@localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	MongoDB = client.Database("test")

	if err = MongoDB.Collection("customer").Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err = MongoDB.Collection("customer_companies").Drop(ctx); err != nil {
		log.Fatal(err)
	}
}

func insertOrders() {
	// Populate the orders
	ordersCSV, err := os.Open("./data/postgres/orders.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(ordersCSV)

	//Skip the first line
	if _, err := r.Read(); err != nil {
		panic(err)
	}

	// Iterate through the csv lines
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		selfID, err := strconv.ParseUint(record[0], 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		var orderCreateTime = record[1]
		var orderName = record[2]
		var orderCustomerName = record[3]
		PostgresDB.Exec("INSERT INTO orders VALUES(?,?,?,?)", selfID, orderCreateTime, orderName, orderCustomerName)
	}
}

func insertOrderItems() {
	// Populate the order_items
	orderItemsCSV, err := os.Open("./data/postgres/order_items.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(orderItemsCSV)

	//Skip the first line
	if _, err := r.Read(); err != nil {
		panic(err)
	}

	// Iterate through the csv lines
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		selfID, err := strconv.ParseUint(record[0], 10, 32)
		orderID, err := strconv.ParseUint(record[1], 10, 32)
		//price,err :=strconv.ParseUint(record[2], 10, 32)
		//pricePerUnit, err := strconv.ParseFloat(fmt.Sprintf("%.4f", price), 32)

		var quantity = record[3]
		var product = record[4]

		PostgresDB.Exec("INSERT INTO order_items VALUES(?,?,?,?,?)", selfID, orderID, record[2], quantity, product)
	}
}

func insertDeliveries() {
	// Populate the deliveries
	deliveriesCSV, err := os.Open("./data/postgres/deliveries.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(deliveriesCSV)

	//Skip the first line
	if _, err := r.Read(); err != nil {
		panic(err)
	}

	// Iterate through the csv lines
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		selfID, err := strconv.ParseUint(record[0], 10, 32)
		orderItemID, err := strconv.ParseUint(record[1], 10, 32)
		var deliveredQuantity = record[2]

		PostgresDB.Exec("INSERT INTO deliveries VALUES(?,?,?)", selfID, orderItemID, deliveredQuantity)
	}
}

func insertCustomers() {
	// Populate the customers
	customersCSV, err := os.Open("./data/mongo/customers.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(customersCSV)

	//Skip the first line
	if _, err := r.Read(); err != nil {
		panic(err)
	}

	// Iterate through the csv lines
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		user_id := record[0]
		login := record[1]
		password := record[2]
		name := record[3]
		company_id, err := strconv.ParseUint(record[4], 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		credit_cards := strings.Split(record[5], ",")

		customer := Customer{Id: user_id, Login: login, Password: password, Name: name, CompanyId: int64(company_id), CreditCards: credit_cards}

		insertResult, err := MongoDB.Collection("customer").InsertOne(context.TODO(), customer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}
}

func insertCustomerCompanies() {
	// Populate the insertCustomerCompanies
	customerCompaniesCSV, err := os.Open("./data/mongo/customer_companies.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(customerCompaniesCSV)

	//Skip the first line
	if _, err := r.Read(); err != nil {
		panic(err)
	}

	// Iterate through the csv lines
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		company_id, err := strconv.ParseUint(record[0], 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		company_name := record[1]

		customer := Company{Id: int64(company_id), CompanyName: company_name}

		insertResult, err := MongoDB.Collection("customer_companies").InsertOne(context.TODO(), customer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}
}
