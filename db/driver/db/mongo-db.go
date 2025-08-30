package db

import (
	"context"
	"fmt"

	"github.com/Daniel-Peace/go-logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDatabase[T any] struct {
	client       *mongo.Client
	databaseName string
	log          *logger.GoLogger
}

// Enforcing explicit interface implementation
var _ IMongoDatabase[any] = (*MongoDatabase[any])(nil)

func NewMongoDatabase[T any](client *mongo.Client, databaseName string, log *logger.GoLogger) *MongoDatabase[T] {
	db := &MongoDatabase[T]{
		client:       client,
		databaseName: databaseName,
		log:          log,
	}
	return db
}

func GetDBClient(log *logger.GoLogger) (*mongo.Client, error) {
	uri := "mongodb://localhost:8000"

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("connecting to DB: %w", err)
	}

	var result bson.M
	err = client.Database("admin").RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("connecting to DB: %w", err)
	} else {
		log.StatusPrintln(logger.SUCCESS, "Successfully connected to DB")
	}

	return client, err
}

func (m *MongoDatabase[T]) InsertDocument(ctx context.Context, collName string, doc T) error {
	m.log.StatusPrintln(logger.IN_PROGRESS, "Inserting document into database...")
	coll := m.client.Database(m.databaseName).Collection(collName)
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		m.log.StatusPrintf(logger.ERROR, "Failed to insert document: %v", err)
	}
	return err //fmt.Errorf("insert %s: %w", collName, err)
}

func (m *MongoDatabase[T]) DeleteDocument(ctx context.Context, collName string, filter Filter) error {
	m.log.StatusPrintln(logger.IN_PROGRESS, "Deleting document from database...")
	coll := m.client.Database(m.databaseName).Collection(collName)
	_, err := coll.DeleteOne(ctx, filter)
	return fmt.Errorf("delete %s: %w", collName, err)
}

func (m *MongoDatabase[T]) RetrieveDocument(ctx context.Context, collName string, filter Filter) (T, error) {
	m.log.StatusPrintln(logger.IN_PROGRESS, "Retreiving document from database...")
	coll := m.client.Database(m.databaseName).Collection(collName)
	var document T
	err := coll.FindOne(ctx, filter).Decode(&document)
	return document, fmt.Errorf("retrieve doc %s: %w", collName, err)
}

func (m *MongoDatabase[T]) RetrieveDocuments(ctx context.Context, collName string, filter Filter) ([]T, error) {
	m.log.StatusPrintln(logger.IN_PROGRESS, "Retreiving documents from database...")
	coll := m.client.Database(m.databaseName).Collection(collName)
	cursor, err := coll.Find(ctx, filter)
	var results []T
	err = cursor.All(ctx, &results)
	return results, fmt.Errorf("retrieve docs %s: %w", collName, err)
}
func (m *MongoDatabase[T]) UpdateDocument(ctx context.Context, collName string, filter Filter, fieldUpdates FieldUpdates) error {
	m.log.StatusPrintln(logger.IN_PROGRESS, "Updating document in database...")
	coll := m.client.Database(m.databaseName).Collection(collName)
	fieldUpdatesCommand := map[string]any{
		"$set": fieldUpdates,
	}
	_, err := coll.UpdateOne(ctx, filter, fieldUpdatesCommand)
	return fmt.Errorf("update %s: %w", collName, err)
}

func (m *MongoDatabase[T]) CountDocuments(ctx context.Context, collName string, filter Filter) (int64, error) {
	opts := options.Count()
	coll := m.client.Database(m.databaseName).Collection(collName)
	count, err := coll.CountDocuments(context.TODO(), filter, opts)
	if err != nil {
		m.log.StatusPrintf(logger.ERROR, "Failed to count documents: %v", err)
	}

	if err != nil {
		err = fmt.Errorf("count %s: %w", collName, err)
	}
	return count, err
}

func (m *MongoDatabase[T]) ReplaceDocument(ctx context.Context, collName string, filter Filter, doc T) error {
	coll := m.client.Database(m.databaseName).Collection(collName)
	_, err := coll.ReplaceOne(ctx, filter, doc)
	return fmt.Errorf("replace %s: %w", collName, err)
}
