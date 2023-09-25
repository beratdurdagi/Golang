package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Karalakrepp/Golang/Ecommerce_GO/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *HandlerAPI) AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*20)

		var address models.Address

		defer cancel()

		uid := c.Query("id")
		if uid == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, "Invalid Code")
			c.Abort()
			return

		}

		addr, err := primitive.ObjectIDFromHex(uid)

		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}

		address.Address_id = primitive.NewObjectID()

		if err := c.BindJSON(&address); err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotAcceptable, err.Error())
			return
		}
		defer cancel()

		//Match  user_id with addr
		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: addr}}}}
		//unwind them
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$addr"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var addressInfo []bson.M

		if err = pointcursor.All(ctx, &addressInfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: address}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.IndentedJSON(400, "Not Allowed ")
		}
		defer cancel()
		ctx.Done()

	}
}

func (s *HandlerAPI) EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, "invalid url")
			c.Abort()
			return
		}

		uid, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		var editAddr models.Address
		if err := c.BindJSON(&editAddr); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())

		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: uid}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editAddr.House}, {Key: "address.0.street_name", Value: editAddr.Street}, {Key: "address.0.city_name", Value: editAddr.City}, {Key: "address.0.pin_code", Value: editAddr.Pincode}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something Went Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Updated the Home address")
	}
}

func (s *HandlerAPI) EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, err)
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House}, {Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: "address.1.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "something Went wrong")
			return
		}
		defer cancel()
		ctx.Done()

		c.IndentedJSON(200, "Successfully updated the Work Address")

	}
}
func (s *HandlerAPI) DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			c.Abort()
			return
		}
		uid, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		newAddr := make([]models.Address, 0)

		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*25)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: uid}}
		update := bson.D{{
			Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: newAddr}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Server Error")
			return
		}

		c.IndentedJSON(http.StatusOK, "Address Deleted")
	}
}
