package db

import (
	"context"
	custom_utils "driver/utils"
	"strconv"

	"github.com/Daniel-Peace/go-logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client     *mongo.Client
	Database   string
	Collection string
	Logger     *logger.GoLogger
}

/*
 * Creates an instance of a mongodb implementation of the BirthdayDatabase interface
 */
func NewMongoDB(client *mongo.Client, database string, collection string, logger *logger.GoLogger) *MongoDB {
	return &MongoDB{
		Client:     client,
		Database:   database,
		Collection: collection,
		Logger:     logger,
	}
}

// Inserts the given birthday document
func (m *MongoDB) InsertBirthday(context context.Context, doc BirthdayDocument) error {
	m.Logger.StatusPrintln(logger.IN_PROGRESS, "Inserting document...")
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	result, err := coll.InsertOne(context, doc)
	if err != nil {
		m.Logger.StatusPrintf(logger.ERROR, "Failed to insert doc: %v", err)
	} else {
		m.Logger.StatusPrintf(logger.SUCCESS, "Inserted doc with id \"%s\"", result.InsertedID)
	}
	return err
}

// Deletes the first document that matches the provided filter
func (m *MongoDB) DeleteBirthday(context context.Context, filter bson.M) error {
	m.Logger.StatusPrintln(logger.IN_PROGRESS, "Deleting document...")
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	result, err := coll.DeleteOne(context, filter)
	if err != nil {
		m.Logger.StatusPrintf(logger.ERROR, "Failed to insert doc: %v", err)
	} else {
		m.Logger.StatusPrintf(logger.SUCCESS, "Deleted doc: Count %d", result.DeletedCount)
	}
	return err
}

// Retreives the first birthday that matches the filter
func (m *MongoDB) FindBirthday(context context.Context, filter bson.M) (BirthdayDocument, error) {
	m.Logger.Printf("[%s] [%s]",
		custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
		custom_utils.ColorizeString(custom_utils.WORKING_STATUS, custom_utils.Yellow),
	)
	opts := options.FindOne()
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	var birthday BirthdayDocument
	err := coll.FindOne(context, filter, opts).Decode(&birthday)
	if err != nil {
		m.Logger.Printf("[%s] [%s] - No birthday found: %v",
			custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
			err,
		)
	} else {
		m.Logger.Printf("[%s] [%s] - Retrieved %s's birthday from guild %s",
			custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
			custom_utils.ColorizeString(birthday.GuildUserPair.UserId, custom_utils.Blue),
			custom_utils.ColorizeString(birthday.GuildUserPair.GuildId, custom_utils.Blue),
		)
	}
	return birthday, err
}

func (m *MongoDB) ReplaceBirthday(context context.Context, doc BirthdayDocument) error {
	m.Logger.StatusPrintln(logger.IN_PROGRESS, "Replacing document...")
	coll := m.Client.Database(m.Database).Collection(m.Collection)

	filter := bson.M{
		"guilduserpair.guildid": doc.GuildUserPair.GuildId,
		"guilduserpair.userid":  doc.GuildUserPair.UserId,
	}

	result, err := coll.ReplaceOne(context, filter, doc)
	if err != nil {
		m.Logger.StatusPrintf(logger.ERROR, "Failed to replace doc: %v", err)
	} else {
		m.Logger.StatusPrintf(logger.SUCCESS, "Replaced %d docs", result.ModifiedCount)
	}
	return err
}

/*
 * Finds all docs matching the filter
 */
func (m *MongoDB) FindAllBirthdays(context context.Context, filter bson.M) ([]BirthdayDocument, error) {
	m.Logger.Printf("[%s] [%s]",
		custom_utils.ColorizeString("FindAll", custom_utils.Magenta),
		custom_utils.ColorizeString(custom_utils.WORKING_STATUS, custom_utils.Yellow),
	)
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	cursor, err := coll.Find(context, filter)
	if err != nil {
		m.Logger.Fatalf("[%s] [%s] - Something went wrong while checking for birthdays: %v",
			custom_utils.ColorizeString("FindAll", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
			err,
		)
	}
	var results []BirthdayDocument
	if err = cursor.All(context, &results); err != nil {
		m.Logger.Fatalf("[%s] [%s] - Failed to retreive found birthdays: %v",
			custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
			err,
		)
	}
	m.Logger.Printf("[%s] [%s] - Found %s active birthdays",
		custom_utils.ColorizeString("FindAll", custom_utils.Magenta),
		custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
		custom_utils.ColorizeString(strconv.Itoa(len(results)), custom_utils.Blue),
	)
	return results, err
}
