package user

import (
	"context"
	"irl-mafia-game/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Interface that handlers/services depend on
type UserRepository interface {
	AddUser(context context.Context, user models.User) error
	FindUserWithUsername(context context.Context, username string) (models.User, error)
}

// Mongo implementation
type mongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) UserRepository {
	return &mongoRepository{collection: collection}
}

func (r *mongoRepository) AddUser(context context.Context, user models.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	_, err = r.collection.InsertOne(context, user)
	return err
}

func (r *mongoRepository) FindUserWithUsername(context context.Context, username string) (models.User, error) {
	var user models.User
	err := r.collection.FindOne(context, bson.M{"username": username}).Decode(&user)
	return user, err
}
