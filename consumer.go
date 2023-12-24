package main

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jackc/pgx/v4/pgxpool"
)

func RunConsumer() {
	db := connectDB()
	defer db.Close()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"myTopic"}, nil)

	defer c.Close()

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			go func(value []byte) {
				_, err = db.Exec(context.Background(), "INSERT INTO messages (message) VALUES ($1)", string(value))
				if err != nil {
					log.Printf("Error inserting message into database: %v\n", err)
				}
			}(msg.Value)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func connectDB() *pgxpool.Pool {
	connStr := "postgres://myuser:mypassword@localhost:5433/mydb?sslmode=disable"
	db, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
