package middleware

import (
	"device-info-sorting-storage/pkg/middleware/kafka"
	"fmt"
)

var middle = new(Middleware)

type Middleware struct {
	kafkaReceiver *kafka.KafkaReceiver
}

func InitMiddleware() (err error) {
	middle.kafkaReceiver, err = kafka.NewKafkaReceiver()
	if err != nil {
		return fmt.Errorf("NewKafkaReceiver: %v", err)
	}

	return nil
}

func CloseMiddleware() {
	middle.kafkaReceiver.Close()
}
