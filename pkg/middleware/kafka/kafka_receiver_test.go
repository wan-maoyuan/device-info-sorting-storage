package kafka

import (
	"context"
	"device-info-sorting-storage/pkg/conf"
	"testing"
)

func TestNewKafkaReceiver(t *testing.T) {
	conf.Get().KafkaBrokers = []string{"47.92.72.122:31594"}
	conf.Get().KafkaGroup = "example"
	conf.Get().KafkaTopic = "example"

	receiver, err := NewKafkaReceiver()
	if err != nil {
		t.Errorf("NewKafkaReceiver: %v", err)
		return
	}
	defer receiver.Close()
}

func TestKafkaReceiverGetMessageChan(t *testing.T) {
	conf.Get().KafkaBrokers = []string{"192.168.2.43:9094"}
	conf.Get().KafkaGroup = "example"
	conf.Get().KafkaTopic = "example"

	receiver, err := NewKafkaReceiver()
	if err != nil {
		t.Errorf("NewKafkaReceiver: %v", err)
		return
	}
	defer receiver.Close()

	ctx, cancel := context.WithCancel(context.TODO())
	msgCh := receiver.GetMessageChan(ctx)

	t.Logf("message: %+v", <-msgCh)
	cancel()
}
