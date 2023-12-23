package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
)

const (
	websocketURL = "ws://localhost:8080/ws"
	kafkaServer  = "localhost:9092"
	kafkaTopic   = "myTopic"
)

func RunProducer() {
	ws, err := connectWebSocket(websocketURL)
	if err != nil {
		panic("Failed to connect to WebSocket server: " + err.Error())
	}
	fmt.Println("Connected to WebSocket server")
	defer ws.Close()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaServer})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	go handleKafkaEvents(p)

	for {
		if err := processWebSocketMessages(ws, p); err != nil {
			fmt.Printf("Error processing WebSocket messages: %v\n", err)
			break
		}
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

func connectWebSocket(url string) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	return c, nil
}

func processWebSocketMessages(ws *websocket.Conn, p *kafka.Producer) error {
	topic := kafkaTopic
	_, message, err := ws.ReadMessage()
	if err != nil {
		return fmt.Errorf("read message: %w", err)
	}
	log.Printf("Received message: %s\n", message)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	if err != nil {
		return fmt.Errorf("produce message: %w", err)
	}
	return nil
}

func handleKafkaEvents(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}
