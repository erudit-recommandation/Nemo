package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var envVariable *EnvVariable = nil

func GetConfig() EnvVariable {
	switch CONFIG_MODE {
	case PRODUCTION:
		return GetEnvVariableFromPath("env.json")
	case DEV_DOCKER:
		return GetEnvVariableFromPath("env_dev_docker.json")
	case DEV:
		return GetEnvVariableFromPath("env_dev.json")
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
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Panicf("Error when opening file: %v", err)
		}

		var payload EnvVariable
		err = json.Unmarshal(content, &payload)
		if err != nil {
			log.Panicf("Error during Unmarshal(): %v", err)
		}
		return payload
	}
	return *envVariable
}
