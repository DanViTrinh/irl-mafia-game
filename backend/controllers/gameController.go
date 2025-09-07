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

func CreateGame(gameCollection *mongo.Collection) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            PlayerIDs []string `json:"playerIds"`
            BoardSize int      `json:"boardSize"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        var players []primitive.ObjectID
        for _, id := range req.PlayerIDs {
            objID, _ := primitive.ObjectIDFromHex(id)
            players = append(players, objID)
        }

        game := models.Game{
            Players:   players,
            BoardSize: req.BoardSize,
            Status:    "active",
            CreatedAt: time.Now(),
        }

        res, _ := gameCollection.InsertOne(context.Background(), game)
        c.JSON(http.StatusOK, gin.H{"gameId": res.InsertedID})
    }
}

func JoinGame(gameCollection *mongo.Collection) gin.HandlerFunc {
    return func(c *gin.Context) {
        gameID := c.Param("id")
        var req struct {
            PlayerID string `json:"playerId"`
        }
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

func GetGame(gameCollection *mongo.Collection) gin.HandlerFunc {
    return func(c *gin.Context) {
        gameID := c.Param("id")
        gameObjID, _ := primitive.ObjectIDFromHex(gameID)
        var game models.Game
        gameCollection.FindOne(context.Background(), bson.M{"_id": gameObjID}).Decode(&game)
        c.JSON(http.StatusOK, game)
    }
}

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
