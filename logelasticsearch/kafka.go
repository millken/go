package main

import (
	"github.com/Shopify/sarama"
	"github.com/millken/logger"
	"os"
	//"time"
)

func startKafkaService() {
	logger.Info("Ready kafka service: broker = %v", config.Kafka.Addrs)
	client, err := sarama.NewClient("logelasticsearch", config.Kafka.Addrs, nil)
	if err != nil {
		logger.Exitf("connect kafka failed.Err = %s", err.Error())
	} else {
		logger.Info("kafka connected")
	}
	defer client.Close()
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "logelasticsearch"
	}
	consumer, err := sarama.NewConsumer(client, config.Kafka.Topic, 0, hostname, sarama.NewConsumerConfig())
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
}
