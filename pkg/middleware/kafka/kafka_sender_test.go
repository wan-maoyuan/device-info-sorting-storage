package kafka

import (
	"device-info-sorting-storage/pkg/conf"
	"testing"
)

func TestNewKafkaSender(t *testing.T) {
	conf.Get().KafkaBrokers = []string{"192.168.2.43:9094"}
	conf.Get().KafkaGroup = "example"
	conf.Get().KafkaTopic = "example"

	sender, err := NewKafkaSender()
	if err != nil {
		t.Errorf("NewKafkaSender: %v", err)
		return
	}
	defer sender.Close()
}
