package config

import (
	"flag"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

var (
	configFilePath = flag.String("config", "./config/config.yaml", "config file")
	nodeIDFlag     = flag.Uint64("nodeID", 1, "node id")
	cfg            Config
)

func GetConfig() Config {
	return cfg
}

type Config struct {
	Retry Retry `json:"retry"`
	//TimeWheel TimeWheel `json:"time_wheel"`
	Broker *Broker `json:"broker"`
	Store  *Store  `json:"store"`
	Log    *Log    `json:"log"`
	Server Server  `json:"server"`
}

func (c Config) GetLog() Log {
	if c.Log == nil {
		return Log{}
	}
	return *c.Log
}

func LoadConfig() error {
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles(*configFilePath)
	if err != nil {
		return err
	}
	return config.BindStruct("", &cfg)
}

func GetPubMaxQos() uint8 {
	return 1
}
