package middleware

import (
	"device-info-sorting-storage/pkg/conf"
	"device-info-sorting-storage/pkg/middleware/file"
	"device-info-sorting-storage/pkg/middleware/kafka"
	"device-info-sorting-storage/pkg/middleware/redis"
	"fmt"
)

var middle = new(Middleware)

type Middleware struct {
	fileSave     *file.FileSaver
	kafkaReader  *kafka.KafkaReader
	androidRedis *redis.RedisCli
	iosRedis     *redis.RedisCli
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

	return nil
}

func CloseMiddleware() {
	middle.fileSave.Close()
	middle.kafkaReader.Close()
	middle.androidRedis.Close()
	middle.iosRedis.Close()
}
