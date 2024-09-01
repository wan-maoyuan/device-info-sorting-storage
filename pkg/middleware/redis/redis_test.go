package redis

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestNewRedisCli(t *testing.T) {
	url := "redis://default:redis123..@192.168.2.43:6379/0"

	client, err := NewRedisCli(url)
	if err != nil {
		t.Errorf("NewRedisCli: %v", err)
		return
	}
	defer client.Close()
}

func TestSetDeviceID(t *testing.T) {
	url := "redis://default:redis123..@192.168.2.43:6379/0"

	client, err := NewRedisCli(url)
	if err != nil {
		t.Errorf("NewRedisCli: %v", err)
		return
	}
	defer client.Close()

	id := uuid.New().String()
	if !client.SetDeviceID(context.TODO(), id) {
		t.Errorf("SetDeviceID: %s 失败", id)
		return
	}

	id = "4917d044-f739-45f6-9836-ba2d070815b1"
	if !client.SetDeviceID(context.TODO(), id) {
		t.Errorf("SetDeviceID: %s 失败", id)
		return
	}
}
