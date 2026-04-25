package endpoints

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"ss14mapdle/models"
	"ss14mapdle/util"
	"strings"
	"time"
)

func init() {
	RegisterEndpointGenerator(challengeEndpointGenerator)
}

func challengeEndpointGenerator(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/challenge/:sessionGuid",
				Handler: func(c *gin.Context) {
					getCurrentChallengeSession(db, c)
				},
			},
			{
				Method: "GET",
				Path:   "/challenge/map/:sessionGuid",
				Handler: func(c *gin.Context) {
					getCurrentChallengeMap(db, c)
				},
			},
			{
				Method: "POST",
				Path:   "/guess/:sessionGuid",
				Handler: func(c *gin.Context) {
					guess(db, c)
				},
			},
		},
	}
}

func getCurrentChallengeSession(db *gorm.DB, c *gin.Context) {
	sessionGuid := c.Param("sessionGuid")

	var session models.Session

	session.Guid = sessionGuid

	err := db.Table("sessions").Where("guid = ?", sessionGuid).FirstOrCreate(&session).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create or get session for challenge"})
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

	if session.ChallengeID == 0 {
		session.ChallengeID = currentChallenge.Id
		err = db.Save(session).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not update session"})
			return
		}
	}

	// return some bullshit when the requested challenge is not the current one
	if session.ChallengeID != currentChallenge.Id {
		c.Status(http.StatusGone)
		return
	}

	var response struct {
		Session   models.Session `json:"session"`
		ExpiresAt time.Time      `json:"expires_at"`
		Message   string         `json:"message"`
	}

	response.Session = session
	response.ExpiresAt = currentChallenge.GeneratedAt.Add(time.Hour * 24)

	if response.Session.Correct {
		response.Message = fmt.Sprintf("Correct! The map is %s", currentChallenge.Map.Name)
	}

	c.JSON(http.StatusOK, response)
}

func getCurrentChallengeMap(db *gorm.DB, c *gin.Context) {
	sessionGuid := c.Param("sessionGuid")

	var session models.Session

	session.Guid = sessionGuid

	err := db.Table("sessions").Where("guid = ?", sessionGuid).FirstOrCreate(&session).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create or get session for challenge"})
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

	if session.ChallengeID == 0 {
		session.ChallengeID = currentChallenge.Id
		err = db.Save(session).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not update session"})
			return
		}
	}

	// return some bullshit when the requested challenge is not the current one
	if session.ChallengeID != currentChallenge.Id {
		c.Status(http.StatusGone)
		return
	}

	zoomLevel := session.ZoomLevel

	if session.Correct {
		zoomLevel = util.MaxZoomLevel
	}

	mapPartCachePath, err := util.GetMapPathAtLevel(
		currentChallenge.Map.Name,
		currentChallenge.Map.Path,
		currentChallenge.Map.Index,
		currentChallenge.X,
		currentChallenge.Y,
		zoomLevel,
	)

	if err != nil {
		slog.Error("Could not get map image", "error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get map image"})
		return
	}

	c.File(*mapPartCachePath)
}

func guess(db *gorm.DB, c *gin.Context) {
	sessionGuid := c.Param("sessionGuid")

	type GuessRequest struct {
		Name string `json:"name"`
	}

	request := GuessRequest{}

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not bind guess request"})
		return
	}

	var currentChallenge models.Challenge

	err = db.Transaction(func(tx *gorm.DB) error {
		var transactionError error

		transactionError = db.Table("challenges").Preload("Map").Last(&currentChallenge).Error

		if transactionError != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not find challenge"})
			return transactionError
		}

		var session models.Session

		transactionError = db.Where("guid = ?", sessionGuid).Find(&session).Error

		if transactionError != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not retrieve session"})
			return transactionError
		}

		if session.ID == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return errors.New("no session found")
		}

		if currentChallenge.Id != session.ChallengeID {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "session-challenge mismatch"})
			return errors.New("challenge session mismatch")
		}

		if session.ZoomLevel > util.MaxZoomLevel {
			c.JSON(http.StatusOK, gin.H{"correct": false, "message": "You can't guess again"})
			return errors.New("tried to guess beyond limit")
		}

		if strings.ToLower(currentChallenge.Map.Name) != strings.ToLower(request.Name) {
			newZoomLevel := session.ZoomLevel + 1

			transactionError = db.Table("sessions").
				Where("guid = ?", sessionGuid).
				Update("zoom_level", newZoomLevel).
				Error

			if transactionError != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not update session"})
				return transactionError
			}

			if newZoomLevel <= util.MaxZoomLevel-1 {
				c.JSON(http.StatusOK, gin.H{"correct": false, "message": "Wrong! Try again", "guesses_remaining": util.MaxZoomLevel - newZoomLevel + 1})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"correct": false,
					"message": fmt.Sprintf("Wrong! The correct map was: %s. Try again tomorrow", currentChallenge.Map.Name),
				})
			}

			return nil
		}

		session.Correct = true

		transactionError = db.Save(session).Error

		if transactionError != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not save successful guess"})
			return transactionError
		}

		c.JSON(http.StatusOK, gin.H{
			"correct": true,
			"message": fmt.Sprintf("Correct! The map is %s", currentChallenge.Map.Name),
		})

		return nil
	})

	if err != nil {
		slog.Error(err.Error())
	}

	return
}
