package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty"`
	Username string               `bson:"username"`
	Password string               `bson:"password"`
	Games    []primitive.ObjectID `bson:"games,omitempty"`
}

type UserResponse struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Games    []string `json:"games,omitempty"`
}
