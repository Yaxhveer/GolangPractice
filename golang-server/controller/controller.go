package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-server/model"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "netflix"
const colName = "watchlist"

var collection *mongo.Collection

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	USER := os.Getenv("USER")
	PASS := os.Getenv("PASS")

	connectionString := "mongodb+srv://"+USER+":"+PASS+"@cluster0.uziio6u.mongodb.net/?retryWrites=true&w=majority"

	fmt.Println(connectionString)

	// client option
	clientOption := options.Client().ApplyURI(connectionString)

	// connection with mongodb
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongodb Connection Sucess")

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection Instance is ready", collection)
}

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")

	movies := getAllMovie()
	json.NewEncoder(w).Encode(movies)
}

func InsertMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		fmt.Println(7, err.Error())
	}

	id := insertOneMovie(movie)
	json.NewEncoder(w).Encode(id)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)

	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode("Done")
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode("Done")
}

func DeleteAllMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	deleteAllMovie()
	json.NewEncoder(w).Encode("Done")
}