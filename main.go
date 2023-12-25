package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Missing argument. Must be one of: producer, consumer, websocket, metrics")
	}

	switch os.Args[1] {
	case "producer":
		RunProducer()
	case "consumer":
		RunConsumer()
	case "websocket":
		RunWebSocket()
	case "metrics":
		// Register Prometheus metrics
		fmt.Println("Registering Prometheus metrics")
	default:
		panic("Invalid argument")
	}
}
