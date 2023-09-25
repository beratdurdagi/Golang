package jwttoken

import (
	"context"
	"log"
	"os"
	"time"

	database "github.com/Karalakrepp/Golang/Ecommerce_GO/db"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db             = database.DBSet()
	client         = db.Client
	UserCollection = db.UserData(client, "Users")
	SECRET_KEY     = os.Getenv("SECRET_KEY")
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

func GenerateToken(email string, first_name string, last_name string, user_id string) (signedToken string, signedRefreshToken string, err error) {

	claims := SignedDetails{
		Email:      email,
		First_Name: first_name,
		Last_Name:  last_name,
		Uid:        user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	reflesh_claims := SignedDetails{
		Email:      email,
		First_Name: first_name,
		Last_Name:  last_name,
		Uid:        user_id,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Fatal(err)
		return
	}

	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, reflesh_claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panicln(err)
		return
	}

	return token, refresh_token, err
}

func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The Token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}

func UpdateToken(token string, reflesh_token string, user_id string) {

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)

	var updatedObj primitive.D
	defer cancel()

	updatedObj = append(updatedObj, bson.E{Key: "token", Value: token})
	updatedObj = append(updatedObj, bson.E{Key: "reflesh_token", Value: reflesh_token})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updatedObj = append(updatedObj, bson.E{Key: "updatedat", Value: updated_at})

	upsert := true
	filter := bson.M{"user_id": user_id}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserCollection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updatedObj},
	}, &opt)

	defer cancel()
	if err != nil {

		log.Fatal(err)
		return
	}
}
