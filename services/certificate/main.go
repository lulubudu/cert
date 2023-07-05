package main

import (
	"certificate/db/postgres"
	"certificate/notifier"
	"certificate/notifier/kafka"
	"certificate/router"
	"fmt"
	"log"
)

func main() {
	// create database instance
	db, err := postgres.Connect()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to db: %w", err))
	}

	// create kafka instance
	k := kafka.New().
		WithNetwork("tcp").
		WithAddress("kafka:29092").
		WithTopic("cert-active-status-toggled").
		WithPartition(0)
	if err := k.Connect(); err != nil {
		log.Fatal(fmt.Errorf("failed to connect to kafka: %w", err))
	}

	// create and start HTTP server
	if err := router.New().
		WithDatabase(db).
		WithNotifier(notifier.New(k)).Start("0.0.0.0:8080"); err != nil {
		log.Fatal(fmt.Errorf("failed to start http server: %w", err))
	}
}
