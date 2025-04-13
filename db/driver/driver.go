package main

import (
	"context"
	"driver/db"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	COLLECTION_NAME = "birthdays"
	DATABASE_NAME   = "wadsworth-birthday"
)

type Color string

/*
 * Ascii colors for logs
 */
var (
	Reset   Color = "\033[0m"
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Blue    Color = "\033[34m"
	Magenta Color = "\033[35m"
	Cyan    Color = "\033[36m"
	Gray    Color = "\033[37m"
	White   Color = "\033[97m"
)

/*
 * Log strings
 */
var (
	Working_status = "WORKING"
	Success_status = "SUCCESS"
	Error_status   = "ERROR"
)

var (
	logger = log.New(os.Stderr, "[DRIVER] ", log.LstdFlags)
)

func colorizeString(s string, c Color) string {
	return string(c) + s + string(Reset)
}

/*
 * Loads the .env
 */
func loadDotEnv() {
	logger.Printf("[%s] [%s] - Loading .env file...", colorizeString("loadDotEnv", Magenta), colorizeString(Working_status, Yellow))
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("[%s] [%s] - %v", colorizeString("loadDotEnv", Magenta), colorizeString(Error_status, Red), err)
	}
	logger.Printf("[%s] [%s]", colorizeString("loadDotEnv", Magenta), colorizeString(Success_status, Green))
}

/*
 * Connects driver to db
 */
func connectToDB() *mongo.Client {
	logger.Printf("[%s] [%s] - Connecting to DB...", colorizeString("connectToDB", Magenta), colorizeString(Working_status, Yellow))

	// getting the URI from the .env
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		logger.Fatalf("[%s] [%s] - %s", colorizeString("connectToDB", Magenta), colorizeString(Error_status, Red), "Failed to find environment variable MONGODB_URI")
	}

	// sedtting API version
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// creating full URI with options
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// creating client and connecting to db
	client, err := mongo.Connect(opts)
	if err != nil {
		logger.Fatalf("[%s] [%s] - %v", colorizeString("connectToDB", Magenta), colorizeString(Error_status, Red), err)
		panic(err)
	}

	// sending a ping to confirm a successful connection
	var result bson.M
	if err := client.Database(DATABASE_NAME).RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		logger.Fatalf("[%s] [%s] - %v", colorizeString("connectToDB", Magenta), colorizeString(Error_status, Red), err)
	}

	logger.Printf("[%s] [%s]", colorizeString("connectToDB", Magenta), colorizeString(Success_status, Green))

	return client
}

func main() {
	// loading .env
	loadDotEnv()

	// connecting to db
	client := connectToDB()

	// creating new instance of mongodb
	database := db.NewMongoDB(client, DATABASE_NAME, COLLECTION_NAME)

	guildUserPair1 := db.GuildUserPair{
		GuildId: "some_server_id_1",
		UserId:  "PacoDaTaco",
	}

	guildUserPair2 := db.GuildUserPair{
		GuildId: "some_server_id_2",
		UserId:  "PacoDaTaco",
	}

	guildUserPair3 := db.GuildUserPair{
		GuildId: "some_server_id_2",
		UserId:  "GenericUser",
	}

	testBirthdayDoc1 := db.BirthdayDocument{
		GuildUserPair: guildUserPair1,
		Day:           31,
		Month:         10,
	}

	testBirthdayDoc2 := db.BirthdayDocument{
		GuildUserPair: guildUserPair2,
		Day:           31,
		Month:         10,
	}

	testBirthdayDoc3 := db.BirthdayDocument{
		GuildUserPair: guildUserPair3,
		Day:           12,
		Month:         4,
	}

	filter1 := bson.M{
		"guilduserpair.guildid": guildUserPair1.GuildId,
		"guilduserpair.userid":  guildUserPair1.UserId,
	}

	filter2 := bson.M{
		"guilduserpair.guildid": guildUserPair2.GuildId,
		"guilduserpair.userid":  guildUserPair2.UserId,
	}

	filter3 := bson.M{
		"guilduserpair.guildid": guildUserPair3.GuildId,
		"guilduserpair.userid":  guildUserPair3.UserId,
	}

	database.DeleteOne(context.TODO(), filter1)
	database.DeleteOne(context.TODO(), filter2)
	database.DeleteOne(context.TODO(), filter3)

	database.FindOne(context.TODO(), filter1)
	database.FindOne(context.TODO(), filter2)
	database.FindOne(context.TODO(), filter3)

	database.InsertOne(context.TODO(), testBirthdayDoc1)
	database.InsertOne(context.TODO(), testBirthdayDoc2)
	database.InsertOne(context.TODO(), testBirthdayDoc3)

	database.FindOne(context.TODO(), filter1)
	database.FindOne(context.TODO(), filter2)
	database.FindOne(context.TODO(), filter3)
}
