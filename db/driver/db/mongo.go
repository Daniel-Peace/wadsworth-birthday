package db

import (
	"context"
	custom_utils "driver/utils"
	"fmt"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

/*
 * Log strings
 */
var (
	WORKING_STATUS = "WORKING"
	SUCCESS_STATUS = "SUCCESS"
	ERROR_STATUS   = "ERROR"
)

type MongoDB struct {
	Client     *mongo.Client
	Database   string
	Collection string
	Logger     *log.Logger
}

/*
 * Creates an instance of a mongodb implementation of the BirthdayDatabase interface
 */
func NewMongoDB(client *mongo.Client, database string, collection string, logger *log.Logger) *MongoDB {
	return &MongoDB{
		Client:     client,
		Database:   database,
		Collection: collection,
		Logger:     logger,
	}
}

/*
 * Inserts a birthday document
 */
func (m *MongoDB) InsertOne(context context.Context, doc BirthdayDocument) error {
	m.Logger.Printf("[%s] [%s]",
		custom_utils.ColorizeString("InsertOne", custom_utils.Magenta),
		custom_utils.ColorizeString(WORKING_STATUS, custom_utils.Yellow),
	)
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	result, err := coll.InsertOne(context, doc)
	if err != nil {
		m.Logger.Printf("[%s] [%s] - Failed to insert doc: %v",
			custom_utils.ColorizeString("InsertOne", custom_utils.Magenta),
			custom_utils.ColorizeString(ERROR_STATUS, custom_utils.Red),
			err,
		)
	} else {
		m.Logger.Printf("[%s] [%s] - Inserted doc with %s",
			custom_utils.ColorizeString("InsertOne", custom_utils.Magenta),
			custom_utils.ColorizeString(SUCCESS_STATUS, custom_utils.Green),
			custom_utils.ColorizeString(fmt.Sprintf("%s", result.InsertedID), custom_utils.Blue),
		)
	}
	return err
}

/*
 * Deletes a birthday document
 */
func (m *MongoDB) DeleteOne(context context.Context, filter bson.M) error {
	m.Logger.Printf("[%s] [%s]",
		custom_utils.ColorizeString("DeleteOne", custom_utils.Magenta),
		custom_utils.ColorizeString(WORKING_STATUS, custom_utils.Yellow),
	)
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	result, err := coll.DeleteOne(context, filter)
	if err != nil {
		m.Logger.Printf("[%s] [%s] - Failed to insert doc: %v",
			custom_utils.ColorizeString("DeleteOne", custom_utils.Magenta),
			custom_utils.ColorizeString(ERROR_STATUS, custom_utils.Red),
			err,
		)
	} else {
		m.Logger.Printf("[%s] [%s] - Deleted doc: Count %s",
			custom_utils.ColorizeString("DeleteOne", custom_utils.Magenta),
			custom_utils.ColorizeString(SUCCESS_STATUS, custom_utils.Green),
			custom_utils.ColorizeString(strconv.FormatInt(result.DeletedCount, 10), custom_utils.Blue),
		)
	}
	return err
}

/*
 * Finds the first doc matching the filter
 */
func (m *MongoDB) FindOne(context context.Context, filter bson.M) (BirthdayDocument, error) {
	m.Logger.Printf("[%s] [%s]",
		custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
		custom_utils.ColorizeString(WORKING_STATUS, custom_utils.Yellow),
	)
	opts := options.FindOne()
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	var birthday BirthdayDocument
	err := coll.FindOne(context, filter, opts).Decode(&birthday)
	if err != nil {
		m.Logger.Printf("[%s] [%s] - No birthday found: %v",
			custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
			custom_utils.ColorizeString(ERROR_STATUS, custom_utils.Red),
			err,
		)
	} else {
		m.Logger.Printf("[%s] [%s] - Retrieved %s's birthday from guild %s",
			custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
			custom_utils.ColorizeString(SUCCESS_STATUS, custom_utils.Green),
			custom_utils.ColorizeString(birthday.GuildUserPair.UserId, custom_utils.Blue),
			custom_utils.ColorizeString(birthday.GuildUserPair.GuildId, custom_utils.Blue),
		)
	}
	return birthday, err
}

/*
 * Finds all docs matching the filter
 */
func (m *MongoDB) FindAll(context context.Context, filter bson.M) ([]BirthdayDocument, error) {
	m.Logger.Printf("[%s] [%s]",
		custom_utils.ColorizeString("FindAll", custom_utils.Magenta),
		custom_utils.ColorizeString(WORKING_STATUS, custom_utils.Yellow),
	)
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	cursor, err := coll.Find(context, filter)
	if err != nil {
		m.Logger.Fatalf("[%s] [%s] - Something went wrong while checking for birthdays: %v",
			custom_utils.ColorizeString("FindAll", custom_utils.Magenta),
			custom_utils.ColorizeString(ERROR_STATUS, custom_utils.Red),
			err,
		)
	}
	var results []BirthdayDocument
	if err = cursor.All(context, &results); err != nil {
		m.Logger.Fatalf("[%s] [%s] - Failed to retreive found birthdays: %v",
			custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
			custom_utils.ColorizeString(ERROR_STATUS, custom_utils.Red),
			err,
		)
	}
	m.Logger.Printf("[%s] [%s] - Found %s active birthdays",
		custom_utils.ColorizeString("FindAll", custom_utils.Magenta),
		custom_utils.ColorizeString(SUCCESS_STATUS, custom_utils.Green),
		custom_utils.ColorizeString(strconv.Itoa(len(results)), custom_utils.Blue),
	)
	return results, err
}
