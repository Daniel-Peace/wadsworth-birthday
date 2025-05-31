package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type BirthdayDatabase interface {
	InsertBirthday(context context.Context, doc BirthdayDocument) error
	DeleteBirthday(context context.Context, filter bson.M) error
	FindBirthday(context context.Context, filter bson.M) (BirthdayDocument, error)
	ReplaceBirthday(context context.Context, doc BirthdayDocument) error
	FindAllBirthdays(context context.Context, filter bson.M) ([]BirthdayDocument, error)
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
