package kafka

import (
	"context"
	"device-info-sorting-storage/pkg/conf"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
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

type KafkaReceiver struct {
	ch        chan *Message
	processor *goka.Processor
}

func NewKafkaReceiver() (*KafkaReceiver, error) {
	ch := make(chan *Message, ChannelCapacity)

	callback := func(ctx goka.Context, msg any) {
		var message = new(Message)
		if err := json.Unmarshal(msg.([]byte), &message); err != nil {
			logrus.Errorf("解析 kafka 消息失败: %v", err)
			return
		}

		ch <- message
	}

	config := goka.DefaultConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	goka.ReplaceGlobalConfig(config)

	topicManager, err := goka.NewTopicManager(conf.Get().KafkaBrokers, goka.DefaultConfig(), tmc)
	if err != nil {
		return nil, fmt.Errorf("创建 TopicManager 失败: %v", err)
	}
	defer topicManager.Close()

	if err := topicManager.EnsureStreamExists(string(conf.Get().KafkaTopic), 8); err != nil {
		return nil, fmt.Errorf("topic: %s 创建失败: %v", string(conf.Get().KafkaTopic), err)
	}

	groupGraph := goka.DefineGroup(
		goka.Group(conf.Get().KafkaGroup),
		goka.Input(goka.Stream(conf.Get().KafkaTopic), new(codec.Bytes), callback),
		goka.Persist(new(codec.Int64)),
	)

	processor, err := goka.NewProcessor(
		conf.Get().KafkaBrokers,
		groupGraph,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(tmc)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		return nil, fmt.Errorf("连接 kafka 创建 Processor 失败: %v", err)
	}

	return &KafkaReceiver{
		ch:        ch,
		processor: processor,
	}, nil
}

func (k *KafkaReceiver) GetMessageChan(ctx context.Context) chan *Message {
	go func() {
		defer close(k.ch)

		if err := k.processor.Run(ctx); err != nil {
			logrus.Errorf("kafka Processor 消息监听器停止运行: %v", err)
		}
	}()

	return k.ch
}

func (k *KafkaReceiver) Close() {

}
