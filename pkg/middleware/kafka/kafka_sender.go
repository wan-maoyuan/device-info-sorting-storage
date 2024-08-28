package kafka

import (
	"device-info-sorting-storage/pkg/conf"
	"encoding/json"
	"fmt"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"github.com/sirupsen/logrus"
)

var tmc *goka.TopicManagerConfig

func init() {
	tmc = goka.NewTopicManagerConfig()
	tmc.Table.Replication = 1
	tmc.Stream.Replication = 1
}

type KafkaSender struct {
	emitter *goka.Emitter
}

func NewKafkaSender() (*KafkaSender, error) {
	brokers := conf.Get().KafkaBrokers
	topic := goka.Stream(conf.Get().KafkaTopic)

	topicManager, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), tmc)
	if err != nil {
		return nil, fmt.Errorf("创建 TopicManager 失败: %v", err)
	}
	defer topicManager.Close()

	if err := topicManager.EnsureStreamExists(string(topic), 8); err != nil {
		return nil, fmt.Errorf("topic: %s 创建失败: %v", string(topic), err)
	}

	emitter, err := goka.NewEmitter(brokers, topic, new(codec.Bytes))
	if err != nil {
		return nil, fmt.Errorf("kafka emitter 创建失败: %v", err)
	}

	return &KafkaSender{
		emitter: emitter,
	}, nil
}

func (k *KafkaSender) SendMsg(key string, msg *Message) error {
	content, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化 Message 失败: %v", err)
	}

	return k.emitter.EmitSync(key, content)
}

func (k *KafkaSender) Close() {
	if err := k.emitter.Finish(); err != nil {
		logrus.Errorf("kafka 连接关闭失败: %v", err)
	}
}
