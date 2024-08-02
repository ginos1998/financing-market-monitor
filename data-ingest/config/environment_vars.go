package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const varsFile = ".env.ingest"
var envsMap map[string]string

func LoadEnvVars() error {
	err := godotenv.Load(varsFile)
	if err != nil {
		return err
	}

	envsMap = make(map[string]string)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
            envsMap[pair[0]] = pair[1]
        }
	}

	logrus.Info("Environment variables loaded successfully")
	
	return nil
}

func GetEnvVar(key string) string {
	return envsMap[key]
}


