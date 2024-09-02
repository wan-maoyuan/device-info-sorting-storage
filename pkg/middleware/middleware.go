package middleware

import (
	"device-info-sorting-storage/pkg/conf"
	"device-info-sorting-storage/pkg/middleware/file"
	"device-info-sorting-storage/pkg/middleware/kafka"
	"device-info-sorting-storage/pkg/middleware/rabbitmq"
	"device-info-sorting-storage/pkg/middleware/redis"
	"fmt"
)

var middle = new(Middleware)

type Middleware struct {
	fileSave     *file.FileSaver
	kafkaReader  *kafka.KafkaReader
	androidRedis *redis.RedisCli
	iosRedis     *redis.RedisCli
	rabbit       *rabbitmq.Rabbit
	rabbitChan   chan *kafka.Message
}

func InitMiddleware() (err error) {
	middle.fileSave, err = file.NewFileSaver()
	if err != nil {
		return fmt.Errorf("NewFileSaver: %v", err)
	}

	middle.kafkaReader, err = kafka.NewKafkaReader()
	if err != nil {
		return fmt.Errorf("NewKafkaReader: %v", err)
	}

	middle.androidRedis, err = redis.NewRedisCli(conf.Get().AndroidRedisUrl)
	if err != nil {
		return fmt.Errorf("连接安卓 redis 失败: %v", err)
	}

	middle.iosRedis, err = redis.NewRedisCli(conf.Get().IOSRedisUrl)
	if err != nil {
		return fmt.Errorf("连接苹果 redis 失败: %v", err)
	}

	if conf.Get().IsNeedSendMQ {
		middle.rabbit, err = rabbitmq.NewRabbitmq()
		if err != nil {
			return fmt.Errorf("连接 rabbit_mq 失败: %v", err)
		}

		middle.rabbitChan = make(chan *kafka.Message, 100)
		middle.rabbit.SendMessage(middle.rabbitChan)
	}

	return nil
}

func CloseMiddleware() {
	if middle.fileSave != nil {
		middle.fileSave.Close()
	}

	if middle.kafkaReader != nil {
		middle.kafkaReader.Close()
	}

	if middle.androidRedis != nil {
		middle.androidRedis.Close()
	}

	if middle.iosRedis != nil {
		middle.iosRedis.Close()
	}

	if middle.rabbitChan != nil {
		close(middle.rabbitChan)
	}

	if middle.rabbit != nil {
		middle.rabbit.Close()
	}
}
