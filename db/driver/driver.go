package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	COLLECTION = "birthdays"
	DATABASE   = "wadsworth-birthday"
)

/*
 * Global variables can go here
 */
var (
	logger = log.New(os.Stderr, "[DRIVER] ", log.LstdFlags)
)

type GuildUserPair struct {
	GuildId string
	UserId  string
}

type BirthdayDocument struct {
	GuildUserPair GuildUserPair
	Day           int
	Month         int
}

/*
 * Loads the .env
 */
func loadDotEnv() {
	logger.Println("Loading .env file...")
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Error loading .env file:", err)
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
		logger.Fatal("Failed to connect to DB:", err)
		panic(err)
	}

	// sending a ping to confirm a successful connection
	var result bson.M
	if err := client.Database(DATABASE).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		logger.Println("Couldn't ping DB:", err)
		panic(err)
	}
	logger.Println("Successfully connected to DB")

	return client
}

func insertBirthday(doc BirthdayDocument, client *mongo.Client) error {
	logger.Println("Inserting birthday document...")
	coll := client.Database(DATABASE).Collection(COLLECTION)
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		logger.Println("Failed to insert birthday document:", err)
	} else {
		logger.Println("Inserted document with _id:", result.InsertedID)
	}
	return err
}

func deleteBirthday(guildUserPair GuildUserPair, client *mongo.Client) error {
	logger.Println("Deleting birthday...")
	filter := bson.M{
		"guilduserpair.guildid": guildUserPair.GuildId,
		"guilduserpair.userid":  guildUserPair.UserId,
	}
	coll := client.Database(DATABASE).Collection(COLLECTION)
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		logger.Println("Failed to delete birthday document:", err)
	} else {
		logger.Println("Deleted document with _id. Count:", result.DeletedCount)
	}
	return err
}

func main() {
	// loading .env
	loadDotEnv()

	// connecting to db
	client := connectToDB()

	guildUserPair1 := GuildUserPair{
		GuildId: "some_server_id_1",
		UserId:  "PacoDaTaco",
	}

	// guildUserPair2 := GuildUserPair{
	// 	GuildId: "some_server_id_2",
	// 	UserId:  "PacoDaTaco",
	// }

	// testBirthdayDoc1 := BirthdayDocument{
	// 	GuildUserPair: guildUserPair1,
	// 	Day:           31,
	// 	Month:         10,
	// }

	// testBirthdayDoc2 := BirthdayDocument{
	// 	GuildUserPair: guildUserPair2,
	// 	Day:           31,
	// 	Month:         10,
	// }

	// insertBirthday(testBirthdayDoc1, client)
	// insertBirthday(testBirthdayDoc2, client)

	deleteBirthday(guildUserPair1, client)
}
