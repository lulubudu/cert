package kafka_test

import (
	"certificate/notifier/kafka"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, kafka.New())
}

func TestKafka_WithNetwork(t *testing.T) {
	mockNetwork := "mock_network"
	k := kafka.New().WithNetwork(mockNetwork)
	assert.Equal(t, mockNetwork, k.Network)
}

func TestKafka_WithAddress(t *testing.T) {
	mockAddr := "mock_addr"
	k := kafka.New().WithAddress(mockAddr)
	assert.Equal(t, mockAddr, k.Address)
}

func TestKafka_WithTopic(t *testing.T) {
	mockTopic := "mock_topic"
	k := kafka.New().WithTopic(mockTopic)
	assert.Equal(t, mockTopic, k.Topic)
}

func TestKafka_WithPartition(t *testing.T) {
	mockPartition := 0
	k := kafka.New().WithPartition(mockPartition)
	assert.Equal(t, mockPartition, k.Partition)
}

func TestKafka_Connect(t *testing.T) {

}
