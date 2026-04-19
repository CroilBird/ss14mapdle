package functions

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"ss14mapdle/config"
	"ss14mapdle/db"
	"ss14mapdle/endpoints"
	"ss14mapdle/models"
	"time"
)

var StartApiFunctionName FunctionName = "start-api"

func init() {
	AddFunction(StartApiFunctionName, startApi, "start the ss14 mapdle API")
}

func newChallengeGenerator(db *gorm.DB) {
	for {

		var lastChallenge models.Challenge
		var challengeCount int64

		err := db.Table("challenges").Count(&challengeCount).Last(&lastChallenge).Error

		if err != nil && challengeCount > 0 {
			slog.Error("[challenge worker] Failed to get last challenge")
			time.Sleep(1 * time.Minute)
			continue
		}

		if challengeCount == 0 {
			slog.Info("Generating new challenge")
			GenerateChallenge(db)
			slog.Info("Generated new challenge")
			continue
		}

		now := time.Now()

		if now.YearDay() == lastChallenge.GeneratedAt.YearDay() {
			time.Sleep(1 * time.Minute)
			continue
		}

		slog.Info("Generating new challenge")
		GenerateChallenge(db)
		slog.Info("Generated new challenge")
	}
}

func startApi(params []string) {
	slog.Info("Starting SS14 mapdle API")

	slog.Info("Setting up database")
	dbConnection, err := db.Connect()

	if err != nil {
		slog.Error("[main] Could not connect to database")
		panic("Panic in main function: Could not connect to database")
	}

	slog.Info("Setting up router")
	// setup router
	router := gin.Default()

	// of cors, we use cors
	router.Use(cors.Default())

	slog.Info("Registering endpoints")
	endpoints.RegisterEndpointHandlers(router, dbConnection)

	slog.Info("Migrating models")
	models.Migrate(dbConnection)

	slog.Info("Loading map JSON")
	err = models.LoadMapJson(dbConnection)

	if err != nil {
		slog.Error("Error loading map json", "error", err.Error())
		return
	}

	slog.Info("Starting server")

	serverPort, err := config.GetConfig("SERVER_PORT")

	if err != nil {
		slog.Error("[main] Could not get env variable for server port")
		panic("Panic in main function: could not start server")
	}

	slog.Info("Starting challenge generator goroutine")

	go newChallengeGenerator(dbConnection)

	err = router.Run(serverPort)
}
