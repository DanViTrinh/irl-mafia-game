package main

import (
	"irl-mafia-game/controllers"
	"irl-mafia-game/db"

	"github.com/gin-gonic/gin"
)

func main() {
    client := db.Connect("mongodb://localhost:27017")
    playerCollection := client.Database("fluGame").Collection("players")
    gameCollection := client.Database("fluGame").Collection("games")

    r := gin.Default()

    r.POST("/games", controllers.CreateGame(gameCollection))
    r.GET("/games/:id", controllers.GetGame(gameCollection))
    r.POST("/games/:id/join", controllers.JoinGame(gameCollection))
    r.POST("/games/:id/action", controllers.PerformAction(playerCollection))

    r.Run(":8080")
}
