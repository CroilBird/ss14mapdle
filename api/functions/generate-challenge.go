package functions

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"ss14mapdle/db"
	"ss14mapdle/models"
	"ss14mapdle/util"
	"time"
)

var GenerateChallengeFunctionName FunctionName = "generate-challenge"

func init() {
	AddFunction(GenerateChallengeFunctionName, generateChallenge, "Generate a new ss14 mapdle challenge")
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func generateChallenge(params []string) {
	slog.Info("Generating SS14 mapdle challenge")

	dbConnection, err := db.Connect()

	if err != nil {
		slog.Error("[main] Could not connect to database", "error", err.Error())
		return
	}

	randomMap, err := models.GetRandomMap(dbConnection)

	if err != nil {
		slog.Error("Could not select random map", "error", err.Error())
		return
	}

	slog.Info(fmt.Sprintf("Selected random map: %s", randomMap.Name))

	minX := util.StartWidth / 2
	maxX := int(randomMap.Width - util.StartWidth/2)
	minY := util.StartHeight / 2
	maxY := int(randomMap.Height - util.StartHeight/2)

	x := randRange(minX, maxX)
	y := randRange(minY, maxY)

	challenge := models.Challenge{
		X:           x,
		Y:           y,
		MapID:       randomMap.ID,
		GeneratedAt: time.Now(),
	}

	err = dbConnection.Create(&challenge).Error

	if err != nil {
		slog.Error("Could not create challenge", "error", err.Error())
		return
	}
}
