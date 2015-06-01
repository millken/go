package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/millken/logger"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type KafkaInput struct {
	processMessageCount    int64
	processMessageFailures int64

	config             KafkaInputConfig
	clientConfig       *sarama.ClientConfig
	consumerConfig     *sarama.ConsumerConfig
	client             *sarama.Client
	consumer           *sarama.Consumer
	checkpointFile     *os.File
	stopChan           chan bool
	checkpointFilename string
}

func (k *KafkaInput) writeCheckpoint(offset int64) (err error) {
	if k.checkpointFile == nil {
		if k.checkpointFile, err = os.OpenFile(k.checkpointFilename,
			os.O_WRONLY|os.O_SYNC|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return
		}
	}
	k.checkpointFile.Seek(0, 0)
	err = binary.Write(k.checkpointFile, binary.LittleEndian, &offset)
	return
}

func (k *KafkaInput) Init() (err error) {
	k.config = config.KafkaInput
	if len(k.config.Addrs) == 0 {
		return errors.New("addrs must have at least one entry")
	}
	if len(k.config.Group) == 0 {
		k.config.Group = k.config.Id
	}

	k.clientConfig = sarama.NewClientConfig()
	k.clientConfig.MetadataRetries = k.config.MetadataRetries
	k.clientConfig.WaitForElection = time.Duration(k.config.WaitForElection) * time.Millisecond
	k.clientConfig.BackgroundRefreshFrequency = time.Duration(k.config.BackgroundRefreshFrequency) * time.Millisecond

	k.clientConfig.DefaultBrokerConf = sarama.NewBrokerConfig()
	k.clientConfig.DefaultBrokerConf.MaxOpenRequests = k.config.MaxOpenRequests
	k.clientConfig.DefaultBrokerConf.DialTimeout = time.Duration(k.config.DialTimeout) * time.Millisecond
	k.clientConfig.DefaultBrokerConf.ReadTimeout = time.Duration(k.config.ReadTimeout) * time.Millisecond
	k.clientConfig.DefaultBrokerConf.WriteTimeout = time.Duration(k.config.WriteTimeout) * time.Millisecond

	k.consumerConfig = sarama.NewConsumerConfig()
	k.consumerConfig.DefaultFetchSize = k.config.DefaultFetchSize
	k.consumerConfig.MinFetchSize = k.config.MinFetchSize
	k.consumerConfig.MaxMessageSize = k.config.MaxMessageSize
	k.consumerConfig.MaxWaitTime = time.Duration(k.config.MaxWaitTime) * time.Millisecond
	k.checkpointFilename = filepath.Join("kafka",
		fmt.Sprintf("%s.%d.offset.bin", k.config.Topic, k.config.Partition))

	switch k.config.OffsetMethod {
	case "Manual":
		k.consumerConfig.OffsetMethod = sarama.OffsetMethodManual
		if fileExists(k.checkpointFilename) {
			if k.consumerConfig.OffsetValue, err = readCheckpoint(k.checkpointFilename); err != nil {
				return fmt.Errorf("readCheckpoint %s", err)
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(k.checkpointFilename), 0766); err != nil {
				return
			}
			if err = k.writeCheckpoint(0); err != nil {
				return
			}
		}
	case "Newest":
		k.consumerConfig.OffsetMethod = sarama.OffsetMethodNewest
		if fileExists(k.checkpointFilename) {
			if err = os.Remove(k.checkpointFilename); err != nil {
				return
			}
		}
	case "Oldest":
		k.consumerConfig.OffsetMethod = sarama.OffsetMethodOldest
		if fileExists(k.checkpointFilename) {
			if err = os.Remove(k.checkpointFilename); err != nil {
				return
			}
		}
	default:
		return fmt.Errorf("invalid offset_method: %s", k.config.OffsetMethod)
	}

	k.consumerConfig.EventBufferSize = k.config.EventBufferSize

	k.client, err = sarama.NewClient(k.config.Id, k.config.Addrs, k.clientConfig)
	if err != nil {
		return
	}
	k.consumer, err = sarama.NewConsumer(k.client, k.config.Topic, k.config.Partition, k.config.Group, k.consumerConfig)
	return
}

func (k *KafkaInput) Run() (err error) {
	defer func() {
		k.consumer.Close()
		k.client.Close()
		if k.checkpointFile != nil {
			k.checkpointFile.Close()
		}
	}()
	k.stopChan = make(chan bool)

	for {
		select {
		case event, ok := <-k.consumer.Events():
			if !ok {
				return
			}
			atomic.AddInt64(&k.processMessageCount, 1)
			if event.Err != nil {
				if event.Err == sarama.OffsetOutOfRange {
					logger.Error(fmt.Errorf("removing the out of range checkpoint file and stopping"))
					if err := os.Remove(k.checkpointFilename); err != nil {
						logger.Error(err)
					}
					return
				}
				atomic.AddInt64(&k.processMessageFailures, 1)
				logger.Error(event.Err)
				break
			}

			if k.consumerConfig.OffsetMethod == sarama.OffsetMethodManual {
				if err = k.writeCheckpoint(event.Offset + 1); err != nil {
					return
				}
			}
			workerCh <- string(event.Value)
			//logger.Debug(string(event.Value))

		case <-k.stopChan:
			return
		}
	}
	return
}

func (k *KafkaInput) Stop() {
	close(k.stopChan)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func readCheckpoint(filename string) (offset int64, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	err = binary.Read(file, binary.LittleEndian, &offset)
	return
}

func startKafkaService() {
	logger.Info("startKafkaService()")
	k := new(KafkaInput)
	if err := k.Init(); err != nil {
		logger.Error("Init Kafka service Fail.Err = %s", err.Error())
	}
	k.Run()
	/*
	   	logger.Info("Ready kafka service: broker = %v", config.Kafka.Addrs)
	   	client, err := sarama.NewClient("logelasticsearch", config.Kafka.Addrs, nil)
	   	if err != nil {
	   		logger.Exitf("connect kafka failed.Err = %s", err.Error())
	   	} else {
	   		logger.Info("kafka connected")
	   	}
	   	defer client.Close()

	   	consumer, err := sarama.NewConsumer(client, config.Kafka.Topic, 0, "id", sarama.NewConsumerConfig())
	   	if err != nil {
	   		panic(err)
	   	} else {
	   		logger.Info("Kafka consumer ready")
	   	}
	   	defer consumer.Close()
	   	msgCount := 0

	   consumerLoop:
	   	for {
	   		select {
	   		case event := <-consumer.Events():
	   			if event.Err != nil {
	   				logger.Error(event.Err)
	   			}
	   			logger.Debug(string(event.Value))
	   			msgCount++
	   			//case <-time.After(5 * time.Second):
	   			//logger.Warn("Kafka Timeout")
	   			continue consumerLoop
	   		}
	   	}
	*/
}
