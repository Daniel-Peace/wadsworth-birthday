package db

// import (
// 	"context"
// 	custom_utils "driver/utils"
// 	"strconv"

// 	"github.com/Daniel-Peace/go-logger"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	"go.mongodb.org/mongo-driver/v2/mongo/options"
// )

// // type BirthdayMongoDB struct {
// // 	Client             *mongo.Client
// // 	Database           string
// // 	BirthdayCollection string
// // 	ConfigCollection   string
// // 	Logger             *logger.GoLogger
// // }

// /*
//  * Creates an instance of a mongodb implementation of the BirthdayDatabase interface
//  */
// // func NewMongoDB(client *mongo.Client, database string, birthdayCollection string, configCollection string, logger *logger.GoLogger) *BirthdayMongoDB {
// // 	return &BirthdayMongoDB{
// // 		Client:             client,
// // 		Database:           database,
// // 		BirthdayCollection: birthdayCollection,
// // 		ConfigCollection:   configCollection,
// // 		Logger:             logger,
// // 	}
// // }

// // Inserts the given birthday document
// func (wadsworthDB *WadsworthDB) InsertBirthday(context context.Context, doc BirthdayDocument) error {
// 	wadsworthDB.Logger.StatusPrintln(logger.IN_PROGRESS, "Inserting document...")
// 	coll := wadsworthDB.Client.Database(wadsworthDB.Database).Collection(wadsworthDB.BirthdayCollection)
// 	result, err := coll.InsertOne(context, doc)
// 	if err != nil {
// 		wadsworthDB.Logger.StatusPrintf(logger.ERROR, "Failed to insert doc: %v", err)
// 	} else {
// 		wadsworthDB.Logger.StatusPrintf(logger.SUCCESS, "Inserted doc with id \"%s\"", result.InsertedID)
// 	}
// 	return err
// }

// // Deletes the first document that matches the provided filter
// func (wadsworthDB *WadsworthDB) DeleteBirthday(context context.Context, filter bson.M) error {
// 	wadsworthDB.Logger.StatusPrintln(logger.IN_PROGRESS, "Deleting document...")
// 	coll := wadsworthDB.Client.Database(wadsworthDB.Database).Collection(wadsworthDB.BirthdayCollection)
// 	result, err := coll.DeleteOne(context, filter)
// 	if err != nil {
// 		wadsworthDB.Logger.StatusPrintf(logger.ERROR, "Failed to insert doc: %v", err)
// 	} else {
// 		wadsworthDB.Logger.StatusPrintf(logger.SUCCESS, "Deleted doc: Count %d", result.DeletedCount)
// 	}
// 	return err
// }

// // // Retreives the first birthday that matches the filter
// // func (m *BirthdayMongoDB) FindBirthday(context context.Context, filter bson.M) (BirthdayDocument, error) {
// // 	m.Logger.Printf("[%s] [%s]",
// // 		custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
// // 		custom_utils.ColorizeString(custom_utils.WORKING_STATUS, custom_utils.Yellow),
// // 	)
// // 	opts := options.FindOne()
// // 	coll := m.Client.Database(m.Database).Collection(m.BirthdayCollection)
// // 	var birthday BirthdayDocument
// // 	err := coll.FindOne(context, filter, opts).Decode(&birthday)
// // 	if err != nil {
// // 		m.Logger.Printf("[%s] [%s] - No birthday found: %v",
// // 			custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
// // 			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
// // 			err,
// // 		)
// // 	} else {
// // 		m.Logger.Printf("[%s] [%s] - Retrieved %s's birthday from guild %s",
// // 			custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
// // 			custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
// // 			custom_utils.ColorizeString(birthday.GuildUserPair.UserId, custom_utils.Blue),
// // 			custom_utils.ColorizeString(birthday.GuildUserPair.GuildId, custom_utils.Blue),
// // 		)
// // 	}
// // 	return birthday, err
// // }

// // func (m *BirthdayMongoDB) ReplaceBirthday(context context.Context, doc BirthdayDocument) error {
// // 	m.Logger.StatusPrintln(logger.IN_PROGRESS, "Replacing document...")
// // 	coll := m.Client.Database(m.Database).Collection(m.BirthdayCollection)

// // 	filter := bson.M{
// // 		"guilduserpair.guildid": doc.GuildUserPair.GuildId,
// // 		"guilduserpair.userid":  doc.GuildUserPair.UserId,
// // 	}

// // 	result, err := coll.ReplaceOne(context, filter, doc)
// // 	if err != nil {
// // 		m.Logger.StatusPrintf(logger.ERROR, "Failed to replace doc: %v", err)
// // 	} else {
// // 		m.Logger.StatusPrintf(logger.SUCCESS, "Replaced %d docs", result.ModifiedCount)
// // 	}
// // 	return err
// // }

// /*
//  * Finds all docs matching the filter
//  */
// func (wadsworthDB *WadsworthDB) FindAllBirthdays(context context.Context, filter bson.M) ([]BirthdayDocument, error) {
// 	wadsworthDB.Logger.Printf("[%s] [%s]",
// 		custom_utils.ColorizeString("FindAll", custom_utils.Magenta),
// 		custom_utils.ColorizeString(custom_utils.WORKING_STATUS, custom_utils.Yellow),
// 	)
// 	coll := wadsworthDB.Client.Database(wadsworthDB.Database).Collection(wadsworthDB.BirthdayCollection)
// 	cursor, err := coll.Find(context, filter)
// 	if err != nil {
// 		wadsworthDB.Logger.Fatalf("[%s] [%s] - Something went wrong while checking for birthdays: %v",
// 			custom_utils.ColorizeString("FindAll", custom_utils.Magenta),
// 			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
// 			err,
// 		)
// 	}
// 	var results []BirthdayDocument
// 	if err = cursor.All(context, &results); err != nil {
// 		wadsworthDB.Logger.Fatalf("[%s] [%s] - Failed to retreive found birthdays: %v",
// 			custom_utils.ColorizeString("FindOne", custom_utils.Magenta),
// 			custom_utils.ColorizeString(custom_utils.ERROR_STATUS, custom_utils.Red),
// 			err,
// 		)
// 	}
// 	wadsworthDB.Logger.Printf("[%s] [%s] - Found %s active birthdays",
// 		custom_utils.ColorizeString("FindAll", custom_utils.Magenta),
// 		custom_utils.ColorizeString(custom_utils.SUCCESS_STATUS, custom_utils.Green),
// 		custom_utils.ColorizeString(strconv.Itoa(len(results)), custom_utils.Blue),
// 	)
// 	return results, err
// }

// func (wadsworthDB *WadsworthDB) InsertConfig(context context.Context, doc GuildConfig) error {
// 	wadsworthDB.Logger.StatusPrintln(logger.IN_PROGRESS, "Inserting document...")
// 	coll := wadsworthDB.Client.Database(wadsworthDB.Database).Collection(wadsworthDB.ConfigCollection)
// 	result, err := coll.InsertOne(context, doc)
// 	if err != nil {
// 		wadsworthDB.Logger.StatusPrintf(logger.ERROR, "Failed to insert document: %v", err)
// 	} else {
// 		wadsworthDB.Logger.StatusPrintf(logger.SUCCESS, "Inserted document with id \"%s\"", result.InsertedID)
// 	}
// 	return err
// }

// func (wadsworthDB *WadsworthDB) DeleteConfig(context context.Context, filter bson.M) error {
// 	wadsworthDB.Logger.StatusPrintln(logger.IN_PROGRESS, "Deleting document...")
// 	coll := wadsworthDB.Client.Database(wadsworthDB.Database).Collection(wadsworthDB.ConfigCollection)
// 	result, err := coll.DeleteOne(context, filter)
// 	if err != nil {
// 		wadsworthDB.Logger.StatusPrintf(logger.ERROR, "Failed to insert document: %v", err)
// 	} else {
// 		wadsworthDB.Logger.StatusPrintf(logger.SUCCESS, "Deleted document: Count %d", result.DeletedCount)
// 	}
// 	return err
// }

// func (wadsworthDB *WadsworthDB) FindConfig(context context.Context, filter bson.M) (GuildConfig, error) {
// 	wadsworthDB.Logger.StatusPrintln(logger.IN_PROGRESS, "Finding config document...")
// 	opts := options.FindOne()
// 	coll := wadsworthDB.Client.Database(wadsworthDB.Database).Collection(wadsworthDB.ConfigCollection)
// 	var config GuildConfig
// 	err := coll.FindOne(context, filter, opts).Decode(&config)
// 	if err != nil {
// 		wadsworthDB.Logger.StatusPrintf(logger.ERROR, "No config found: %v", err)
// 	} else {
// 		wadsworthDB.Logger.StatusPrintln(logger.IN_PROGRESS, "Retrieved config file")
// 	}
// 	return config, err
// }

// func (wadsworthDB *WadsworthDB) ReplaceConfig(context context.Context, doc GuildConfig) error {
// 	wadsworthDB.Logger.StatusPrintln(logger.IN_PROGRESS, "Replacing document...")
// 	coll := wadsworthDB.Client.Database(wadsworthDB.Database).Collection(wadsworthDB.BirthdayCollection)

// 	filter := bson.M{
// 		"guildid": doc.GuildId,
// 	}

// 	result, err := coll.ReplaceOne(context, filter, doc)
// 	if err != nil {
// 		wadsworthDB.Logger.StatusPrintf(logger.ERROR, "Failed to replace doc: %v", err)
// 	} else {
// 		wadsworthDB.Logger.StatusPrintf(logger.SUCCESS, "Replaced %d docs", result.ModifiedCount)
// 	}
// 	return err
// }
