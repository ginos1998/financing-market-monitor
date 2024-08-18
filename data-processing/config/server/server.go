package server

import (
	"errors"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/sirupsen/logrus"
)

type Server struct {
	EnvVars         map[string]string
	Logger          *logrus.Logger
	MongoRepository mongod.MongoRepository
}

func NewServer() *Server {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
	})

	envVars, err := loadEnvVars()
	if err != nil {
		logger.Fatal("Error loading environment variables: ", err)
	}

	mongoRepository, err := mongod.CreateMongoClient(envVars)
	if err != nil {
		logger.Fatal("Error creating MongoDB client: ", err)
	}

	return &Server{
		EnvVars:         envVars,
		Logger:          logger,
		MongoRepository: *mongoRepository,
	}
}

func loadEnvVars() (map[string]string, error) {
	envVars, err := LoadEnvVars()
	if err != nil {
		return nil, errors.New("error loading environment variables: " + err.Error())
	}

	return envVars, nil
}
