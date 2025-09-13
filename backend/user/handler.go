// internal/user/handler.go
package user

import (
	"irl-mafia-game/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// SignupHandler godoc
// @Summary Signup a new user
// @Description Create a new user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body SignupRequest true "User info"
// @Success 200 {object} map[string]string
// @Router /signup [post]
func SignupHandler(repo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignupRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := User{
			Username: req.Username,
			Password: req.Password,
		}

		err := repo.AddUser(c.Request.Context(), user)
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "signup successful"})
	}
}

// LoginHandler godoc
// @Summary Login a user
// @Description Authenticate user and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param user body LoginRequest true "User info"
// @Success 200 {object} LoginResponse
// @Router /login [post]
func LoginHandler(repo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repo.FindUserWithUsername(c.Request.Context(), req.Username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token, err := auth.GenerateToken(user.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Token:    token,
		})
	}
}

// GetAllUsersHandler godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
// @Security BearerAuth
func GetAllUsersHandler(repo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := repo.GetAllUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

// GetCurrentUserHandler godoc
// @Summary Get current user
// @Description Retrieve details of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UserResponse
// @Router /users/me [get]
// @Security BearerAuth
func GetCurrentUserHandler(repo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID not found in context"})
			return
		}

		userObjectId, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
			return
		}

		user, err := repo.FindUserWithID(c.Request.Context(), userObjectId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Games: func() []string {
				ids := make([]string, len(user.Games))
				for i, id := range user.Games {
					ids[i] = id.Hex()
				}
				return ids
			}(),
		})
	}
}
