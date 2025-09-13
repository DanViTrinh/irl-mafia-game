package game

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Requests
type CreateGameRequest struct {
	PlayerIDs []string `json:"playerIds"`
	BoardSize int      `json:"boardSize"`
}

type JoinGameRequest struct {
	PlayerID string `json:"playerId"`
}

type ActionRequest struct {
	TargetID string `json:"targetId"`
	Action   string `json:"action"`
}

// Handlers

// CreateGameHandler godoc
// @Summary Create a new game
// @Description Create a new game with a list of player IDs and board size
// @Tags games
// @Accept json
// @Produce json
// @Param game body CreateGameRequest true "Game info"
// @Success 200 {object} map[string]interface{}
// @Router /games [post]
// @Security BearerAuth
func CreateGameHandler(repo GameRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateGameRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var players []primitive.ObjectID
		for _, id := range req.PlayerIDs {
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
				return
			}
			players = append(players, objID)
		}

		game := Game{
			Players:   players,
			BoardSize: req.BoardSize,
			Status:    "active",
			CreatedAt: time.Now(),
		}

		insertedID, err := repo.Create(context.Background(), game)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"gameId": insertedID.Hex()})
	}
}

// JoinGameHandler godoc
// @Summary Join an existing game
// @Description Add a player to an existing game by game ID
// @Tags games
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param game body JoinGameRequest true "Join info"
// @Success 200 {object} map[string]string
// @Router /games/{id}/join [post]
// @Security BearerAuth
func JoinGameHandler(repo GameRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameID := c.Param("id")
		var req JoinGameRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		gameObjID, err := primitive.ObjectIDFromHex(gameID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
			return
		}

		playerObjID, err := primitive.ObjectIDFromHex(req.PlayerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player ID"})
			return
		}

		if err := repo.AddPlayer(context.Background(), gameObjID, playerObjID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "joined"})
	}
}

// GetGameHandler godoc
// @Summary Get game details
// @Description Retrieve game details by game ID
// @Tags games
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} Game
// @Router /games/{id} [get]
// @Security BearerAuth
func GetGameHandler(repo GameRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		gameID := c.Param("id")
		gameObjID, err := primitive.ObjectIDFromHex(gameID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
			return
		}

		game, err := repo.GetByID(context.Background(), gameObjID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
			return
		}

		c.JSON(http.StatusOK, game)
	}
}

// TODO: Implement ActionHandler and other game-related handlers
