package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tile struct {
    FriendID primitive.ObjectID `bson:"friendId"`
    Claimed  bool               `bson:"claimed"`
}

type Player struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    Name          string             `bson:"name"`
    Email         string             `bson:"email"`
    Board         []Tile             `bson:"board"`
    Cooties       bool               `bson:"cooties"`
    LastAction    time.Time          `bson:"lastAction"`
    ClaimedCount  int                `bson:"claimedCount"`
}

type Game struct {
    ID        primitive.ObjectID   `bson:"_id,omitempty"`
    Players   []primitive.ObjectID `bson:"players"`
    BoardSize int                  `bson:"boardSize"`
    Status    string               `bson:"status"` // active, finished
    CreatedAt time.Time            `bson:"createdAt"`
}
