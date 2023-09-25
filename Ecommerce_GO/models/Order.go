package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Order_ID primitive.ObjectID `bson:"_id"`

	//Order_Cart is []ProductUser cause it include ProductUser
	Order_Cart     []ProductUser `json:"order_list"  bson:"order_list"`
	Orderered_At   time.Time     `json:"ordered_on"  bson:"ordered_on"`
	Price          int           `json:"total_price" bson:"total_price"`
	Discount       *int          `json:"discount"    bson:"discount"`
	Payment_Method Payment       `json:"payment_method" bson:"payment_method"`
}
