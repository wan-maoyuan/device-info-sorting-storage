package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	ExpirationDuration = time.Hour * 24 * 7
)

type RedisCli struct {
	client *redis.Client
}

func NewRedisCli(url string) (*RedisCli, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("redis url 解析失败: %v", err)
	}

	return &RedisCli{
		client: redis.NewClient(opt),
	}, nil
}

func (r *RedisCli) SetDeviceID(ctx context.Context, key string) bool {
	cmd := r.client.SetNX(ctx, key, nil, ExpirationDuration)
	return cmd.Val()
}

func (r *RedisCli) Close() {
	if err := r.client.Close(); err != nil {
		logrus.Errorf("关闭 redis 客户端失败: %v", err)
	}
}
