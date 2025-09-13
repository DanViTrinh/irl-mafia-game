package game

import (
	"context"
	"irl-mafia-game/user"
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
// @Router /games/create [post]
// @Security BearerAuth
func CreateGameHandler(gameRepo GameRepository, userRepo user.UserRepository) gin.HandlerFunc {
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

		insertedID, err := gameRepo.Create(context.Background(), game)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update each user's game list
		for _, playerID := range players {
			err = userRepo.AddGameToUser(context.Background(), playerID, insertedID)
			if err != nil {
				println("Failed to add game to user:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
				return
			}
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
func JoinGameHandler(gameRepo GameRepository, userRepo user.UserRepository) gin.HandlerFunc {
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

		if err := gameRepo.AddPlayer(context.Background(), gameObjID, playerObjID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = userRepo.AddGameToUser(context.Background(), playerObjID, gameObjID)
		if err != nil {
			println("Failed to add game to user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
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

// GetAllGamesHandler godoc
// @Summary Get all games
// @Description Retrieve a list of all games
// @Tags games
// @Accept json
// @Produce json
// @Success 200 {array} Game
// @Router /games [get]
// @Security BearerAuth
func GetAllGamesHandler(repo GameRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		games, err := repo.GetAllGames(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, games)
	}
}

// TODO: Implement ActionHandler and other game-related handlers
