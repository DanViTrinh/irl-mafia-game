package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Interface that handlers/services depend on
type UserRepository interface {
	AddUser(context context.Context, user User) error
	FindUserWithID(context context.Context, id primitive.ObjectID) (User, error)
	FindUserWithUsername(context context.Context, username string) (User, error)
	GetAllUsers(context context.Context) ([]UserResponse, error)
	AddGameToUser(context context.Context, userID primitive.ObjectID, gameID primitive.ObjectID) error
}

// Mongo implementation
type mongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) UserRepository {
	return &mongoRepository{collection: collection}
}

func (r *mongoRepository) AddUser(context context.Context, user User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	_, err = r.collection.InsertOne(context, user)
	return err
}

func (r *mongoRepository) AddGameToUser(context context.Context, userID primitive.ObjectID, gameID primitive.ObjectID) error {
	// Convert string ID to ObjectID
	r.collection.UpdateByID(context, userID, bson.M{
		"$addToSet": bson.M{"friends": "asdflkj"},
	})

	result, err := r.collection.UpdateByID(context, userID, bson.M{
		"$addToSet": bson.M{"games": gameID},
	})

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *mongoRepository) FindUserWithID(context context.Context, id primitive.ObjectID) (User, error) {
	var user User
	err := r.collection.FindOne(context, bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (r *mongoRepository) FindUserWithUsername(context context.Context, username string) (User, error) {
	var user User
	err := r.collection.FindOne(context, bson.M{"username": username}).Decode(&user)
	return user, err
}

func (r *mongoRepository) GetAllUsers(context context.Context) ([]UserResponse, error) {
	var dbUsers []User
	cursor, err := r.collection.Find(context, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context, &dbUsers); err != nil {
		return nil, err
	}

	var responseUsers []UserResponse = make([]UserResponse, len(dbUsers))
	for i, user := range dbUsers {
		responseUsers[i] = UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
		}
	}

	return responseUsers, nil
}
