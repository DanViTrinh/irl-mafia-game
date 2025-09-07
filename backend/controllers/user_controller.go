package controllers

import (
	"context"
	"irl-mafia-game/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserDTO true "User info"
// @Success 200 {object} map[string]string
// @Router /users/register [post]
func RegisterUser(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input UserDTO

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := models.User{
			Username: input.Username,
			Password: hashedPassword,
		}

		// Suppose db is a MongoDB collection
		_, err = userCollection.InsertOne(c, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

// LoginUser godoc
// @Summary Login a user
// @Description Authenticate user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserDTO true "User info"
// @Success 200 {object} map[string]string
// @Router /users/login [post]
func LoginUser(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input UserDTO

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user models.User
		// TODO: fix context
		err := userCollection.FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&user)
		if err != nil || !CheckPasswordHash(input.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	}
}

// helpers

// HashPassword hashes a plain password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a plain password with a hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
