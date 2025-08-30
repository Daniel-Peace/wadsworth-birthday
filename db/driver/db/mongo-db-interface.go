package db

import (
	"context"
)

type IMongoDatabase[T any] interface {
	InsertDocument(context context.Context, collName string, doc T) error
	DeleteDocument(context context.Context, collName string, filter Filter) error
	RetrieveDocument(context context.Context, collName string, filter Filter) (T, error)
	RetrieveDocuments(context context.Context, collName string, filter Filter) ([]T, error)
	UpdateDocument(context context.Context, collName string, filter Filter, fieldUpdates FieldUpdates) error
	CountDocuments(context context.Context, collName string, filter Filter) (int64, error)
	ReplaceDocument(context context.Context, collName string, filter Filter, doc T) error
}

type Filter map[string]any

type FieldUpdates map[string]any
