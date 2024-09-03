package redis

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(envVars map[string]string) (*RedisClient, error) {
	host := envVars["REDIS_HOST"]
	port := envVars["REDIS_PORT"]
	password := envVars["REDIS_PASSWORD"]
	username := envVars["REDIS_USERNAME"]
	db, _ := strconv.Atoi(envVars["REDIS_DB"])
	if host == "" || port == "" || password == "" || username == "" || db < 0 || db > 15 {
		return nil, errors.New("missing redis env vars")
	}

	address := host + ":" + port
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Username: username,
		Password: password,
		DB:       db,
	})
	return &RedisClient{Client: client}, nil
}
