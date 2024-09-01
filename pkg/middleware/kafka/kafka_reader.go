package kafka

import (
	"context"
	"device-info-sorting-storage/pkg/conf"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

const (
	ChannelCapacity int = 1
)

type Message struct {
	Timestamp string `json:"@timestamp"`
	ADX       string `json:"adx"`
	APP       struct {
		Cat    []string `json:"cat"`
		Bundle string   `json:"bundle"`
		Name   string   `json:"name"`
		Ver    string   `json:"ver"`
	} `json:"app"`
	Device struct {
		Country  string `json:"country"`
		IP       string `json:"ip"`
		Language string `json:"language"`
		OS       string `json:"os"`
		OSV      string `json:"osv"`
		UA       string `json:"ua"`
		IFA      string `json:"ifa"`
	} `json:"device"`
}

type KafkaReader struct {
	reader *kafka.Reader
}

func NewKafkaReader() (*KafkaReader, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: conf.Get().KafkaBrokers,
		GroupID: conf.Get().KafkaGroup,
		Topic:   conf.Get().KafkaTopic,
		MaxWait: time.Second,
		// StartOffset: kafka.FirstOffset,
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaReader{
		reader: reader,
	}, nil
}

func (k *KafkaReader) ReadMessage(ctx context.Context) chan *Message {
	ch := make(chan *Message, ChannelCapacity)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				logrus.Info("停止接收 kafka 消息")
				return
			default:
				msg, err := k.reader.ReadMessage(context.TODO())
				if err != nil {
					logrus.Errorf("kafka 读取消息失败: %v", err)
					continue
				}

				var message = new(Message)
				if err := json.Unmarshal(msg.Value, &message); err != nil {
					logrus.Errorf("kafka 消息反序列化失败: %v", err)
					continue
				}

				ch <- message
				logrus.Debugf("kafka 收到一条消息: %s", string(msg.Value))
			}
		}
	}()

	return ch
}

func (k *KafkaReader) Close() {
	if err := k.reader.Close(); err != nil {
		logrus.Errorf("kafka 读取服务关闭失败: %v", err)
	}
}
