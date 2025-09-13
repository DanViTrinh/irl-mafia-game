// internal/db/dbmanager.go
package db

import (
	"context"
	"fmt"
	"irl-mafia-game/game"
	"irl-mafia-game/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBManager struct {
	Client   *mongo.Client
	Database *mongo.Database
	UserRepo user.UserRepository
	GameRepo game.GameRepository
}

// NewDBManager connects to Mongo and sets up repositories
func NewDBManager(uri, dbName string) (*DBManager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping Mongo: %w", err)
	}

	db := client.Database(dbName)

	// Ensure unique index on username in users collection
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = db.Collection("users").Indexes().CreateOne(context.Background(), indexModel)

	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	return &DBManager{
		Client:   client,
		Database: db,
		UserRepo: user.NewMongoRepository(db.Collection("users")),
		GameRepo: game.NewMongoRepository(db.Collection("games")),
	}, nil
}

// Close disconnects the Mongo client
func (dbm *DBManager) Close(ctx context.Context) error {
	return dbm.Client.Disconnect(ctx)
}
