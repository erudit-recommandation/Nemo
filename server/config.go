package server

import (
	"os"

	"github.com/joho/godotenv"
)

var envVariable *EnvVariable = nil

func GetConfig() EnvVariable {
	return GetEnvVariableFromPath(".env")
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
			Port:             os.Getenv("PORT"),
			TemplateRootPath: os.Getenv("TEMPLATE_ROOT_PATH"),
		}
	}
	return *envVariable
}

type EnvVariable struct {
	Port             string
	TemplateRootPath string
}
