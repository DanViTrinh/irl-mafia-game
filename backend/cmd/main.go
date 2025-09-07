package main

import (
	"irl-mafia-game/controllers"
	"irl-mafia-game/db"

	"github.com/gin-gonic/gin"

	_ "irl-mafia-game/docs" // swagger generated docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	client := db.Connect("mongodb://localhost:27017")
	playerCollection := client.Database("fluGame").Collection("players")
	gameCollection := client.Database("fluGame").Collection("games")
	userCollection := client.Database("fluGame").Collection("users")

	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	r.POST("/users/register", controllers.RegisterUser(userCollection))
	r.POST("/users/login", controllers.LoginUser(userCollection))

	r.POST("/games", controllers.CreateGame(gameCollection))
	r.GET("/games/:id", controllers.GetGame(gameCollection))
	r.POST("/games/:id/join", controllers.JoinGame(gameCollection))
	r.POST("/games/:id/action", controllers.PerformAction(playerCollection))

	r.Run(":8080")
}
