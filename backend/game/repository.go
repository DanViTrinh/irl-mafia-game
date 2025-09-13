package game

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameRepository interface {
	Create(ctx context.Context, g Game) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (Game, error)
	AddPlayer(ctx context.Context, gameID, playerID primitive.ObjectID) error
	GetAllGames(ctx context.Context) ([]Game, error)
}

type mongoRepository struct {
	col *mongo.Collection
}

func NewMongoRepository(col *mongo.Collection) GameRepository {
	return &mongoRepository{col: col}
}

func (r *mongoRepository) Create(ctx context.Context, g Game) (primitive.ObjectID, error) {
	res, err := r.col.InsertOne(ctx, g)
	if err != nil {
		return primitive.NilObjectID, err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return id, nil
}

func (r *mongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (Game, error) {
	var g Game
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&g)
	if err != nil {
		return g, errors.New("game not found")
	}
	return g, nil
}

func (r *mongoRepository) AddPlayer(ctx context.Context, gameID, playerID primitive.ObjectID) error {
	_, err := r.col.UpdateOne(
		ctx,
		bson.M{"_id": gameID},
		bson.M{"$addToSet": bson.M{"players": playerID}},
	)
	return err
}

func (r *mongoRepository) GetAllGames(ctx context.Context) ([]Game, error) {
	var games []Game
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	return games, nil
}
