package redis

import (
	"device-info-sorting-storage/pkg/conf"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisCli struct {
	client *redis.Client
}

func NewRedisCli() (*RedisCli, error) {
	opt, err := redis.ParseURL(conf.Get().RedisUrl)
	if err != nil {
		return nil, fmt.Errorf("redis url 解析失败: %v", err)
	}

	client := redis.NewClient(opt)

	return &RedisCli{
		client: client,
	}, nil
}

func (r *RedisCli) Close() {
	if err := r.client.Close(); err != nil {
		logrus.Errorf("关闭 redis 客户端失败: %v", err)
	}
}
