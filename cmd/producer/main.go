package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	producer "github.com/thatjimmi/go-kafka/internal/producer"
)

func main() {
	go producer.RunProducer()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
