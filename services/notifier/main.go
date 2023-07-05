package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"notifier/kafka"
	"os"
)

// handleCertToggle takes in message content and sends it to the specified
// HTTP endpoint.
func handleCertToggle(msg []byte, endpoint string) {
	log.Println("received message", string(msg))
	resp, err := http.Post(endpoint+"/cert-active-status-toggled", "application/json", bytes.NewReader(msg))
	if err != nil {
		// TODO: add retry logic by not acknowledging receipt of the kafka message
		log.Println(fmt.Errorf("failed to send notify via HTTP POST: %w", err))
	}
	log.Println("got status code", resp.StatusCode)
}

func main() {
	// get configured endpoint
	endpoint := os.Getenv("ENDPOINT")
	if endpoint == "" {
		log.Fatal("ENDPOINT ENV not set")
	}

	// get configured kafka address
	addr := os.Getenv("KAFKA_ADDR")
	if addr == "" {
		log.Fatal("KAFKA_ADDR ENV not set")
	}

	// establish kafka reader
	messageChan, err := kafka.ReadCertActiveStatusToggled(addr)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get kafka consumer: %w", err))
	}
	log.Println("listening for messages")

	// process messages
	for m := range messageChan {
		// spawn a go routine here for concurrent processing of messages
		// I'd do some time based backoff in real code so that requests
		//are not sent too close to each other
		go handleCertToggle(m.Value, endpoint)
	}
}
