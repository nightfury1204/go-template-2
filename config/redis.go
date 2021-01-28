package config

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Redis holds mongo config
type Redis struct {
	URL          string        `yaml:"url"`
	RedisTimeOut time.Duration `yaml:"time_out"`
}

var redisOnce = sync.Once{}
var redisConfig *Redis

// loadRedis loads config from path
func loadRedis(fileName string) error {
	readConfig(fileName)
	viper.AutomaticEnv()

	redisConfig = &Redis{
		URL:          viper.GetString("redis.url"),
		RedisTimeOut: viper.GetDuration("redis.time_out") * time.Second,
	}

	return nil
}

// GetRedis returns redis config
func GetRedis(fileName string) *Redis {
	redisOnce.Do(func() {
		loadRedis(fileName)
	})

	return redisConfig
}
