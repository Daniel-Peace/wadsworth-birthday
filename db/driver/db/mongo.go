package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

/*
 * Creates an instance of a mongodb implementation of the BirthdayDatabase interface
 */
func NewMongoDB(client *mongo.Client, database string, collection string) *MongoDB {
	return &MongoDB{
		Client:     client,
		Database:   database,
		Collection: collection,
	}
}

/*
 * Inserts a birthday document
 */
func (m *MongoDB) InsertOne(context context.Context, doc BirthdayDocument) error {
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	result, err := coll.InsertOne(context, doc)
	if err != nil {
		fmt.Println("Failed to insert doc:", err)
	} else {
		fmt.Println("Inserted doc with _id:", result.InsertedID)
	}
	return err
}

/*
 * Deletes a birthday document
 */
func (m *MongoDB) DeleteOne(context context.Context, filter bson.M) error {
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	result, err := coll.DeleteOne(context, filter)
	if err != nil {
		fmt.Printf("Failed to delete doc:%v", err)
	} else {
		fmt.Printf("Deleted doc: Count %d", result.DeletedCount)
	}
	return err
}

/*
 * Finds the first doc matching the filter
 */
func (m *MongoDB) FindOne(context context.Context, filter bson.M) (BirthdayDocument, error) {
	opts := options.FindOne()
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	var birthday BirthdayDocument
	err := coll.FindOne(context, filter, opts).Decode(&birthday)
	if err != nil {
		fmt.Printf("No birthday found: %v", err)
	} else {
		fmt.Printf("Retrieved %s's birthday from guild %s.\n",
			birthday.GuildUserPair.UserId,
			birthday.GuildUserPair.GuildId,
		)
	}
	return birthday, err
}

/*
 * Finds all docs matching the filter
 */
func (m *MongoDB) FindAll(context context.Context, filter bson.M) ([]BirthdayDocument, error) {
	coll := m.Client.Database(m.Database).Collection(m.Collection)
	cursor, err := coll.Find(context, filter)
	if err != nil {
		panic(err)
	}
	var results []BirthdayDocument
	if err = cursor.All(context, &results); err != nil {
		panic(err)
	}
	fmt.Printf(
		"Found %d active birthday(s)",
		len(results),
	)
	return results, err
}
