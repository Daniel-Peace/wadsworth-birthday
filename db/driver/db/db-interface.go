package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type BirthdayDatabase interface {
	InsertOne(context context.Context, doc BirthdayDocument) error
	DeleteOne(context context.Context, filter bson.M) error
	FindOne(context context.Context, filter bson.M) error
	UpdateOne(context context.Context, filter bson.M) error
	FindAll(context context.Context, filter bson.M) error
}

type GuildUserPair struct {
	GuildId string
	UserId  string
}

type BirthdayDocument struct {
	GuildUserPair GuildUserPair
	Day           int
	Month         int
}
