package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

type tomlConfig struct {
	Kafka KafkaInputConfig `toml:"kafka"`
}

type KafkaInputConfig struct {
	Id                         string
	Addrs                      []string
	MetadataRetries            int    `toml:"metadata_retries"`
	WaitForElection            uint32 `toml:"wait_for_election"`
	BackgroundRefreshFrequency uint32 `toml:"background_refresh_frequency"`

	// Broker Config
	MaxOpenRequests int    `toml:"max_open_reqests"`
	DialTimeout     uint32 `toml:"dial_timeout"`
	ReadTimeout     uint32 `toml:"read_timeout"`
	WriteTimeout    uint32 `toml:"write_timeout"`

	// Consumer Config
	Topic            string
	Partition        int32
	Group            string
	DefaultFetchSize int32  `toml:"default_fetch_size"`
	MinFetchSize     int32  `toml:"min_fetch_size"`
	MaxMessageSize   int32  `toml:"max_message_size"`
	MaxWaitTime      uint32 `toml:"max_wait_time"`
	OffsetMethod     string `toml:"offset_method"` // Manual, Newest, Oldest
	EventBufferSize  int    `toml:"event_buffer_size"`
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
	hn, err := os.Hostname()
	if err != nil {
		hn = "logelasticsearch"
	}
	config.Kafka = KafkaInputConfig{
		Id:                         hn,
		MetadataRetries:            3,
		WaitForElection:            250,
		BackgroundRefreshFrequency: 10 * 60 * 1000,
		MaxOpenRequests:            4,
		DialTimeout:                60 * 1000,
		ReadTimeout:                60 * 1000,
		WriteTimeout:               60 * 1000,
		DefaultFetchSize:           1024 * 32,
		MinFetchSize:               1,
		MaxWaitTime:                250,
		OffsetMethod:               "Manual",
		EventBufferSize:            16,
	}

	if _, err = toml.Decode(string(contents), &config); err != nil {
		return fmt.Errorf("Error decoding config file: %s", err)
	}

	return nil
}
