package db

import (
	"context"
	"fmt"

	"github.com/Daniel-Peace/go-logger"
)

const (
	DB_NAME             = "wadsworth"
	BIRTHDAY_COLLECTION = "birthdays"
	CONFIG_COLLECTION   = "configs"
)

type WadsworthDB struct {
	birthdayCollection string
	configCollection   string
	birthdayDatabase   *MongoDatabase[BirthdayDocument]
	configDatabase     *MongoDatabase[GuildConfigDocument]
	log                *logger.GoLogger
}

func NewWadsworthDB(log *logger.GoLogger) *WadsworthDB {
	client, err := GetDBClient(log)
	if err != nil {
		log.StatusFatalf(logger.ERROR, "%v", err)
	}

	birthdayDatabase := NewMongoDatabase[BirthdayDocument](client, DB_NAME, log)
	configDatabase := NewMongoDatabase[GuildConfigDocument](client, DB_NAME, log)

	return &WadsworthDB{
		birthdayCollection: BIRTHDAY_COLLECTION,
		configCollection:   CONFIG_COLLECTION,
		birthdayDatabase:   birthdayDatabase,
		configDatabase:     configDatabase,
		log:                log,
	}
}

type GuildId string

type GuildUserPair struct {
	GuildId GuildId
	UserId  string
}

type BirthdayDocument struct {
	GuildUserPair GuildUserPair
	Day           int
	Month         int
}

type GuildConfigDocument struct {
	GuildId   string
	RoleId    string
	ChannelId string
}

//=====================================================================================================
//                                          BIRTHDAY DOC METHODS
//=====================================================================================================

func (wadsworthDB *WadsworthDB) InsertBirthdayDocument(ctx context.Context, doc BirthdayDocument) error {
	err := wadsworthDB.birthdayDatabase.InsertDocument(ctx, wadsworthDB.birthdayCollection, doc)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return err
}

func (wadsworthDB *WadsworthDB) DeleteBirthdayDocument(ctx context.Context, filter Filter) error {
	err := wadsworthDB.birthdayDatabase.DeleteDocument(ctx, wadsworthDB.birthdayCollection, filter)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return err
}

func (wadsworthDB *WadsworthDB) RetrieveBirthdayDocument(ctx context.Context, filter Filter) (BirthdayDocument, error) {
	doc, err := wadsworthDB.birthdayDatabase.RetrieveDocument(ctx, wadsworthDB.birthdayCollection, filter)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return doc, err
}

func (wadsworthDB *WadsworthDB) RetrieveBirthdayDocuments(ctx context.Context, filter Filter) ([]BirthdayDocument, error) {
	docs, err := wadsworthDB.birthdayDatabase.RetrieveDocuments(ctx, wadsworthDB.birthdayCollection, filter)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return docs, err
}
func (wadsworthDB *WadsworthDB) UpdateBirthdayDocument(ctx context.Context, filter Filter, fieldUpdates FieldUpdates) error {
	err := wadsworthDB.birthdayDatabase.UpdateDocument(ctx, wadsworthDB.birthdayCollection, filter, fieldUpdates)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return err
}

func (wadsworthDB *WadsworthDB) CountBirthdayDocuments(ctx context.Context, filter Filter) (int64, error) {
	count, err := wadsworthDB.birthdayDatabase.CountDocuments(ctx, wadsworthDB.birthdayCollection, filter)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
		count = 0
	}
	return count, err
}

func (wadsworthDB *WadsworthDB) ReplaceBirthdayDocument(ctx context.Context, filter Filter, doc BirthdayDocument) error {
	err := wadsworthDB.birthdayDatabase.ReplaceDocument(ctx, wadsworthDB.birthdayCollection, filter, doc)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return err
}

// =====================================================================================================
//
//	CONFIG DOC METHODS
//
// =====================================================================================================
func (wadsworthDB *WadsworthDB) InsertConfigDocument(ctx context.Context, doc GuildConfigDocument) error {
	err := wadsworthDB.configDatabase.InsertDocument(ctx, wadsworthDB.configCollection, doc)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return err
}

func (wadsworthDB *WadsworthDB) DeleteConfigDocument(ctx context.Context, filter Filter) error {
	err := wadsworthDB.configDatabase.DeleteDocument(ctx, wadsworthDB.configCollection, filter)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return err
}

func (wadsworthDB *WadsworthDB) RetrieveConfigDocument(ctx context.Context, filter Filter) (GuildConfigDocument, error) {
	doc, err := wadsworthDB.configDatabase.RetrieveDocument(ctx, wadsworthDB.configCollection, filter)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return doc, err
}

func (wadsworthDB *WadsworthDB) RetrieveConfigDocuments(ctx context.Context, filter Filter) ([]GuildConfigDocument, error) {
	docs, err := wadsworthDB.configDatabase.RetrieveDocuments(ctx, wadsworthDB.configCollection, filter)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return docs, err
}
func (wadsworthDB *WadsworthDB) UpdateConfigDocument(ctx context.Context, filter Filter, fieldUpdates FieldUpdates) error {
	err := wadsworthDB.configDatabase.UpdateDocument(ctx, wadsworthDB.configCollection, filter, fieldUpdates)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return err
}

func (wadsworthDB *WadsworthDB) CountConfigDocuments(ctx context.Context, filter Filter) (int64, error) {
	count, err := wadsworthDB.configDatabase.CountDocuments(ctx, wadsworthDB.configCollection, filter)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return count, fmt.Errorf("failed config count for filter %s: %w", filter, err)
}

func (wadsworthDB *WadsworthDB) ReplaceConfigDocument(ctx context.Context, filter Filter, doc GuildConfigDocument) error {
	err := wadsworthDB.configDatabase.ReplaceDocument(ctx, wadsworthDB.configCollection, filter, doc)
	if err != nil {
		wadsworthDB.log.StatusPrintf(logger.ERROR, "%v", err)
	}
	return fmt.Errorf("failed config replace for guild %s: %w", doc.GuildId, err)
}
