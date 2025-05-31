package main

import (
	"context"
	"driver/db"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Daniel-Peace/go-logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	COLLECTION_NAME = "birthdays"
	DATABASE_NAME   = "wadsworth-birthday"
)

type DBResponseStatus int

const (
	SUCCESS DBResponseStatus = iota
	CONFLICT
	NOT_POSSIBLE
	ERROR
)

var (
	log = logger.NewGoLogger("DRIVER", os.Stdout, true, true, true)
)

type Server struct {
	Database *db.MongoDB
}

type ResponseBody struct {
	status_code int
	Status      DBResponseStatus
	Description string
	Data        string
}

// Loads the .env
func loadDotEnv() {
	log.StatusPrintln(logger.IN_PROGRESS, "Loading .env file...")
	err := godotenv.Load(".env")
	if err != nil {
		log.StatusFatalf(logger.ERROR, "%v", err)
	}
	log.StatusPrintln(logger.SUCCESS, "Successfully loaded .env")
}

// Connects driver to db
func connectToDB() *mongo.Client {
	log.StatusPrint(logger.IN_PROGRESS, "Connecting to DB...")

	// getting the URI from the .env
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.StatusFatalln(logger.ERROR, "Failed to find environment variable MONGODB_URI")
	}

	// sedtting API version
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	// creating full URI with options
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// creating client and connecting to db
	client, err := mongo.Connect(opts)
	if err != nil {
		log.StatusFatalf(logger.ERROR, "%v", err)
	}

	// sending a ping to confirm a successful connection
	var result bson.M
	if err := client.Database(DATABASE_NAME).RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		log.StatusFatalf(logger.ERROR, "%v", err)
	}
	log.StatusPrint(logger.SUCCESS, "Successfully connected to the DB")
	return client
}

func sendJsonResponse[T any](status_code int, body T, w http.ResponseWriter) error {
	// marshalling data
	bodyAsJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// setting and writing haeder
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)

	// writing body
	_, err = w.Write(bodyAsJson)
	return err
}

func buildAndSendResponse(w http.ResponseWriter, status_code int, status DBResponseStatus, description string, data string) error {
	// building response
	body := ResponseBody{
		status_code: status_code,
		Status:      status,
		Description: description,
		Data:        data,
	}

	// sending response
	return sendJsonResponse(status_code, body, w)
}

func sendFallbackResponse(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func sendFallbackIfError(w http.ResponseWriter, err error) {
	if err != nil {
		log.StatusPrintf(logger.ERROR, "%v", err)
		sendFallbackResponse(w)
	}
}

// Reads the body of the request and parses it into a struct
func parseRequestBody[T any](r *http.Request) (T, error) {
	var document T

	// reading body from request
	log.StatusPrintln(logger.IN_PROGRESS, "Reading request body")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.StatusPrintf(logger.ERROR, "%v", err)
		return document, err
	}
	defer r.Body.Close()

	// unmarshalling the body

	d := json.NewDecoder(strings.NewReader(string(body)))
	d.DisallowUnknownFields()
	err = d.Decode(&document)
	if err != nil {
		log.StatusPrintf(logger.ERROR, "%v", err)
		return document, err
	} else {
		log.StatusPrintln(logger.SUCCESS, "Good json!")
		log.Printf("\n--- JSON ---\n%s\n--- END ----", string(body))
		log.Printf("\n--- Struct ---\n%v\n--- END ----", document)
	}
	return document, err
}

// Checks if a given birthday exists
func birthdayExists(s *Server, document db.GuildUserPair) (bool, error) {
	log.StatusPrintln(logger.IN_PROGRESS, "Checking if birthday exists...")

	// creating filter
	filter := bson.M{
		"guilduserpair.guildid": document.GuildId,
		"guilduserpair.userid":  document.UserId,
	}

	// checking if birthday exists
	_, err := s.Database.FindBirthday(context.TODO(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// adds a birthday to the db if it DNE
func (s *Server) insertBirthday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.StatusPrintf(logger.ERROR, "%s", http.StatusText(http.StatusMethodNotAllowed))
		w.Header().Set("Allow", http.MethodPut)
		err := buildAndSendResponse(w, http.StatusMethodNotAllowed, ERROR, http.StatusText(http.StatusMethodNotAllowed), "")
		sendFallbackIfError(w, err)
		return
	}

	birthdayDocument, err := parseRequestBody[db.BirthdayDocument](r)
	if err != nil {
		err := buildAndSendResponse(w, http.StatusBadRequest, ERROR, "Failed to parse body of request", "")
		sendFallbackIfError(w, err)
		return
	}

	exists, err := birthdayExists(s, birthdayDocument.GuildUserPair)
	if err != nil {
		err := buildAndSendResponse(w, http.StatusInternalServerError, ERROR, "Failed when checking db for birthday", "")
		sendFallbackIfError(w, err)
		return
	}

	if exists {
		err := buildAndSendResponse(w, http.StatusOK, CONFLICT, "Birthday already exists", "")
		sendFallbackIfError(w, err)
		return
	}

	err = s.Database.InsertBirthday(context.TODO(), birthdayDocument)
	if err != nil {
		err = buildAndSendResponse(w, http.StatusInternalServerError, ERROR, err.Error(), "")
		sendFallbackIfError(w, err)
		return
	}

	err = buildAndSendResponse(w, http.StatusOK, SUCCESS, "Successfully added birthday", "")
	sendFallbackIfError(w, err)
}

// deletes a birthday from the database for a server if it exists
func (s *Server) deleteBirthday(w http.ResponseWriter, r *http.Request) {
	// checking for correct http method
	if r.Method != http.MethodDelete {
		log.StatusPrintf(logger.ERROR, "%s", http.StatusText(http.StatusMethodNotAllowed))
		w.Header().Set("Allow", http.MethodDelete)
		err := buildAndSendResponse(w, http.StatusMethodNotAllowed, ERROR, http.StatusText(http.StatusMethodNotAllowed), "")
		sendFallbackIfError(w, err)
		return
	}

	// parsing body of request
	guildUserPair, err := parseRequestBody[db.GuildUserPair](r)
	if err != nil {
		err = buildAndSendResponse(w, http.StatusBadRequest, ERROR, "Failed to parse body of request", "")
		sendFallbackIfError(w, err)
		return
	}

	exists, err := birthdayExists(s, guildUserPair)
	if err != nil {
		err = buildAndSendResponse(w, http.StatusInternalServerError, ERROR, "Failed when checking db for birthday", "")
		sendFallbackIfError(w, err)
		return
	}

	if !exists {
		err = buildAndSendResponse(w, http.StatusOK, CONFLICT, "Birthday does not exist", "")
		sendFallbackIfError(w, err)
		return
	}

	filter := bson.M{
		"guilduserpair.guildid": guildUserPair.GuildId,
		"guilduserpair.userid":  guildUserPair.UserId,
	}

	err = s.Database.DeleteBirthday(context.TODO(), filter)
	if err != nil {
		err = buildAndSendResponse(w, http.StatusInternalServerError, ERROR, err.Error(), "")
		sendFallbackIfError(w, err)
		return
	}

	err = buildAndSendResponse(w, http.StatusOK, SUCCESS, "Successfully deleted birthday", "")
	sendFallbackIfError(w, err)
}

func (s *Server) checkForBirthday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.StatusPrintf(logger.ERROR, "%s", http.StatusText(http.StatusMethodNotAllowed))
		w.Header().Set("Allow", http.MethodGet)
		err := buildAndSendResponse(w, http.StatusMethodNotAllowed, ERROR, http.StatusText(http.StatusMethodNotAllowed), "")
		sendFallbackIfError(w, err)
		return
	}

	params := r.URL.Query()
	guildId := params.Get("GuildId")
	userId := params.Get("UserId")
	log.Printf("Guild Id:\n%s", guildId)
	log.Printf("User Id:\n%s", userId)

	// creating filter
	filter := bson.M{
		"guilduserpair.guildid": guildId,
		"guilduserpair.userid":  userId,
	}

	// checking if birthday exists
	result, err := s.Database.FindBirthday(context.TODO(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = buildAndSendResponse(w, http.StatusOK, SUCCESS, "Birthday does not exist", "")
			sendFallbackIfError(w, err)
			return
		}
		err = buildAndSendResponse(w, http.StatusInternalServerError, ERROR, "Failed when checking db for birthday", "")
		sendFallbackIfError(w, err)
		return
	}
	bodyAsJson, err := json.Marshal(result)
	log.Printf("Marhsaled data:\n%s", string(bodyAsJson))
	if err != nil {
		err = buildAndSendResponse(w, http.StatusInternalServerError, ERROR, "Failed to marshal birthday", "")
		sendFallbackIfError(w, err)
		return
	}
	err = buildAndSendResponse(w, http.StatusOK, SUCCESS, "Birthday exists", string(bodyAsJson))
	sendFallbackIfError(w, err)
}

func (s *Server) getActiveBirthdays(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.StatusPrintf(logger.ERROR, "%s", http.StatusText(http.StatusMethodNotAllowed))
		w.Header().Set("Allow", http.MethodGet)
		buildAndSendResponse(w, http.StatusMethodNotAllowed, ERROR, http.StatusText(http.StatusMethodNotAllowed), "")
		return
	}

	// creating filter
	// filter := bson.M{
	// 	"guilduserpair.guildid": guildId,
	// 	"guilduserpair.userid":  userId,
	// }

	filter := bson.M{
		"day":   31,
		"month": 10,
	}

	result, err := s.Database.FindAllBirthdays(context.TODO(), filter)
	if err != nil {
		err = buildAndSendResponse(w, http.StatusInternalServerError, ERROR, err.Error(), "")
		sendFallbackIfError(w, err)
		return
	}

	if len(result) > 0 {
		bodyAsJson, err := json.Marshal(result)
		log.Printf("Marhsaled data:\n%s", string(bodyAsJson))
		if err != nil {
			err = buildAndSendResponse(w, http.StatusInternalServerError, ERROR, "Failed to marshal birthday", "")
			sendFallbackIfError(w, err)
			return
		}
		err = buildAndSendResponse(w, http.StatusOK, SUCCESS, "Found some active birthdays", string(bodyAsJson))
		sendFallbackIfError(w, err)
	} else {
		err = buildAndSendResponse(w, http.StatusOK, SUCCESS, "No active birthdays found", "")
		sendFallbackIfError(w, err)
	}

}

func (s *Server) updateBirthday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.StatusPrintf(logger.ERROR, "%s", http.StatusText(http.StatusMethodNotAllowed))
		w.Header().Set("Allow", http.MethodPut)
		buildAndSendResponse(w, http.StatusMethodNotAllowed, ERROR, http.StatusText(http.StatusMethodNotAllowed), "")
		return
	}

	birthdayDocument, err := parseRequestBody[db.BirthdayDocument](r)
	if err != nil {
		buildAndSendResponse(w, http.StatusBadRequest, ERROR, "Failed to parse body of request", "")
		return
	}

	exists, err := birthdayExists(s, birthdayDocument.GuildUserPair)
	if err != nil {
		buildAndSendResponse(w, http.StatusInternalServerError, ERROR, "Failed when checking db for birthday", "")
		return
	}

	if !exists {
		err := buildAndSendResponse(w, http.StatusOK, CONFLICT, "Birthday does not exist", "")
		sendFallbackIfError(w, err)
		return
	}

	err = s.Database.ReplaceBirthday(context.TODO(), birthdayDocument)
	if err != nil {
		err = buildAndSendResponse(w, http.StatusInternalServerError, ERROR, err.Error(), "")
		sendFallbackIfError(w, err)
		return
	}

	err = buildAndSendResponse(w, http.StatusOK, SUCCESS, "Successfully replaced birthday", "")
	sendFallbackIfError(w, err)
}

func main() {
	// loading .env
	loadDotEnv()

	// connecting to db
	client := connectToDB()

	// creating new instance of mongodb
	database := db.NewMongoDB(client, DATABASE_NAME, COLLECTION_NAME, log)

	// creating an instance of the server struct to pass the db onto the handlers
	server := &Server{
		Database: database,
	}

	// setting up handlers
	http.HandleFunc("/check-for-bday", server.checkForBirthday)
	http.HandleFunc("/get-active-bday", server.getActiveBirthdays)
	http.HandleFunc("/update-bday", server.updateBirthday)
	http.HandleFunc("/insert-bday", server.insertBirthday)
	http.HandleFunc("/delete-bday", server.deleteBirthday)

	// starting http server
	log.StatusPrintln(logger.IN_PROGRESS, "Starting http server")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.StatusFatalf(logger.ERROR, "%v", err)
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
	fmt.Println()
}

// {"guilduserpair": { "guildid": "some_server_id_2", "userid": "GenericUser" },"day": 31,"month": 10}
