package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hyqe/timber"
)

const (
	env_PORT                 = "PORT"
	env_LEVEL                = "LEVEL"
	env_MONGO                = "MONGO"
	env_DO_SPACES_ACCESS_KEY = "DO_SPACES_ACCESS_KEY"
	env_DO_SPACES_SECRET_KEY = "DO_SPACES_SECRET_KEY"
	env_DO_SPACES_ENDPOINT   = "DO_SPACES_ENDPOINT"
	env_DO_SPACES_REGION     = "DO_SPACES_REGION"
	env_DO_SPACES_BUCKET     = "DO_SPACES_BUCKET"
	env_SUDO_PASSWORD        = "SUDO_PASSWORD"
	env_SUDO_USERNAME        = "SUDO_USERNAME"
)

type Config struct {
	Port              string
	Level             timber.Level
	MongoUri          string
	DOSpacesAccessKey string
	DOSpacesSecretKey string
	DOSpacesEndpoint  string
	DOSpacesRegion    string
	DOSpacesBucket    string
	Sudo              Credentials
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

	doAccessKey, err := getEnv(env_DO_SPACES_ACCESS_KEY)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get access key: %v", err)
	}

	doSecretKey, err := getEnv(env_DO_SPACES_SECRET_KEY)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get secret key: %v", err)
	}

	doEndpoint, err := getEnv(env_DO_SPACES_ENDPOINT)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get S3_endpoint: %v", err)
	}

	doRegion, err := getEnv(env_DO_SPACES_REGION)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get bucket_name: %v", err)
	}

	doBucket, err := getEnv(env_DO_SPACES_BUCKET)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get bucket_name: %v", err)
	}

	sudoPassword, err := getEnv(env_SUDO_PASSWORD)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get sudo password: %v", err)
	}

	sudoUsername, err := getEnv(env_SUDO_USERNAME)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get sudo username: %v", err)
	}

	sudo := Credentials{
		Username: sudoUsername,
		Password: sudoPassword,
	}

	return Config{
		Port:              getPort(),
		Level:             level,
		MongoUri:          mongoUri,
		DOSpacesAccessKey: doAccessKey,
		DOSpacesSecretKey: doSecretKey,
		DOSpacesEndpoint:  doEndpoint,
		DOSpacesRegion:    doRegion,
		DOSpacesBucket:    doBucket,
		Sudo:              sudo,
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

type Credentials struct {
	Username string
	Password string
}

func SudoMiddlware(sudo Credentials) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			switch username, password, ok := r.BasicAuth(); {
			case username == sudo.Username && password == sudo.Password && ok:
				next.ServeHTTP(w, r)
				return
			default:
				w.Header().Set("WWW-Authenticate", "Basic")
				http.Error(w, "no authorization tokens", http.StatusUnauthorized)
			}

		}
	}
}
