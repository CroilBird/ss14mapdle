package endpoints

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"ss14mapdle/models"
	"ss14mapdle/util"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(challengeEndpointGenerator)
}

func challengeEndpointGenerator(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/challenge/:zoom",
				Handler: func(c *gin.Context) {
					getCurrentChallenge(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/guess",
				Handler: func(c *gin.Context) {
					guess(db, c)
				},
			},
		},
	}
}

func getCurrentChallenge(db *gorm.DB, c *gin.Context) {
	zoomString := c.Param("zoom")

	zoom, err := strconv.Atoi(zoomString)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse zoom level"})
		return
	}

	var currentChallenge models.Challenge

	err = db.Table("challenges").
		Preload("Map").
		Last(&currentChallenge).
		Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get last challenge"})
		return
	}

	mapPartCachePath, err := util.GetMapPathAtLevel(
		currentChallenge.Map.Name,
		currentChallenge.Map.Path,
		currentChallenge.Map.Index,
		currentChallenge.X,
		currentChallenge.Y,
		zoom,
	)

	if err != nil {
		slog.Error("Could not get map image", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get map image"})
		return
	}

	c.File(*mapPartCachePath)
}

func guess(db *gorm.DB, c *gin.Context) {
	type GuessRequest struct {
		Name string `json:"name"`
	}

}
