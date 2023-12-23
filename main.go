package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Missing argument. Please specify 'producer' or 'consumer'")
	}

	switch os.Args[1] {
	case "producer":
		RunProducer()
	case "consumer":
		RunConsumer()
	default:
		panic("Invalid argument")
	}
}
