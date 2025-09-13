package main

import (
	"context"
	"irl-mafia-game/auth"
	"irl-mafia-game/db"
	"irl-mafia-game/game"
	"irl-mafia-game/user"
	"log"
	"time"

	_ "irl-mafia-game/docs" // Swagger docs

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	dbm, err := db.NewDBManager("mongodb://localhost:27017", "irl-mafia-game")
	if err != nil {
		log.Fatal(err)
	}
	defer dbm.Close(context.Background())

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:19006", "http://localhost:8081"}, // Expo dev servers
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/signup", user.SignupHandler(dbm.UserRepo))
	r.POST("/login", user.LoginHandler(dbm.UserRepo))

	// Protected routes
	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())

	// User routes
	protected.GET("/users", user.GetAllUsersHandler(dbm.UserRepo))
	protected.GET("/users/me", user.GetCurrentUserHandler(dbm.UserRepo))

	// Game routes
	protected.GET("/games", game.GetAllGamesHandler(dbm.GameRepo))
	protected.POST("/games/create", game.CreateGameHandler(dbm.GameRepo, dbm.UserRepo))
	protected.GET("/games/:id", game.GetGameHandler(dbm.GameRepo))
	protected.POST("/games/:id/join", game.JoinGameHandler(dbm.GameRepo, dbm.UserRepo))
	protected.GET("/games/:id/players", game.GetPlayersUsernamesHandler(dbm.GameRepo, dbm.UserRepo))

	r.Run(":8080")
}
