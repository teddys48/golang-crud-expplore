package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisConfig(viper *viper.Viper) *redis.Client {

	config := redis.NewClient(&redis.Options{Addr: viper.GetString("redis.address"), Password: viper.GetString("redis.password")})

	return config
}
