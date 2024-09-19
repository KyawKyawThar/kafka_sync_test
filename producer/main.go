package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"strconv"
)

const (
	broker1Address = "localhost:9093"
	//broker2Address = "localhost:9094"
	//broker3Address = "localhost:9095"
	topic = "message-logs"
)

func main() {
	// Setup configuration for the producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll

	// Connect to the Kafka broker (replace with your Kafka broker's address)
	brokerList := []string{broker1Address}

	// Create a new async Kafka producer
	producer, err := sarama.NewSyncProducer(brokerList, config)

	if err != nil {
		log.Fatalln("Failed to start Kafka producer:", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln("Failed to close Kafka producer:", err)
		}
	}()

	numMessages := 10

	for i := 0; i < numMessages; i++ {

		message := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder("Message number from kafka: " + strconv.Itoa(i)),
		}

		partition, offset, err := producer.SendMessage(message)

		if err != nil {
			log.Println("Failed to send message:", err)
		} else {
			fmt.Printf("Message %d stored in topic(%s)/partition(%d)/offset(%d)\n", i, topic, partition, offset)
		}
	}
	fmt.Println("All messages sent successfully!")

}
