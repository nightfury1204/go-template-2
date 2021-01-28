package config

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Mongo holds mongo config
type Mongo struct {
	URL       string        `yaml:"url"`
	DBName    string        `yaml:"db_name"`
	DBTimeOut time.Duration `yaml:"time_out"`
}

var mongoOnce = sync.Once{}
var mongoConfig *Mongo

// loadMongo loads config from path
func loadMongo(fileName string) error {
	readConfig(fileName)
	viper.AutomaticEnv()

	mongoConfig = &Mongo{
		URL:       viper.GetString("mongo.url"),
		DBName:    viper.GetString("mongo.db_name"),
		DBTimeOut: viper.GetDuration("mongo.time_out") * time.Second,
	}

	return nil
}

// GetMongo returns mongo config
func GetMongo(fileName string) *Mongo {
	mongoOnce.Do(func() {
		loadMongo(fileName)
	})

	return mongoConfig
}
