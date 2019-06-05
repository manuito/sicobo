package store

import (
	"biblio/application"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// StartConnect init use of MongoDb instance
func StartConnect() {

	// Set client options
	clientOptions := options.Client().ApplyURI(application.State.Config.MongoDb)

	// Connect to MongoDB
	cli, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		application.Error("Cannot connect to mongo instance", err)
	}

	// Check the connection
	err = cli.Ping(context.TODO(), nil)

	if err != nil {
		application.Error("No connection available with mongo instance", err)
	}

	application.Info("Connected to MongoDB!")

	client = cli
}

func loadOrInitDatabaseSpec(database *mongo.Database, name string) BookDatabaseSpec {

	var result BookDatabaseSpec
	collection := database.Collection("spec")

	findErr := collection.FindOne(context.TODO(), bson.D{{}}).Decode(&result)
	cnt, err := database.Collection("books").CountDocuments(context.TODO(), bson.D{{}})

	// Always updated fields
	result.TotalBooks = cnt
	result.LastLoadTime = time.Now().String()

	if findErr != nil || result.Name == "" {
		application.Info("Cannot find database spec "+name+", will create it", err)
		result.Name = name
		result.CreateTime = result.LastLoadTime
		insertResult, err := collection.InsertOne(context.TODO(), result)
		if err != nil {
			application.Error("Cannot insert spec to mongo instance", err)
		}

		application.Debug("Initialized mongo database "+name+" : ", insertResult.InsertedID)
	} else {
		_, err := collection.UpdateOne(
			context.TODO(),
			bson.D{{"name", name}},
			bson.D{{"$set", bson.D{
				{"totalbooks", result.TotalBooks}, {"lastloadtime", result.LastLoadTime}}}})
		if err != nil {
			application.Error("Cannot update spec to mongo instance", err)
		}
	}

	return result
}

func loadMongoDatabase(name string) *mongo.Database {
	return client.Database(name)
}

func insertBook(book Book) {
	insertResult, err := onCollection("books").InsertOne(context.TODO(), book)
	if err != nil {
		application.Error("Cannot insert to mongo instance", err)
	}

	// Incr locally
	active.BookDatabaseSpec.TotalBooks++

	application.Debug("Inserted a single document: ", insertResult.InsertedID)
}

func loadBook(isbn string) Book {

	var result Book
	findErr := onCollection("books").FindOne(context.TODO(), bson.D{{"isbn", isbn}}).Decode(&result)
	if findErr != nil {
		application.Debug("Book with isbn " + isbn + " doesn't exist yet")
	}

	return result
}

func loadBooks() []Book {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var results = make([]Book, 0)

	cur, findErr := onCollection("books").Find(ctx, bson.D{})
	if findErr != nil {
		application.Error(findErr)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result Book
		decErr := cur.Decode(&result)
		if decErr != nil {
			application.Error("cannot decode result", decErr)
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return results
}

func updateBook(book *Book) {
	_, err := onCollection("books").UpdateOne(
		context.TODO(),
		bson.D{{"isbn", book.Isbn}},
		bson.D{{"$set",
			bson.D{{"picture", book.Picture},
				{"publisheddate", book.PublishedDate},
				{"category", book.Category},
				{"snippet", book.Snippet},
				{"pagecount", book.PageCount},
				{"candidatedetails", bson.D{{"collections", book.CandidateDetails.Collections},
					{"titles", book.CandidateDetails.Titles},
					{"pictures", book.CandidateDetails.Pictures}}}}}})
	if err != nil {
		application.Error("Cannot update book to mongo instance", err)
	}
}

func onCollection(name string) *mongo.Collection {
	return active.Database.Collection(name)
}
