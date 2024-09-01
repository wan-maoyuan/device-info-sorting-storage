package kafka

import (
	"context"
	"device-info-sorting-storage/pkg/conf"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter() (*KafkaWriter, error) {
	brokers := conf.Get().KafkaBrokers
	topic := conf.Get().KafkaTopic

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           time.Second,
		AllowAutoTopicCreation: true,
		Async:                  true,
		RequiredAcks:           kafka.RequiredAcks(0),
	}

	return &KafkaWriter{
		writer: writer,
	}, nil
}

func (k *KafkaWriter) WriteMessage(key string, msg *Message) (err error) {
	if msg == nil {
		return fmt.Errorf("发送到 kafka 的消息为空")
	}

	value, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("发送到 kafka 的消息序列化失败: %v", err)
	}

	err = k.writer.WriteMessages(context.TODO(), kafka.Message{
		Key:   []byte(key),
		Value: value,
	})

	if err != nil {
		return fmt.Errorf("发送到 kafka 的消息失败: %v", err)
	}

	return nil
}

func (k *KafkaWriter) Close() {
	if err := k.writer.Close(); err != nil {
		logrus.Errorf("关闭 kafka writer 失败: %v", err)
	}
}
