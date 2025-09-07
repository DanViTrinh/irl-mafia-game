package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"irl-mafia-game/models"
)

type CreateGameRequest struct {
	PlayerIDs []string `json:"playerIds" example:"[\"id1\",\"id2\"]"`
	BoardSize int      `json:"boardSize" example:"3"`
}

// CreateGame godoc
// @Summary Create a new game
// @Description Create a new game with a list of player IDs and board size
// @Tags games
// @Accept json
// @Produce json
// @Param game body CreateGameRequest true "Game info"
// @Success 200 {object} map[string]interface{}
// @Router /games [post]
func CreateGame(gameCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateGameRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var players []primitive.ObjectID
		for _, id := range req.PlayerIDs {
			print(id)
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
				return
			}
			print(objID.Hex())
			players = append(players, objID)
		}

		game := models.Game{
			Players:   players,
			BoardSize: req.BoardSize,
			Status:    "active",
			CreatedAt: time.Now(),
		}

		res, err := gameCollection.InsertOne(context.Background(), game)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"gameId": res.InsertedID})
	}
}

type JoinGameRequest struct {
	PlayerID string `json:"playerId" example:"player_object_id"`
}

// JoinGame godoc
// @Summary Join an existing game
// @Description Add a player to an existing game by game ID
// @Tags games
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param player body JoinGameRequest true "Player info"
// @Success 200 {object} map[string]string
// @Router /games/{id}/join [post]
func JoinGame(gameCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameID := c.Param("id")
		var req JoinGameRequest
		c.ShouldBindJSON(&req)

		gameObjID, _ := primitive.ObjectIDFromHex(gameID)
		playerObjID, _ := primitive.ObjectIDFromHex(req.PlayerID)

		gameCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": gameObjID},
			bson.M{"$addToSet": bson.M{"players": playerObjID}},
		)

		c.JSON(http.StatusOK, gin.H{"status": "joined"})
	}
}

// GetGame godoc
// @Summary Get current game state
// @Description Retrieve the current state of a game by ID, including all players and board size
// @Tags games
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} models.Game
// @Failure 404 {object} map[string]string
// @Router /games/{id} [get]
func GetGame(gameCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameID := c.Param("id")
		gameObjID, _ := primitive.ObjectIDFromHex(gameID)
		var game models.Game
		gameCollection.FindOne(context.Background(), bson.M{"_id": gameObjID}).Decode(&game)
		c.JSON(http.StatusOK, game)
	}
}

type ActionRequest struct {
	PlayerID string `json:"playerId" example:"player_object_id"`
	Action   string `json:"action" example:"claim"`
	TargetID string `json:"targetId" example:"target_player_id"`
}

// PerformAction godoc
// @Summary Perform a daily action
// @Description Player performs "claim" or "guess" on a target player
// @Tags games
// @Accept json
// @Produce json
// @Param action body ActionRequest true "Action info"
// @Success 200 {object} map[string]string
// @Router /games/{id}/action [post]
func PerformAction(playerCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			PlayerID string `json:"playerId"`
			Action   string `json:"action"`   // "claim" or "guess"
			TargetID string `json:"targetId"` // target player
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		playerObjID, _ := primitive.ObjectIDFromHex(req.PlayerID)
		targetObjID, _ := primitive.ObjectIDFromHex(req.TargetID)

		var player models.Player
		playerCollection.FindOne(context.Background(), bson.M{"_id": playerObjID}).Decode(&player)

		// Daily limit check
		if sameDay(player.LastAction) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Already performed action today"})
			return
		}

		// Update last action date
		player.LastAction = time.Now()
		playerCollection.UpdateOne(context.Background(),
			bson.M{"_id": playerObjID},
			bson.M{"$set": bson.M{"lastAction": player.LastAction}},
		)

		// Action logic
		switch req.Action {
		case "claim":
			handleClaim(playerCollection, playerObjID, targetObjID)
		case "guess":
			handleGuess(playerCollection, playerObjID, targetObjID)
		}

		c.JSON(http.StatusOK, gin.H{"status": "action performed"})
	}
}

// Helpers
func sameDay(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}

func handleClaim(playerCollection *mongo.Collection, playerID, targetID primitive.ObjectID) {
	var target models.Player
	playerCollection.FindOne(context.Background(), bson.M{"_id": targetID}).Decode(&target)

	if target.Cooties {
		// Transfer Cooties
		playerCollection.UpdateOne(context.Background(),
			bson.M{"_id": targetID},
			bson.M{"$set": bson.M{"cooties": false}},
		)
		playerCollection.UpdateOne(context.Background(),
			bson.M{"_id": playerID},
			bson.M{"$set": bson.M{"cooties": true}},
		)
	} else {
		// Gain a tile
		playerCollection.UpdateOne(context.Background(),
			bson.M{"_id": playerID},
			bson.M{"$inc": bson.M{"claimedCount": 1}},
		)
	}
}

func handleGuess(playerCollection *mongo.Collection, playerID, targetID primitive.ObjectID) {
	var target models.Player
	playerCollection.FindOne(context.Background(), bson.M{"_id": targetID}).Decode(&target)

	if target.Cooties {
		// Correct guess: gain 2 tiles
		playerCollection.UpdateOne(context.Background(),
			bson.M{"_id": playerID},
			bson.M{"$inc": bson.M{"claimedCount": 2}},
		)
	} else {
		// Incorrect guess: lose 1 random tile (simplified)
		playerCollection.UpdateOne(context.Background(),
			bson.M{"_id": playerID},
			bson.M{"$inc": bson.M{"claimedCount": -1}},
		)
	}
}
