package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	database "github.com/Karalakrepp/Golang/Ecommerce_GO/db"
	jwttoken "github.com/Karalakrepp/Golang/Ecommerce_GO/jwt_token"
	"github.com/Karalakrepp/Golang/Ecommerce_GO/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HandlerAPI struct {
	l      *log.Logger
	storer database.Storer
}

var (
	db                = database.DBSet()
	client            = db.Client
	UserCollection    = db.UserData(client, "Users")
	ProductCollection = db.ProductData(client, "Products")
	Validate          = validator.New()
)

func NewHandlerAPI(l *log.Logger, storer database.Storer) *HandlerAPI {
	return &HandlerAPI{
		l:      l,
		storer: storer,
	}
}

func (s *HandlerAPI) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var req models.LoginReq
		var founduser models.User
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}
		//Compile your password

		ok, msg := ComparePassword(*founduser.Password, *req.Password)

		if !ok {
			c.JSON(http.StatusBadRequest, "Email or Password not match")
			s.l.Fatalln(msg)
			return
		}
		//Create token
		token, refleshToken, err := jwttoken.GenerateToken(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"err": err})
			s.l.Fatal("Token Error")
			return

		}
		jwttoken.UpdateToken(token, refleshToken, founduser.User_ID)

		c.JSON(http.StatusOK, founduser)

	}
}
func (s *HandlerAPI) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()
		var req = new(models.SignUpUser)

		err := c.BindJSON(req)
		if err != nil {
			msg := map[string]error{
				"Error": err,
			}
			c.JSON(http.StatusBadRequest, msg)
			s.l.Fatal(msg)
			return
		}
		validateError := Validate.Struct(req)

		if validateError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateError,
			})
			return
		}

		//If it is exist, find the User
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": req.Email})

		if err != nil {
			s.l.Panic(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			msg := "User already exist"
			s.l.Panic(msg)

			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": req.Phone})
		if err != nil {
			s.l.Panic(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			msg := "the phone is being used by another user"
			s.l.Panic(msg)

			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Match with user
		var user = models.NewUser(*req.Email, *req.Password)

		encPassword := HashPassword(*user.Password)
		user.Password = &encPassword

		user.First_Name = req.First_Name
		user.Last_Name = req.Last_Name
		user.Phone = req.Phone

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		//token section
		token, refleshToken, _ := jwttoken.GenerateToken(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refleshToken

		//continue
		user.UserCart = make([]models.ProductUser, 0)
		user.Order_Status = make([]models.Order, 0)
		user.Address_Details = make([]models.Address, 0)

		result, insertErr := UserCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"user": result})

	}
}

//********************************************Product Admin And View ////////////////********************************///////////////////////***************

func (s *HandlerAPI) ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		var product models.Product

		defer cancel()

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Err": err.Error(),
			})
			s.l.Fatal(err)
			return

		}
		product.Product_ID = primitive.NewObjectID()

		_, productErr := ProductCollection.InsertOne(ctx, product)

		if productErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Database inserting error",
			})
			s.l.Fatal("Database irserting Error")

			return
		}
		defer cancel()

		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")

	}
}

func (s *HandlerAPI) SeachProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		cur, err := ProductCollection.Find(ctx, bson.D{{}})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
			s.l.Fatal("Product Search Err")
			return

		}

		if err := cur.All(ctx, &productList); err != nil {
			s.l.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return

		}
		defer cur.Close(ctx)
		if err := cur.Err(); err != nil {
			s.l.Fatal(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productList)

	}
}

func (s *HandlerAPI) SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {

		var searchingProduct models.Product

		var query = c.Query("name")

		if query == "" {
			s.l.Println("the given query is empty")
			c.IndentedJSON(http.StatusBadRequest, "the given query is empty")
			c.Abort()
			return

		}
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*15)

		defer cancel()
		cur, err := ProductCollection.Find(ctx, bson.M{
			"product_name": bson.M{"$regex": query}})

		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}

		curErr := cur.All(ctx, &searchingProduct)

		if curErr != nil {

			s.l.Println(curErr)
			c.IndentedJSON(http.StatusBadRequest, "Invalid Request")
			c.Abort()
			return
		}

		defer cur.Close(ctx)

		if err := cur.Err(); err != nil {
			s.l.Println(err)
			c.IndentedJSON(http.StatusBadRequest, "invalid request")
			return
		}

		defer cancel()

		c.IndentedJSON(http.StatusOK, searchingProduct)

	}
}

//********************************************PASSWORD HASHED ////////////////********************************///////////////////////***************//

func HashPassword(password string) (encpass string) {

	//hash password this give you bytes and error
	bytePass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	//error handling
	if err != nil {
		log.Fatal(err)
	}

	//return string hashed password
	return string(bytePass)

}

// Compare your given and user password
func ComparePassword(encpassword, password string) (bool, string) {

	valid := true
	msg := ""

	//this func take 2 bytes parameters our passwords r string. we need to convert these passwords into bytes
	if err := bcrypt.CompareHashAndPassword([]byte(encpassword), []byte(password)); err != nil {

		msg = "Email or Password not match"
		valid = false
		return valid, msg
	}

	return valid, msg

}
