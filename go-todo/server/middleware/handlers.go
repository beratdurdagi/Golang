package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/karalakrepp/Golang/go-react-todo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	loadTheEnv()
	ConnectDb()
}

func loadTheEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the .env file")
	}
}

func ConnectDb() *mongo.Client {
	collectionName := os.Getenv("DB_COLLECTION_NAME")
	dbName := os.Getenv("DB_NAME")
	MongoDb := os.Getenv("MONGODB_URL")

	clientOptions := options.Client().ApplyURI(MongoDb)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb!")

	collection = client.Database(dbName).Collection(collectionName)
	fmt.Println("collection instance created", collection)

	return client
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {

	/* HTML form verilerinin sunucuya gönderildiği durumlarda kullanılır. Bu biçim, verilerin URL kodlaması kullanılarak aktarıldığı bir formattır.
	Örneğin, "name=John&age=25" gibi veriler bu biçimde kodlanır. */
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	/*  Bu, tarayıcının, bu kaynağa erişimi olan herhangi bir etki alanından gelen isteklere izin verilmesi gerektiğini belirtir.
	Yıldız karakteri "" tüm etki alanlarına izin verir.
	 Bu, Cross-Origin Resource Sharing (CORS) politikalarını kullanarak farklı etki alanları arasında kaynak paylaşımını sağlar.*/

	w.Header().Set("Access-Control-Allow-Origin", "*")

	payload := getAllTasks()

	json.NewEncoder(w).Encode(payload)

}

func CreateTasks(w http.ResponseWriter, r *http.Request) {
	//HTTP isteklerini kontrol et
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var task models.ToDoList

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		log.Fatal(err.Error())
	}

	insertOneTask(task)

	json.NewEncoder(w).Encode(task)

}
func TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www.form-urlencoded")
	w.Header().Set("Access-Control-Allow", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	taskCompleted(params["id"])

	json.NewEncoder(w).Encode(params["id"])

}
func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www.form-urlencoded")
	w.Header().Set("Access-Control-Allow", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)

	undoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www.form-urlencoded")
	w.Header().Set("Access-Control-Allow", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	deleteOneTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}
func DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www.form-urlencoded")
	w.Header().Set("Access-Control-Allow", "*")

	count := deleteAllTasks()

	json.NewEncoder(w).Encode(count)
}

// /////////////////////// dbhandler///////////////////////
func getAllTasks() []primitive.M {

	cur, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	var results []primitive.M

	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func insertOneTask(task models.ToDoList) {
	res, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("insertedId:", res.InsertedID)
}

func taskCompleted(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	updated := bson.M{"$set": bson.M{"status": "true"}}
	res, err := collection.UpdateByID(context.Background(), filter, updated)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count:", res.ModifiedCount)

}

func deleteOneTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted count:", res.DeletedCount)

}

// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.12.0/mongo#Collection.DeleteMany
func deleteAllTasks() int64 {
	res, err := collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)

	return res.DeletedCount
}

func undoTask(task string) {

	id, _ := primitive.ObjectIDFromHex(task)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": "false"}}

	res, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count:", res.ModifiedCount)

}
