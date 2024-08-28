package redis

import (
	"device-info-sorting-storage/pkg/conf"
	"testing"
)

func TestNewRedisCli(t *testing.T) {
	conf.Get().RedisUrl = "redis://default:redis123..@192.168.2.43:6379/0"

	client, err := NewRedisCli()
	if err != nil {
		t.Errorf("NewRedisCli: %v", err)
		return
	}
	defer client.Close()
}
