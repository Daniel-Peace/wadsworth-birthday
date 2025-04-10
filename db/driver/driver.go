package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

/*
 * Global variables can go here
 */
var (
	logger = log.New(os.Stderr, "[DRIVER] ", log.LstdFlags)
)

type BirthdayDocument struct {
	GuildId string
	UserId  string
	Day     int
	Month   int
}

/*
 * Loads the .env
 */
func loadDotEnv() {
	logger.Println("Loading .env file...")
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}
	logger.Println("Successfully loaded \".env\" file")
}

/*
 * Connects driver to db
 */
func connectToDB() *mongo.Client {
	logger.Println("Connecting to DB...")

	// getting the URI from the .env
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		logger.Fatal("Failed to find environment variable MONGODB_URI")
	}

	// sedtting API version
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// creating full URI with options
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// creating client and connecting to db
	client, err := mongo.Connect(opts)
	if err != nil {
		logger.Fatalf("Failed to connect to DB %v", err)
		panic(err)
	}

	// sending a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		fmt.Println("Something failed:", err)
		panic(err)
	}
	logger.Println("Successfully connected to DB")

	return client
}

func insertBirthday(doc BirthdayDocument, client *mongo.Client) error {
	logger.Println("Inserting birthday document...")
	coll := client.Database("wadsworth").Collection(doc.GuildId)
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Println("Failed to insert birthday document: ", err)
	} else {
		logger.Println("Inserted document with _id:", result.InsertedID)
	}
	return err
}

func main() {
	// loading .env
	loadDotEnv()

	// connecting to db
	client := connectToDB()

	testBirthdayDoc := BirthdayDocument{
		GuildId: "some_server_id",
		UserId:  "PacoDaTaco",
		Day:     31,
		Month:   10,
	}
	insertBirthday(testBirthdayDoc, client)
}
