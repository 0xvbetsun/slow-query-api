package configs

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const appName = "Slow-QUeery-API"

// Config represents app's configuration.
type Config struct {
	AppName           string
	Port              string
	PostgresDSN       string
	SlowQueryDuration uint64
}

// New constructor for application config
func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("application PORT not provided")
	}
	dbPort := os.Getenv("PG_PORT")
	dbHost := os.Getenv("PG_HOST")
	dbUser := os.Getenv("PG_USER")
	dbPass := os.Getenv("PG_PASSWORD")
	dbName := os.Getenv("PG_DATABASE")

	if dbPort == "" || dbHost == "" || dbUser == "" || dbPass == "" || dbName == "" {
		return nil, errors.New("database env not provided")
	}

	slowQuery, err := strconv.ParseUint(os.Getenv("SLOW_QUERY_DURATION"), 10, 64)

	if err != nil {
		return nil, err
	}

	dbDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName,
	)
	return &Config{
		AppName:           appName,
		Port:              port,
		PostgresDSN:       dbDSN,
		SlowQueryDuration: slowQuery,
	}, nil
}
