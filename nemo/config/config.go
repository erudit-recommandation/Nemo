package config

import (
	"os"

	"github.com/joho/godotenv"
)

var envVariable *EnvVariable = nil

func GetConfig() EnvVariable {
	switch CONFIG_MODE {
	case PRODUCTION:
		return GetEnvVariableFromPath(".env")
	case DEV_DOCKER:
		return GetEnvVariableFromPath(".env_dev_docker")
	case DEV:
		return GetEnvVariableFromPath(".env_dev")
	}
	return EnvVariable{}
}

func SetConfig(config *EnvVariable) {
	envVariable = config
}

func ClearEnvVariable() {
	envVariable = nil
}

func GetEnvVariableFromPath(path string) EnvVariable {
	if envVariable == nil {
		if err := godotenv.Load(path); err != nil {
			panic("was not able to load config check the current path in relation to the .env file")
		}
		return EnvVariable{
			Port:                    os.Getenv("PORT"),
			ArangoPort:              os.Getenv("ARANGO_PORT"),
			ArangoPassword:          os.Getenv("ARANGO_PASSWORD"),
			ArangoUsername:          os.Getenv("ARANGO_USERNAME"),
			ArangoDatabase:          os.Getenv("ARANGO_DATABASE"),
			ArangoArticleCollection: os.Getenv("ARANGO_ARTICLE_COLLECTION"),
			GemsimServiceAddr:       os.Getenv("GEMSIM_SERVICE_ADDR"),
		}
	}
	return *envVariable
}
