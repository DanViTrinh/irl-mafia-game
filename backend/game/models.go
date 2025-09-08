package game

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Game struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Players   []primitive.ObjectID `bson:"players"`
	BoardSize int                  `bson:"boardSize"`
	Status    string               `bson:"status"` // active, finished
	CreatedAt time.Time            `bson:"createdAt"`
}

type Tile struct {
	FriendID primitive.ObjectID `bson:"friendId"`
	Claimed  bool               `bson:"claimed"`
}

type Board struct {
	Tiles []Tile `bson:"tiles"`
}

type Player struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	PlayerName   string             `bson:"playerName"`
	User         primitive.ObjectID `bson:"user"`
	Board        []Tile             `bson:"board"`
	Cooties      bool               `bson:"cooties"`
	LastAction   time.Time          `bson:"lastAction"`
	ClaimedCount int                `bson:"claimedCount"`
}
