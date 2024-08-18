package server

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const varsFile = ".env.processing"

func LoadEnvVars() (map[string]string, error) {
	err := godotenv.Load(varsFile)
	if err != nil {
		return nil, err
	}

	envsMap := make(map[string]string)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envsMap[pair[0]] = pair[1]
		}
	}

	return envsMap, nil
}
