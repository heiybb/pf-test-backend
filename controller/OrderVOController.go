package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math"
	"net/http"
	"pf-test-backend/models"
	"strconv"
	"strings"
	"time"
)

func FindAll(c *gin.Context) {
	var orders []models.Order
	models.PostgresDB.Find(&orders)

	var orderVO []models.OrderVO
	for i := range orders {
		var voOrderName string
		var voCustomerCompany string
		var voCustomerName string
		var voOrderDate string
		var voDeliveryAmount float64 = 0
		var voTotalAmount float64 = 0

		order := orders[i]
		var customerName = order.CustomerName

		voOrderName = order.OrderName
		t, _ := time.Parse(time.RFC3339, order.CreateTime)

		voOrderDate = suffixFormat(t.Format("Jan 2, 3:04 PM"))
		var voFormatDate = t.String()

		//Find the customer record in the mongo db by the customer name
		var customerResult models.Customer
		filter := bson.D{{"user_id", customerName}}

		err := models.MongoDB.Collection("customer").FindOne(context.TODO(), filter).Decode(&customerResult)
		if err != nil {
			log.Fatal(err)
		}

		voCustomerName = customerResult.Name

		//Find the company name in the mongo db by the company id
		var companyResult models.Company
		var companyID = customerResult.CompanyId
		filter = bson.D{{"company_id", companyID}}
		err = models.MongoDB.Collection("customer_companies").FindOne(context.TODO(), filter).Decode(&companyResult)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(companyResult)

		voCustomerCompany = companyResult.CompanyName

		var itemsInOrder []models.OrderItem

		models.PostgresDB.Where("order_id=?", order.Id).Find(&itemsInOrder)
		if len(itemsInOrder) > 0 {
			for itemIndex := range itemsInOrder {
				var item models.OrderItem
				item = itemsInOrder[itemIndex]
				price, _ := strconv.ParseFloat(item.PricePerUnit, 64)
				voTotalAmount += price * float64(item.Quantity)

				var delivery []models.Delivery
				models.PostgresDB.Where("order_item_id=?", item.Id).Find(&delivery)
				if len(delivery) > 0 {
					for i2 := range delivery {
						voDeliveryAmount += price * float64(delivery[i2].Quantity)
					}
				}
			}
		}

		test := models.OrderVO{
			OrderName:       voOrderName,
			CustomerCompany: voCustomerCompany,
			CustomerName:    voCustomerName,
			OrderDate:       voOrderDate,
			DeliveryAmount:  toFixed(voDeliveryAmount, 2),
			TotalAmount:     toFixed(voTotalAmount, 2),
			FormatDate:      voFormatDate,
		}
		orderVO = append(orderVO, test)
		fmt.Println(test)
	}

	c.JSON(http.StatusOK, gin.H{"data": orderVO})
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func suffixFormat(timeString string) string {
	monPart := strings.Split(strings.Split(timeString, ",")[0], " ")[0]
	timePart := strings.Split(timeString, ",")[1]
	day := strings.Split(strings.Split(timeString, ",")[0], " ")[1]

	suffix := "th"
	switch day {
	case "1", "21", "31":
		suffix = "st"
	case "2", "22":
		suffix = "nd"
	case "3", "23":
		suffix = "rd"
	}

	return monPart + " " + day + suffix + "," + timePart
}
