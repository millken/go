package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

type tomlConfig struct {
	Kafka kafkaConf `toml:"kafka"`
}

type kafkaConf struct {
	Addrs []string
	Topic string
}

func LoadConfig(configPath string) (err error) {

	p, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("Error opening config file: %s", err)
	}
	contents, err := ioutil.ReadAll(p)
	if err != nil {
		return fmt.Errorf("Error reading config file: %s", err)
	}

	if _, err = toml.Decode(string(contents), &config); err != nil {
		return fmt.Errorf("Error decoding config file: %s", err)
	}

	return nil
}
