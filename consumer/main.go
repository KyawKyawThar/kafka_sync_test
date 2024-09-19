package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	broker1Address = "localhost:9093"
	//broker2Address = "localhost:9094"
	//broker3Address = "localhost:9095"
	topic     = "message-logs"
	partition = 0
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Connect to the Kafka broker (replace with your Kafka broker's address)
	brokerList := []string{broker1Address}

	// Create a new consumer
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Kafka consumer:", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln("Failed to close Kafka consumer:", err)
		}
	}()

	partition, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)

	if err != nil {
		log.Fatalln("Failed to consume partition:", err)
	}

	// Handle OS signals to safely exit the consumer
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Starting to consume messages...")
	var msgCount = 0
	for {
		select {
		case msg := <-partition.Messages():
			fmt.Printf("Consumed message: %s, offset: %d\n", string(msg.Value), msg.Offset)
			msgCount++
			fmt.Printf("Consumed %d messages. Exiting...\n", msgCount)
			return
		case err := <-partition.Errors():
			fmt.Println("Failed to consume message:", err)
		case <-signals:
			fmt.Println("Interrupt signal received, shutting down...")
			return
		}
	}
}
