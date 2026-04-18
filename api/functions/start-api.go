package functions

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"ss14mapdle/config"
	"ss14mapdle/db"
	"ss14mapdle/endpoints"
	"ss14mapdle/models"
)

var StartApiFunctionName FunctionName = "start-api"

func init() {
	AddFunction(StartApiFunctionName, startApi, "start the ss14 mapdle API")
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

	err = router.Run(serverPort)
}
