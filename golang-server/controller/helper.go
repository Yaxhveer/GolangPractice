package controller

import (
	"context"
	"fmt"
	"golang-server/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// helper function not exported

func insertOneMovie(movie model.Netflix) any {

	fmt.Println(movie.Movie, movie.Watched)
	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted 1 movie in db ", inserted.InsertedID)
	return inserted.InsertedID
}

func updateOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Modified Count ", result.ModifiedCount)
}

func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId);

	filter := bson.M{"_id": id};
	result, err := collection.DeleteOne(context.Background(), filter);

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted Count ", result.DeletedCount)
}

func deleteAllMovie() {

	result, err := collection.DeleteMany(context.Background(), bson.M{});

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted count ", result.DeletedCount)
}

func getAllMovie() []bson.M {
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	var movies []bson.M

	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)

		if err != nil {
			fmt.Println(err)
		}

		movies = append(movies, movie)
	}

	defer cursor.Close(context.Background())

	return movies
}