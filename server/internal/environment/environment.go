package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var envInstance *Environment

func getEnvFileFromType(etype string) string {
	return ApiServiceOverrideEnvFile
}

func NewEnvironment(etype string) (*Environment, error) {
	if envInstance == nil {
		pwd, pwdErr := os.Getwd()
		if pwdErr != nil {
			panic(pwdErr)
		}

		envFile := getEnvFileFromType(etype)

		envInstance = &Environment{}

		env, err := godotenv.Read(filepath.Join(pwd, "/internal/environment"+envFile))
		if err != nil {
			fmt.Println("failed to read apiservice environment file")
			return nil, err
		}

		for key := range env {
			findEnvVAlue := os.Getenv(key)
			if findEnvVAlue != "" {
				env[key] = findEnvVAlue
			} else {
				log.Error().Msg("failed to get environment variable from system")
				return nil, fmt.Errorf("failed to get environment variable from system")
			}
		}

		envDefaults, err := godotenv.Read(filepath.Join(pwd, "/internal/environment"+ApiServiceEnvDefaultsFile))
		if err != nil {
			fmt.Println("failed to read apiservice default environment file")
			return nil, err
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

	return envInstance, nil
}

func GetEnvironment() *Environment {
	if envInstance == nil {
		panic("environment is not initialized")
	}

	return envInstance
}

func (env *Environment) GetEnv(key string) string {
	if env == nil {
		panic("environment not initialized")
	}

	value, exists := env.data[key]

	if !exists {
		panic(fmt.Errorf("%s not found in environment", key))
	}

	return value
}

func (env *Environment) GetStrEnv(key string) string {
	return env.GetEnv(key)
}

func (env *Environment) GetBoolEnv(key string) bool {
	boolValue, err := strconv.ParseBool(env.GetEnv(key))
	if err != nil {
		panic(fmt.Sprintf("invalid boolean value for %s in environment", key))
	}

	return boolValue
}

func (env *Environment) GetIntEnv(key string) int {
	intValue, err := strconv.Atoi(env.GetEnv(key))
	if err != nil {
		panic(fmt.Sprintf("invalid int value for %s in environment", key))
	}

	return intValue
}

func (env *Environment) GetDurationEnv(key string) time.Duration {
	durationValue, err := time.ParseDuration(env.GetEnv(key))
	if err != nil {
		panic(fmt.Sprintf("invalid duration value for %s in environment", key))
	}

	return durationValue
}

func (env *Environment) GetStringListEnv(key string) []string {
	return strings.Split(env.GetEnv(key), ",")
}

func (env *Environment) GetTrimmedStringListEnv(key string) []string {
	values := []string{}
	for _, value := range env.GetStringListEnv(key) {
		values = append(values, strings.TrimSpace(value))
	}

	return values
}
