package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/millken/logger"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	//"sync"
	"time"
)

const TS = "2006-01-02T15:04:05.000Z"

// Output plugin that index messages to an elasticsearch cluster.
// Largely based on FileOutput plugin.
type ElasticSearchOutput struct {
	flushInterval uint32
	flushCount    int
	batchChan     chan []byte
	backChan      chan []byte
	// The BulkIndexer used to index documents
	bulkIndexer BulkIndexer

	// Specify an overall timeout value in milliseconds for bulk request to complete.
	// Default is 0 (infinite)
	http_timeout uint32
	//Disable both TCP and HTTP keepalives
	http_disable_keepalives bool
	// Specify a resolve and connect timeout value in milliseconds for bulk request.
	// It's always included in overall request timeout (see 'http_timeout' option).
	// Default is 0 (infinite)
	connect_timeout uint32
}

func (o *ElasticSearchOutput) Init() (err error) {
	conf := config.ElasticSearchOutput
	o.flushInterval = conf.FlushInterval
	o.flushCount = conf.FlushCount
	o.batchChan = make(chan []byte)
	o.backChan = make(chan []byte, 2)
	o.http_timeout = conf.HTTPTimeout
	o.http_disable_keepalives = conf.HTTPDisableKeepalives
	o.connect_timeout = conf.ConnectTimeout
	var serverUrl *url.URL
	if serverUrl, err = url.Parse(conf.Server); err == nil {
		var scheme string = strings.ToLower(serverUrl.Scheme)
		switch scheme {
		case "http":

			o.bulkIndexer = NewHttpBulkIndexer(scheme,
				serverUrl.Host, o.flushCount, conf.Username, conf.Password,
				o.http_timeout, o.http_disable_keepalives, o.connect_timeout, nil)

		default:
			err = errors.New("Server URL must specify one of `udp`, `http`, or `https`.")
		}
	} else {
		err = fmt.Errorf("Unable to parse ElasticSearch server URL [%s]: %s", conf.Server, err)
	}
	return
}

// Runs in a separate goroutine, accepting incoming messages, buffering output
// data until the ticker triggers the buffered data should be put onto the
// committer channel.
func (o *ElasticSearchOutput) receiver(inChan chan map[string]interface{}) {
	var (
		count int
	)
	ok := true
	ticker := time.Tick(time.Duration(o.flushInterval) * time.Millisecond)
	outBatch := make([]byte, 0, 10000)

	for ok {
		select {
		case pack := <-inChan:
			packBytes := make([]byte, 0, 10000)

			outBytes, _ := json.Marshal(pack)

			t := time.Now()
			h1 := fmt.Sprintf("{\"index\":{\"_index\":\"nginx-%s\",\"_type\":\"nginx.access\"}}", t.Format("2006-01-02"))

			H := []byte(h1)
			packBytes = append(packBytes, H...)
			packBytes = append(packBytes, []byte("\n")...)
			packBytes = append(packBytes, outBytes...)
			packBytes = append(packBytes, []byte("\n")...)
			/*
				for k, v := range pack {
					switch vv := v.(type) {
					case string:
						fmt.Println(k, "is string", vv)
					case int, float64:
						fmt.Println(k, "is int", vv)
					case []interface{}:
						fmt.Println(k, "is an array:")
						for i, u := range vv {
							fmt.Println(i, u)
						}
					default:
						fmt.Println(k, "is of a type I don't know how to handle")
					}
				}*/

			outBatch = append(outBatch, packBytes...)
			if count = count + 1; o.bulkIndexer.CheckFlush(count, len(outBatch)) {
				if len(outBatch) > 0 {
					// This will block until the other side is ready to accept
					// this batch, so we can't get too far ahead.
					o.batchChan <- outBatch
					outBatch = <-o.backChan
					count = 0
				}
			}

		case <-ticker:
			if len(outBatch) > 0 {
				// This will block until the other side is ready to accept
				// this batch, freeing us to start on the next one.
				o.batchChan <- outBatch
				outBatch = <-o.backChan
				count = 0
			}
		}
	}
}

// Runs in a separate goroutine, waits for buffered data on the committer
// channel, bulk index it out to the elasticsearch cluster, and puts the now
// empty buffer on the return channel for reuse.
func (o *ElasticSearchOutput) committer() {
	initBatch := make([]byte, 0, 10000)
	o.backChan <- initBatch
	var outBatch []byte

	for outBatch = range o.batchChan {
		if err := o.bulkIndexer.Index(outBatch); err != nil {
			logger.Error(err)
		}
		outBatch = outBatch[:0]
		o.backChan <- outBatch
	}
}

func (o *ElasticSearchOutput) Run() (err error) {

	go o.receiver(esCh)
	go o.committer()
	return
}

// A BulkIndexer is used to index documents in ElasticSearch
type BulkIndexer interface {
	// Index documents
	Index(body []byte) error
	// Check if a flush is needed
	CheckFlush(count int, length int) bool
}

// A HttpBulkIndexer uses the HTTP REST Bulk Api of ElasticSearch
// in order to index documents
type HttpBulkIndexer struct {
	// Protocol (http or https).
	Protocol string
	// Host name and port number (default to "localhost:9200").
	Domain string
	// Maximum number of documents.
	MaxCount int
	// Internal HTTP Client.
	client *http.Client
	// Optional username for HTTP authentication
	username string
	// Optional password for HTTP authentication
	password string
}

func NewHttpBulkIndexer(protocol string, domain string, maxCount int,
	username string, password string, httpTimeout uint32, httpDisableKeepalives bool,
	connectTimeout uint32, tlsConf *tls.Config) *HttpBulkIndexer {

	tr := &http.Transport{
		TLSClientConfig:   tlsConf,
		DisableKeepAlives: httpDisableKeepalives,
		Dial: func(network, address string) (net.Conn, error) {
			return net.DialTimeout(network, address, time.Duration(connectTimeout)*time.Millisecond)
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(httpTimeout) * time.Millisecond,
	}
	return &HttpBulkIndexer{
		Protocol: protocol,
		Domain:   domain,
		MaxCount: maxCount,
		client:   client,
		username: username,
		password: password,
	}
}

func (h *HttpBulkIndexer) CheckFlush(count int, length int) bool {
	if count >= h.MaxCount {
		return true
	}
	return false
}

func (h *HttpBulkIndexer) Index(body []byte) error {
	var response_body []byte
	var response_body_json map[string]interface{}

	url := fmt.Sprintf("%s://%s%s", h.Protocol, h.Domain, "/_bulk")

	// Creating ElasticSearch Bulk HTTP request
	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("Can't create bulk request: %s", err.Error())
	}
	request.Header.Add("Accept", "application/json")
	if h.username != "" && h.password != "" {
		request.SetBasicAuth(h.username, h.password)
	}

	request_start_time := time.Now()
	response, err := h.client.Do(request)
	request_time := time.Since(request_start_time)
	if err != nil {
		if (h.client.Timeout > 0) && (request_time >= h.client.Timeout) &&
			(strings.Contains(err.Error(), "use of closed network connection")) {

			return fmt.Errorf("HTTP request was interrupted after timeout. It lasted %s",
				request_time.String())
		} else {
			return fmt.Errorf("HTTP request failed: %s", err.Error())
		}
	}
	if response != nil {
		defer response.Body.Close()
		if response.StatusCode > 304 {
			return fmt.Errorf("HTTP response error status: %s", response.Status)
		}
		if response_body, err = ioutil.ReadAll(response.Body); err != nil {
			return fmt.Errorf("Can't read HTTP response body: %s", err.Error())
		}
		err = json.Unmarshal(response_body, &response_body_json)
		if err != nil {
			return fmt.Errorf("HTTP response didn't contain valid JSON. Body: %s",
				string(response_body))
		}
		json_errors, ok := response_body_json["errors"].(bool)
		if ok && json_errors {
			return fmt.Errorf("ElasticSearch server reported error within JSON: %s",
				string(response_body))
		}
	}
	return nil
}

func startElasticSearchService() {
	logger.Info("startElasticSearchService()")
	e := new(ElasticSearchOutput)
	if err := e.Init(); err != nil {
		logger.Error("Init ElasticSearch service Fail.Err = %s", err.Error())
	}
	go e.Run()
}
