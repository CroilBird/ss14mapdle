// Package db provides a method to connect to a database using GORM.
package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log/slog"
	"ss14mapdle/config"
	"ss14mapdle/models"
)

func connectPostres() (*gorm.DB, error) {

	// behold, golang

	var err error

	var DbHost string
	var DbPort string
	var DbName string
	var DbUser string
	var DbPass string
	if DbHost, err = config.GetConfig(config.EnvDbHost); err != nil {
		slog.Error("[db] Could not get env variable for postgres host address", "error", err.Error())
		return nil, err
	}
	if DbPort, err = config.GetConfig(config.EnvDbPort); err != nil {
		slog.Error("[db] Could not get env variable for postgres port", "error", err.Error())
		return nil, err
	}
	if DbName, err = config.GetConfig(config.EnvDbName); err != nil {
		slog.Error("[db] Could not get env variable for postgres DB name", "error", err.Error())
		return nil, err
	}
	if DbUser, err = config.GetConfig(config.EnvDbUser); err != nil {
		slog.Error("[db] Could not get env variable for postgres user", "error", err.Error())
		return nil, err
	}
	if DbPass, err = config.GetConfig(config.EnvDbPass); err != nil {
		slog.Error("[db] Could not get env variable for postgres password", "error", err.Error())
		return nil, err
	}

	connectionUrl := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		DbHost,
		DbPort,
		DbName,
		DbUser,
		DbPass,
	)

	db, err := gorm.Open(postgres.Open(connectionUrl))

	if err != nil {
		slog.Error("[db] Could not open postgres connection", "error", err.Error())
	}

	return db, err
}

func connectSqlite() (*gorm.DB, error) {
	dbPath, err := config.GetConfig(config.EnvDbPath)

	if err != nil {
		slog.Error("[db] Could not get env variable for sqlite DB path", "error", err.Error())
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath))
	return db, err
}

// Connect uses GORM to connect to a database based on the environment variables.
// This function relies on the user to have set the environment variables correctly.
// The function will panic if the dialect is not supported.
func Connect() (*gorm.DB, error) {
	slog.Info("[db] Opening database connection")

	var db *gorm.DB
	var err error
	var DbDialect string

	if DbDialect, err = config.GetConfig(config.EnvDbDialect); err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("[db] Dialect: %s", DbDialect))

	// We want to be specific about what we are expecting to use
	// Postgres for production and SQLite for testing
	if DbDialect == string(config.DbDialectPostgres) {
		db, err = connectPostres()
	} else if DbDialect == string(config.DbDialectSqlite) {
		db, err = connectSqlite()
	} else {
		slog.Error(fmt.Sprintf("Unsupported database dialect: %s", DbDialect))
		panic(fmt.Sprintf("Unsupported database dialect: %s", DbDialect))
	}

	if err != nil {
		slog.Error(fmt.Sprintf("[db] Failed to connect to database: %s", err))
		panic(err)
	}

	slog.Info("[db] Database connection established")

	slog.Info("[db] Migrating models")
	models.Migrate(db)

	return db, nil
}
