package producer

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	protomessage "github.com/thatjimmi/go-kafka/proto"
	"google.golang.org/protobuf/proto"
)

var (
	websocketURL = "ws://localhost:8080/ws"
	kafkaServer  = "localhost:9092"
	kafkaTopic   = "myTopic"
)

var (
	// Example of a counter metric
	ProducerCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "producer_counter",
		Help: "The total number of processed messages",
	})
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
	_, message, err := ws.ReadMessage()
	if err != nil {
		return fmt.Errorf("read message: %w", err)
	}
	log.Printf("Received message: %s\n", message)

	// Create new Protobuf message
	protoMsg := &protomessage.Message{
		Content: string(message),
	}

	// serialize message
	serializedMsg, err := proto.Marshal(protoMsg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          serializedMsg,
	}, nil)
	if err != nil {
		return fmt.Errorf("produce message: %w", err)
	}

	fmt.Println("Produced message to Kafka")
	ProducerCounter.Inc()
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
