package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Kafka struct {
	*kafka.Conn
	Topic     string
	Network   string
	Address   string
	Partition int
}

// New returns a new Kafka instance.
func New() *Kafka {
	return &Kafka{}
}

// Connect sets `k.Conn` to a new connection, and sets up graceful exit to
// close the connection when exiting.
func (k *Kafka) Connect() error {
	// dial connection
	conn, err := kafka.DialLeader(context.Background(), k.Network, k.Address, k.Topic, k.Partition)
	if err != nil {
		return fmt.Errorf("failed to dial leader: %w", err)
	}
	k.Conn = conn

	// handle exit signals to gracefully shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exit
		log.Println("closing kafka writer")
		if err := k.Close(); err != nil {
			log.Println(fmt.Errorf("failed to close kafka writer: %w", err))
		} else {
			log.Println("kafka writer closed")
		}
	}()
	return nil
}

// WithTopic sets k.Topic.
func (k *Kafka) WithTopic(topic string) *Kafka {
	k.Topic = topic
	return k
}

// WithAddress sets k.Address.
func (k *Kafka) WithAddress(address string) *Kafka {
	k.Address = address
	return k
}

// WithNetwork sets k.Network.
func (k *Kafka) WithNetwork(network string) *Kafka {
	k.Network = network
	return k
}

// WithPartition sets k.Partition.
func (k *Kafka) WithPartition(partition int) *Kafka {
	k.Partition = partition
	return k
}

// WriteMessage writes a message with the current timestamp through its Kafka
// connection.
func (k *Kafka) WriteMessage(value []byte) error {
	_, err := k.WriteMessages(kafka.Message{Value: value, Time: time.Now()})
	return err
}
