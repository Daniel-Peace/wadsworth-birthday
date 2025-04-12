package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

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
	Working_status = "Working"
	Success_status = "SUCCESS"
	Error_status   = "ERROR"
)

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
	logger.Printf("[%s] [%s] Connecting to DB...", colorizeString("connectToDB", Magenta), colorizeString(Working_status, Yellow))

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

func insertBirthday(doc BirthdayDocument, client *mongo.Client) error {
	var err error
	logger.Println("Inserting birthday document...")
	if result, _ := birthdayExists(doc.GuildUserPair, client); result {
		logger.Println("birthday already entered")
	} else {
		coll := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
		result, err := coll.InsertOne(context.TODO(), doc)
		if err != nil {
			logger.Println("Failed to insert birthday document:", err)
		} else {
			logger.Println("Inserted document with _id:", result.InsertedID)
		}
	}
	return err
}

func deleteBirthday(guildUserPair GuildUserPair, client *mongo.Client) error {
	logger.Println("Deleting birthday...")
	filter := bson.M{
		"guilduserpair.guildid": guildUserPair.GuildId,
		"guilduserpair.userid":  guildUserPair.UserId,
	}
	coll := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		logger.Println("Failed to delete birthday document:", err)
	} else {
		logger.Printf("Attempted to delete %s's birthday from guild %s. Count: %s\n",
			colorizeString(guildUserPair.UserId, Magenta),
			colorizeString(guildUserPair.GuildId, Cyan),
			colorizeString(strconv.FormatInt(result.DeletedCount, 10), Blue),
		)
	}
	return err
}

func getBirhtday(guildUserPair GuildUserPair, client *mongo.Client) (BirthdayDocument, error) {
	logger.Println("Retrieving birthday...")
	filter := bson.M{
		"guilduserpair.guildid": guildUserPair.GuildId,
		"guilduserpair.userid":  guildUserPair.UserId,
	}
	opts := options.FindOne()
	coll := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
	var birthday BirthdayDocument
	err := coll.FindOne(context.TODO(), filter, opts).Decode(&birthday)
	if err != nil {
		logger.Printf("No birthday found: %s\n", colorizeString(err.Error(), Red))
	} else {
		logger.Printf("Retrieved %s's birthday from guild %s.\n",
			colorizeString(guildUserPair.UserId, Magenta),
			colorizeString(guildUserPair.GuildId, Cyan),
		)
	}
	return birthday, err
}

func birthdayExists(guildUserPair GuildUserPair, client *mongo.Client) (bool, error) {
	filter := bson.M{
		"guilduserpair.guildid": guildUserPair.GuildId,
		"guilduserpair.userid":  guildUserPair.UserId,
	}
	coll := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		logger.Printf("An error accured while checking if that birthday existed: %s\n", colorizeString(err.Error(), Red))
		return false, err
	} else {
		logger.Printf("Checked for %s's birthday from guild %s. Found: %s\n",
			colorizeString(guildUserPair.UserId, Magenta),
			colorizeString(guildUserPair.GuildId, Cyan),
			colorizeString(strconv.FormatInt(count, 10), Blue),
		)
		return (count > 0), err
	}
}

func getActiveBirthdays(guildId string, client *mongo.Client) []BirthdayDocument {
	currentTime := time.Now()

	currentMonth := currentTime.Month()
	currentDay := currentTime.Day()
	filter := bson.M{
		"guilduserpair.guildid": guildId,
		"month":                 currentMonth,
		"day":                   currentDay,
	}
	coll := client.Database(DATABASE_NAME).Collection(COLLECTION_NAME)

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	var results []BirthdayDocument
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	logger.Printf(
		"Found %s active birthday(s) from guild %s",
		colorizeString(strconv.Itoa(len(results)), Blue),
		colorizeString(guildId, Cyan))
	return results
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

	guildUserPair2 := GuildUserPair{
		GuildId: "some_server_id_2",
		UserId:  "PacoDaTaco",
	}

	guildUserPair3 := GuildUserPair{
		GuildId: "some_server_id_2",
		UserId:  "GenericUser",
	}

	testBirthdayDoc1 := BirthdayDocument{
		GuildUserPair: guildUserPair1,
		Day:           31,
		Month:         10,
	}

	testBirthdayDoc2 := BirthdayDocument{
		GuildUserPair: guildUserPair2,
		Day:           31,
		Month:         10,
	}

	testBirthdayDoc3 := BirthdayDocument{
		GuildUserPair: guildUserPair3,
		Day:           12,
		Month:         4,
	}

	deleteBirthday(guildUserPair3, client)

	insertBirthday(testBirthdayDoc1, client)
	insertBirthday(testBirthdayDoc2, client)
	insertBirthday(testBirthdayDoc3, client)

	// exists, err := birthdayExists(guildUserPair2, client)
	// if exists && err == nil {
	// 	getBirhtday(guildUserPair1, client)
	// 	fmt.Println(getBirhtday(guildUserPair2, client))
	// }
	getActiveBirthdays(guildUserPair2.GuildId, client)
}
