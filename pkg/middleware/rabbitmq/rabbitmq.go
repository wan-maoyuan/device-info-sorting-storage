package rabbitmq

import (
	"device-info-sorting-storage/pkg/conf"
	"device-info-sorting-storage/pkg/middleware/kafka"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var RabbitConcurrentCount int = 100

type Rabbit struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	queueMap sync.Map
}

func NewRabbitmq() (*Rabbit, error) {
	config := conf.Get()

	conn, err := amqp.Dial(config.MQURI)
	if err != nil {
		return nil, fmt.Errorf("连接 rabbitmq 服务器失败: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("创建一个 rabbitmq 连接通道失败: %v", err)
	}

	return &Rabbit{
		conn:     conn,
		channel:  channel,
		queueMap: sync.Map{},
	}, nil
}

func (mq *Rabbit) SendMessage(msgChan chan *kafka.Message) {
	wg := &sync.WaitGroup{}

	for i := 0; i < RabbitConcurrentCount; i++ {
		wg.Add(1)
		go mq.parseMessageAndSend(msgChan, wg)
	}

	wg.Wait()
}

func (mq *Rabbit) parseMessageAndSend(msgChan chan *kafka.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range msgChan {
		key := fmt.Sprintf("%s_%s", msg.Device.Country, msg.Device.OS)
		if _, ok := mq.queueMap.Load(key); !ok {
			queue, err := mq.channel.QueueDeclare(
				fmt.Sprintf("%s_%s_%s", conf.Get().MQQueuePrefix, msg.Device.Country, msg.Device.OS), // 队列名
				true,  // 是否持续
				false, // 是否自动删除
				false, // 是否独占
				false, // 是否阻塞
				nil,   // args
			)
			if err != nil {
				logrus.Errorf("创建 rabbit_mq 消息队列: %s 失败: %v", queue.Name, err)
				continue
			}

			mq.queueMap.Store(key, queue)
		}

		queue, _ := mq.queueMap.Load(key)
		publishQueue := queue.(amqp.Queue)

		body, err := json.Marshal(msg)
		if err != nil {
			logrus.Errorf("解析 Message 结构体 to json 失败: %v", err)
			continue
		}

		err = mq.channel.Publish(
			"",                // 交换机的名称
			publishQueue.Name, // 需要发送的消息队列
			false,             // 消息发送失败是否需要收到回复
			false,             // 设置为true，当消息无法直接投递到消费者时，会返回一个Basic.Return消息给生产者。如果设置为false，则消息会被存储在队列中，等待消费者连接。
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})

		if err != nil {
			logrus.Errorf("往 rabbit_mq 发送消息: %s 失败: %v", string(body), err)
			continue
		}
	}
}

func (mq *Rabbit) Close() {
	if mq.channel != nil {
		if err := mq.channel.Close(); err != nil {
			logrus.Errorf("关闭 rabbitmq 连接通道失败: %v", err)
		}
	}

	if mq.conn != nil {
		if err := mq.conn.Close(); err != nil {
			logrus.Errorf("关闭 rabbitmq 连接失败: %v", err)
		}
	}
}
