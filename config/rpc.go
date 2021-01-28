package config

import (
	"sync"

	"github.com/spf13/viper"
)

// RPC holds RPC config
type RPC struct {
	URL string `yaml:"rpc_url"`
}

var rpcOnce = sync.Once{}
var rpcConfig *RPC

// loadRPC loads config from path
func loadRPC(fileName string) error {
	readConfig(fileName)
	viper.AutomaticEnv()

	rpcConfig = &RPC{
		URL: viper.GetString("rpc_url"),
	}

	return nil
}

// GetRPC returns RPC config
func GetRPC(fileName string) *RPC {
	rpcOnce.Do(func() {
		loadRPC(fileName)
	})

	return rpcConfig
}
