package server

import (
	"fmt"
	"os"

	"github.com/hyqe/timber"
)

const (
	env_PORT  = "PORT"
	env_LEVEL = "LEVEL"
	env_MONGO = "MONGO"
)

type Config struct {
	Port     string
	Level    timber.Level
	MongoUri string
}

func (c Config) Addr() string {
	return fmt.Sprintf(":%v", c.Port)
}

func GetConfig() (Config, error) {
	mongoUri, err := getMongoUri()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get mongo uri: %v", err)
	}

	level, err := getLogLevel()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get log level: %v", err)
	}

	return Config{
		Port:     getPort(),
		Level:    level,
		MongoUri: mongoUri,
	}, nil
}

// getAddr gets the servers address to bind.
func getPort() string {
	PORT, ok := os.LookupEnv(env_PORT)
	if ok {
		return PORT
	}
	return "8080"
}

// getMongoUri gets the mongodb connection string uri.
func getMongoUri() (string, error) {
	MONGO, ok := os.LookupEnv(env_MONGO)
	if !ok {
		return "", fmt.Errorf("missing env: '%v'", env_MONGO)
	}
	return MONGO, nil
}

// getMongoUri gets the mongodb connection string uri.
func getLogLevel() (timber.Level, error) {
	LEVEL, ok := os.LookupEnv(env_LEVEL)
	if ok {
		return timber.ParseLevel(LEVEL), nil
	}
	return timber.DEBUG, nil
}
