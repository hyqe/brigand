package server

import (
	"fmt"
	"os"

	"github.com/hyqe/timber"
)

const (
	env_PORT        = "PORT"
	env_LEVEL       = "LEVEL"
	env_MONGO       = "MONGO"
	env_ACCESS_KEY  = "ACCESS_KEY"
	env_SECRET_KEY  = "SECRET_KEY"
	env_S3_ENDPOINT = "S3_ENDPOINT"
	env_REGION      = "REGION"
	env_BUCKET      = "BUCKET"
)

type Config struct {
	Port        string
	Level       timber.Level
	MongoUri    string
	Access_key  string
	Secret_key  string
	S3_endpoint string
	Region      string
	Bucket      string
}

func (c Config) Addr() string {
	return fmt.Sprintf(":%v", c.Port)
}

func getEnv(env_name string) (string, error) {
	env, ok := os.LookupEnv(env_name)
	if !ok {
		return "", fmt.Errorf("There is an error! No %s!!!!", env_name)
	}

	return env, nil
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

	access_key, err := getEnv(env_ACCESS_KEY)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get access key: %v", err)
	}

	secret_key, err := getEnv(env_SECRET_KEY)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get secret key: %v", err)
	}

	s3_endpoint, err := getEnv(env_S3_ENDPOINT)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get S3_endpoint: %v", err)
	}

	region, err := getEnv(env_REGION)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get bucket_name: %v", err)
	}

	bucket, err := getEnv(env_BUCKET)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get bucket_name: %v", err)
	}

	return Config{
		Port:        getPort(),
		Level:       level,
		MongoUri:    mongoUri,
		Access_key:  access_key,
		Secret_key:  secret_key,
		S3_endpoint: s3_endpoint,
		Region:      region,
		Bucket:      bucket,
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
