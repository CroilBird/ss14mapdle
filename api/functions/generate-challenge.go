package functions

import (
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"math/rand/v2"
	"ss14mapdle/db"
	"ss14mapdle/models"
	"ss14mapdle/util"
	"time"
)

var GenerateChallengeFunctionName FunctionName = "generate-challenge"

func init() {
	AddFunction(GenerateChallengeFunctionName, generateChallengeFunction, "Generate a new ss14 mapdle challenge")
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func GenerateChallenge(db *gorm.DB, forceMap string) {
	var selectedMap *models.Map
	var err error

	if forceMap == "" {
		selectedMap, err = models.GetRandomMap(db)
	} else {
		selectedMap, err = models.GetMap(db, forceMap)
	}

	if err != nil {
		slog.Error("Could not select map", "error", err.Error())
		return
	}

	slog.Info(fmt.Sprintf("Selected map: %s", selectedMap.Name))

	slog.Info(fmt.Sprintf("Map dimensions: %dx%d", selectedMap.Width, selectedMap.Height))

	minX := util.StartWidth / 2
	maxX := int(selectedMap.Width - util.StartWidth/2)
	minY := util.StartHeight / 2
	maxY := int(selectedMap.Height - util.StartHeight/2)

	slog.Info(fmt.Sprintf("Random bounds - Min: x%d y%d - Max x%d y%d", minX, minY, maxX, maxY))

	x := randRange(minX, maxX)
	y := randRange(minY, maxY)

	challenge := models.Challenge{
		X:           x,
		Y:           y,
		MapID:       selectedMap.ID,
		GeneratedAt: time.Now(),
	}

	err = db.Create(&challenge).Error

	if err != nil {
		slog.Error("Could not create challenge", "error", err.Error())
		return
	}
}

func generateChallengeFunction(params []string) {
	slog.Info("Generating SS14 mapdle challenge")

	dbConnection, err := db.Connect()

	if err != nil {
		slog.Error("[main] Could not connect to database", "error", err.Error())
		return
	}

	forceMap := ""

	if len(params) > 0 {
		forceMap = params[0]
	}

	GenerateChallenge(dbConnection, forceMap)
}
