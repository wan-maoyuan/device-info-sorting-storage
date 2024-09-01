package kafka

import (
	"context"
	"device-info-sorting-storage/pkg/conf"
	"testing"
	"time"
)

func TestNewKafkaReader(t *testing.T) {
	conf.Get().KafkaBrokers = []string{"192.168.2.43:9094"}
	conf.Get().KafkaGroup = "example"
	conf.Get().KafkaTopic = "example"

	reader, err := NewKafkaReader()
	if err != nil {
		t.Errorf("NewKafkaReader: %v", err)
		return
	}
	defer reader.Close()
}

func TestKafkaReaderReadMessage(t *testing.T) {
	conf.Get().KafkaBrokers = []string{"192.168.2.43:9094"}
	conf.Get().KafkaGroup = "example"
	conf.Get().KafkaTopic = "example"

	reader, err := NewKafkaReader()
	if err != nil {
		t.Errorf("NewKafkaReader: %v", err)
		return
	}
	defer reader.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()

	ch := reader.ReadMessage(ctx)

	for msg := range ch {
		t.Logf("%v", msg)
	}
}
