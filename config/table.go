package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// Table holds table configurations
type Table struct {
	BrandCollectionName string `yaml:"brand"`
}

var tableOnce = sync.Once{}
var tableConfig *Table

// loadTable loads config from path
func loadTable(fileName string) error {
	readConfig(fileName)
	viper.AutomaticEnv()

	tableConfig = &Table{
		BrandCollectionName: viper.GetString("collection.brand"),
	}

	log.Println("table config ", appConfig)

	return nil
}

// GetTable returns table config
func GetTable(fileName string) *Table {
	tableOnce.Do(func() {
		loadTable(fileName)
	})

	return tableConfig
}
