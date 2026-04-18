// Package config provides ways to read the runtime configuration of the
// application using environment variables.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

type EnvVar struct {
	DefaultValue string
	Optional     bool
}

type EnvKey string

const (
	EnvServerPort  EnvKey = "SERVER_PORT"
	EnvDbDialect          = "DB_DIALECT"
	EnvDbPath             = "DB_PATH"
	EnvDbHost             = "DB_HOST"
	EnvDbPort             = "DB_PORT"
	EnvDbName             = "DB_DBNAME"
	EnvDbUser             = "DB_USER"
	EnvDbPass             = "DB_PASS"
	EnvMapBasePath        = "MAP_BASE_PATH"
)

// Envs is a map of environment variables and their values
// by default, it is initialized with the default values
var envs = map[EnvKey]string{
	EnvServerPort:  ":8000",
	EnvDbDialect:   "sqlite",
	EnvDbPath:      "/tmp/test.db",
	EnvDbHost:      "localhost",
	EnvDbPort:      "5432",
	EnvDbName:      "ss14mapdle",
	EnvDbUser:      "postgres",
	EnvDbPass:      "",
	EnvMapBasePath: "maps/",
}

type EnvValue string

const (
	DbDialectSqlite   EnvValue = "sqlite"
	DbDialectPostgres EnvValue = "postgres"
)

func GetConfig(key EnvKey) (val string, err error) {
	val, contains := envs[key]

	if !contains {
		errorMsg := fmt.Sprintf("Tried to retrieve nonexistent env item %s", key)

		slog.Error(errorMsg)
		return "", errors.New(errorMsg)
	}

	return val, nil
}

func OverrideEnv(key EnvKey, val string) {
	slog.Warn(fmt.Sprintf("[env] Overriding %s with %s", string(key), val))

	envs[key] = val
}

// init initializes the Envs map with the environment variables,
// replacing any default values with actual values from the system environment.
func init() {
	slog.Info("[env] Environment variables:")

	// replace with actual values from system environment
	for name := range envs {
		sysValue := os.Getenv(string(name))
		// if the environment variable is not set, use default value
		if sysValue == "" {
			slog.Info(fmt.Sprintf("%s: %s (DEFAULT)", name, envs[name]))
			continue
		}

		envs[name] = sysValue

		slog.Info(fmt.Sprintf("%s: %s", name, envs[name]))
	}
}

// TODO: Add a function to validate the environment variables
// TODO: Add a bool/check to see if the environment variables are initialized
// necessary to prevent the application from running if the environment variables are not set
