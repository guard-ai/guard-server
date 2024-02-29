package pkg

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ServerAddress            string
	ServerPort               string
	WorkerAuthToken          string
	PostgresConnectionString string
}

var env *EnvConfig = nil

func Env() *EnvConfig {
	if env == nil {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}

		env = &EnvConfig{
			ServerAddress:            os.Getenv("SERVER_ADDRESS"),
			ServerPort:               os.Getenv("SERVER_PORT"),
			WorkerAuthToken:          os.Getenv("WORKER_AUTH_TOKEN"),
			PostgresConnectionString: os.Getenv("POSTGRES_CONNECTION_STRING"),
		}
	}
	return env
}
