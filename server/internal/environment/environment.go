package environment

import (
	"fmt"

	"github.com/joho/godotenv"
)

var envInstance *Environment

func getEnvFileFromType(etype string) string {
	if etype == "dev" {
		return ApiServiceDevEnvFile
	}
	return ApiServiceProdEnvFile
}

func NewEnvironment(etype string) *Environment {
	if envInstance == nil {
		envFile := getEnvFileFromType(etype)

		if envInstance == nil {
			envInstance = &Environment{}

			env, err := godotenv.Read(envFile)
			if err != nil {
				fmt.Errorf("failed to read apiservice environment file %s", envFile)
				return nil
			}

			envDefaults, err := godotenv.Read(ApiServiceEnvDefaultsFile)
			if err != nil {
				fmt.Errorf("failed to read apiservice default environment file %s", ApiServiceEnvDefaultsFile)
				return nil
			} else {
				for key, defaultValue := range envDefaults {
					_, found := env[key]

					if !found {
						env[key] = defaultValue
					}
				}
			}

			envInstance.data = env
		}
	}

	return envInstance
}
