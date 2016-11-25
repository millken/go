package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/millken/logger"
	"os"
	"path/filepath"
	//"sync/atomic"
	"time"
)

type KafkaInput struct {
	processMessageCount    int64
	processMessageFailures int64

	config             KafkaInputConfig
	clientConfig       *sarama.Config
	consumer           sarama.Consumer
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

	k.clientConfig = sarama.NewConfig()
	k.clientConfig.Metadata.Retry.Max = k.config.MetadataRetries
	k.clientConfig.Metadata.Retry.Backoff = time.Duration(k.config.WaitForElection) * time.Millisecond
	k.clientConfig.Metadata.RefreshFrequency = time.Duration(k.config.BackgroundRefreshFrequency) * time.Millisecond

	k.clientConfig.Net.MaxOpenRequests = k.config.MaxOpenRequests
	k.clientConfig.Net.DialTimeout = time.Duration(k.config.DialTimeout) * time.Millisecond
	k.clientConfig.Net.ReadTimeout = time.Duration(k.config.ReadTimeout) * time.Millisecond
	k.clientConfig.Net.WriteTimeout = time.Duration(k.config.WriteTimeout) * time.Millisecond

	k.clientConfig.Consumer.Fetch.Default = k.config.DefaultFetchSize
	k.clientConfig.Consumer.Fetch.Min = k.config.MinFetchSize
	k.clientConfig.Consumer.Fetch.Max = k.config.MaxMessageSize
	k.clientConfig.Consumer.MaxWaitTime = time.Duration(k.config.MaxWaitTime) * time.Millisecond
	k.checkpointFilename = filepath.Join("kafka",
		fmt.Sprintf("%s.%d.offset.bin", k.config.Topic, k.config.Partition))

	switch k.config.OffsetMethod {
	case "Manual":
		if fileExists(k.checkpointFilename) {
			if k.config.OffsetValue, err = readCheckpoint(k.checkpointFilename); err != nil {
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
		k.config.OffsetValue = sarama.OffsetNewest
		if fileExists(k.checkpointFilename) {
			if err = os.Remove(k.checkpointFilename); err != nil {
				return
			}
		}
	case "Oldest":
		k.config.OffsetValue = sarama.OffsetOldest
		if fileExists(k.checkpointFilename) {
			if err = os.Remove(k.checkpointFilename); err != nil {
				return
			}
		}
	default:
		return fmt.Errorf("invalid offset_method: %s", k.config.OffsetMethod)
	}

	k.consumer, err = sarama.NewConsumer(k.config.Addrs, k.clientConfig)
	return

}

func (k *KafkaInput) Run() (err error) {
	consumer, err := k.consumer.ConsumePartition(k.config.Topic, k.config.Partition, k.config.OffsetValue)
	if err != nil {
		logger.Error("Run Kafka service Fail.Err = %s", err.Error())
		//return
	}

	defer func() {
		k.consumer.Close()
		consumer.Close()
		if k.checkpointFile != nil {
			k.checkpointFile.Close()
		}
	}()
	k.stopChan = make(chan bool)

	for {
		select {
		case message := <-consumer.Messages():
			workerCh <- string(message.Value)
			logger.Info(string(message.Value))

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
