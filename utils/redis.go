package utils

import (
	"github.com/go-redis/redis"
)

func GetRedisClient(address string, db int, password string) (redisClient *redis.Client) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		DB:       db,
		Password: password,
	})
	return
}
