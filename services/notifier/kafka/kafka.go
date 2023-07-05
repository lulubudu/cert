package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// ReadCertActiveStatusToggled returns a channel for receiving messages
func ReadCertActiveStatusToggled(addr string) (chan kafka.Message, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{addr},
		Topic:    "cert-active-status-toggled",
		MaxBytes: 10e6,
	})
	messageChan := make(chan kafka.Message)
	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Println(fmt.Errorf("failed to read message: %w", err))
			}
			messageChan <- m
		}
	}()
	// handle exit signals to gracefully shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exit
		log.Println("closing kafka reader")
		if err := r.Close(); err != nil {
			log.Println(fmt.Errorf("failed to close kafka reader: %w", err))
		} else {
			log.Println("kafka reader closed")
		}
	}()
	return messageChan, nil
}
