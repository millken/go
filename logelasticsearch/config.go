package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

type tomlConfig struct {
	Ipdb                string
	KafkaInput          KafkaInputConfig
	ElasticSearchOutput ElasticSearchOutputConfig
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

// ConfigStruct for ElasticSearchOutput plugin.
type ElasticSearchOutputConfig struct {
	// Interval at which accumulated messages should be bulk indexed to
	// ElasticSearch, in milliseconds (default 1000, i.e. 1 second).
	FlushInterval uint32 `toml:"flush_interval"`
	// Number of messages that triggers a bulk indexation to ElasticSearch
	// (default to 10)
	FlushCount int `toml:"flush_count"`
	// ElasticSearch server address. This address also defines the Bulk
	// indexing mode. For example, "http://localhost:9200" defines a server
	// accessible on localhost and the indexing will be done with the HTTP
	// Bulk API, whereas "udp://192.168.1.14:9700" defines a server accessible
	// on the local network and the indexing will be done with the UDP Bulk
	// API. (default to "http://localhost:9200")
	Server string
	// Optional subsection for TLS configuration of ElasticSearch connections. If
	// unspecified, the default ElasticSearch settings will be used.
	//Tls tcp.TlsConfig
	// Optional ElasticSearch username for HTTP authentication. This is useful
	// if you have put your ElasticSearch cluster behind a proxy like nginx.
	// and turned on authentication.
	Username string `toml:"username"`
	// Optional password for HTTP authentication.
	Password string `toml:"password"`
	// Overall timeout
	HTTPTimeout uint32 `toml:"http_timeout"`
	// Disable both TCP and HTTP keepalives
	HTTPDisableKeepalives bool `toml:"http_disable_keepalives"`
	// Resolve and connect timeout only
	ConnectTimeout uint32 `toml:"connect_timeout"`
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
	config.KafkaInput = KafkaInputConfig{
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
	config.ElasticSearchOutput = ElasticSearchOutputConfig{
		FlushInterval:         1000,
		FlushCount:            10,
		Server:                "http://localhost:9200",
		Username:              "",
		Password:              "",
		HTTPTimeout:           0,
		HTTPDisableKeepalives: false,
		ConnectTimeout:        0,
	}

	if _, err = toml.Decode(string(contents), &config); err != nil {
		return fmt.Errorf("Error decoding config file: %s", err)
	}

	return nil
}
