package config

import (
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisConfig(viper *viper.Viper) *redis.Client {

	config := redis.NewClient(&redis.Options{Addr: os.Getenv("redisHost"), Password: os.Getenv("redisPassword")})

	return config
}
