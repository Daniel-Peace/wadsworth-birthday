package main

import (
	"context"
	"driver/db"
	custom_utils "driver/utils"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
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

var (
	logger = log.New(os.Stderr, "[DRIVER] ", log.LstdFlags)
)

/*
 * Loads the .env
 */
func loadDotEnv() {
	logger.Printf("[%s] [%s] - Loading .env file...",
		custom_utils.ColorizeString("loadDotEnv", custom_utils.Magenta),
		custom_utils.ColorizeString(custom_utils.WORKING_STATUS, custom_utils.Yellow),
	)
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("[%s] [%s] - %v",
			custom_utils.ColorizeString("loadDotEnv", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
			err,
		)
	}
	logger.Printf("[%s] [%s]",
		custom_utils.ColorizeString("loadDotEnv", custom_utils.Magenta),
		custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
	)
}

/*
 * Connects driver to db
 */
func connectToDB() *mongo.Client {
	logger.Printf("[%s] [%s] - Connecting to DB...",
		custom_utils.ColorizeString("connectToDB", custom_utils.Magenta),
		custom_utils.ColorizeString(custom_utils.WORKING_STATUS, custom_utils.Yellow),
	)

	// getting the URI from the .env
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		logger.Fatalf("[%s] [%s] - %s",
			custom_utils.ColorizeString("connectToDB", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
			"Failed to find environment variable MONGODB_URI",
		)
	}

	// sedtting API version
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// creating full URI with options
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// creating client and connecting to db
	client, err := mongo.Connect(opts)
	if err != nil {
		logger.Fatalf("[%s] [%s] - %v",
			custom_utils.ColorizeString("connectToDB", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
			err,
		)
	}

	// sending a ping to confirm a successful connection
	var result bson.M
	if err := client.Database(DATABASE_NAME).RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		logger.Fatalf("[%s] [%s] - %v",
			custom_utils.ColorizeString("connectToDB", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
			err,
		)
	}

	logger.Printf("[%s] [%s]",
		custom_utils.ColorizeString("connectToDB", custom_utils.Magenta),
		custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
	)

	return client
}

type BirthdayPostRequest struct {
	ServerId string
	UserId   string
	Day      int
	Month    int
}

type Server struct {
	Database *db.MongoDB
}

func (s *Server) checkForBirthday(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) getActiveBirthdays(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) updateBirthday(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) insertBirthday(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Printf("[%s] [%s] - Unmarshling json...",
		custom_utils.ColorizeString("insertBirthday", custom_utils.Magenta),
		custom_utils.ColorizeString(custom_utils.WORKING_STATUS, custom_utils.Yellow),
	)
	var birthdayDocument db.BirthdayDocument
	err = json.Unmarshal(body, &birthdayDocument)
	if err != nil {
		logger.Printf("[%s] [%s] - %v",
			custom_utils.ColorizeString("insertBirthday", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
			err,
		)
	} else {
		logger.Printf("[%s] [%s] - Good json!",
			custom_utils.ColorizeString("insertBirthday", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
		)
		logger.Printf("[%s] [%s]\n--- JSON ---\n%s\n--- END ----",
			custom_utils.ColorizeString("insertBirthday", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.DATA, custom_utils.Blue),
			custom_utils.ColorizeString(string(body), custom_utils.Cyan),
		)
	}

	filter := bson.M{
		"guilduserpair.guildid": birthdayDocument.GuildUserPair.GuildId,
		"guilduserpair.userid":  birthdayDocument.GuildUserPair.UserId,
	}

	_, err = s.Database.FindOne(context.TODO(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Printf("[%s] [%s] - Birthday not found",
				custom_utils.ColorizeString("insertBirthday", custom_utils.Magenta),
				custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
			)

			err = s.Database.InsertOne(context.TODO(), birthdayDocument)
			if err != nil {
				logger.Printf("[%s] [%s] - %v",
					custom_utils.ColorizeString("insertBirthday", custom_utils.Magenta),
					custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
					err,
				)
			} else {
				logger.Printf("[%s] [%s] - Added birthday to database",
					custom_utils.ColorizeString("insertBirthday", custom_utils.Magenta),
					custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
				)
			}
		} else {
			logger.Printf("[%s] [%s] - %v",
				custom_utils.ColorizeString("insertBirthday", custom_utils.Magenta),
				custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
				err,
			)
		}
	} else {
		logger.Printf("[%s] [%s] - Birthday already exists",
			custom_utils.ColorizeString("insertBirthday", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
		)
	}
}

func main() {
	// loading .env
	loadDotEnv()

	// connecting to db
	client := connectToDB()

	// creating new instance of mongodb
	database := db.NewMongoDB(client, DATABASE_NAME, COLLECTION_NAME, logger)

	// creating an instance of the server struct to pass the db onto the handlers
	server := &Server{
		Database: database,
	}

	// setting up handlers
	http.HandleFunc("/check-for-bday", server.checkForBirthday)
	http.HandleFunc("/get-active-bday", server.getActiveBirthdays)
	http.HandleFunc("/update-bday", server.updateBirthday)
	http.HandleFunc("/insert-bday", server.insertBirthday)
	http.HandleFunc("/delete-bday", server.insertBirthday)

	// starting http server
	logger.Printf("[%s] [%s] - Starting http srever...",
		custom_utils.ColorizeString("main", custom_utils.Magenta),
		custom_utils.ColorizeString(custom_utils.WORKING_STATUS, custom_utils.Yellow),
	)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logger.Fatalf("[%s] [%s] - %v",
			custom_utils.ColorizeString("main", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
			err,
		)
	}

	// guildUserPair1 := db.GuildUserPair{
	// 	GuildId: "some_server_id_1",
	// 	UserId:  "PacoDaTaco",
	// }

	// guildUserPair2 := db.GuildUserPair{
	// 	GuildId: "some_server_id_2",
	// 	UserId:  "PacoDaTaco",
	// }

	// guildUserPair3 := db.GuildUserPair{
	// 	GuildId: "some_server_id_2",
	// 	UserId:  "GenericUser",
	// }

	// testBirthdayDoc1 := db.BirthdayDocument{
	// 	GuildUserPair: guildUserPair1,
	// 	Day:           31,
	// 	Month:         10,
	// }

	// testBirthdayDoc2 := db.BirthdayDocument{
	// 	GuildUserPair: guildUserPair2,
	// 	Day:           31,
	// 	Month:         10,
	// }

	// testBirthdayDoc3 := db.BirthdayDocument{
	// 	GuildUserPair: guildUserPair3,
	// 	Day:           31,
	// 	Month:         10,
	// }

	// filter1 := bson.M{
	// 	"guilduserpair.guildid": guildUserPair1.GuildId,
	// 	"guilduserpair.userid":  guildUserPair1.UserId,
	// }

	// filter2 := bson.M{
	// 	"guilduserpair.guildid": guildUserPair2.GuildId,
	// 	"guilduserpair.userid":  guildUserPair2.UserId,
	// }

	// filter3 := bson.M{
	// 	"guilduserpair.guildid": guildUserPair3.GuildId,
	// 	"guilduserpair.userid":  guildUserPair3.UserId,
	// }

	// filter4 := bson.M{
	// 	"guilduserpair.guildid": "some_server_id_2",
	// 	"day":                   31,
	// 	"month":                 10,
	// }

	// database.DeleteOne(context.TODO(), filter1)
	// database.DeleteOne(context.TODO(), filter2)
	// database.DeleteOne(context.TODO(), filter3)

	// database.FindOne(context.TODO(), filter1)
	// database.FindOne(context.TODO(), filter2)
	// database.FindOne(context.TODO(), filter3)

	// database.InsertOne(context.TODO(), testBirthdayDoc1)
	// database.InsertOne(context.TODO(), testBirthdayDoc2)
	// database.InsertOne(context.TODO(), testBirthdayDoc3)

	// database.FindOne(context.TODO(), filter1)
	// database.FindOne(context.TODO(), filter2)
	// database.FindOne(context.TODO(), filter3)

	// database.FindAll(context.TODO(), filter4)
}

// {"guilduserpair": { "guildid": "some_server_id_2", "userid": "GenericUser" },"day": 31,"month": 10}
