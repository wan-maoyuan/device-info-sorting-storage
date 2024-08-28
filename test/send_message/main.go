package main

import (
	"device-info-sorting-storage/pkg/conf"
	"device-info-sorting-storage/pkg/middleware/kafka"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	JSON_PATH = "./resources/demo.json"
)

func init() {
	c := conf.New()
	c.Show()
}

func main() {
	sender, err := kafka.NewKafkaSender()
	if err != nil {
		logrus.Error(err)
		return
	}
	defer sender.Close()

	msg, err := getDemoMessage()
	if err != nil {
		logrus.Error(err)
		return
	}

	for i := 0; i < 10000; i++ {
		if err := sender.SendMsg("", msg); err != nil {
			logrus.Error(err)
			continue
		}
	}
}

func getDemoMessage() (*kafka.Message, error) {
	content, err := os.ReadFile(JSON_PATH)
	if err != nil {
		return nil, fmt.Errorf("读取 demo.json 文件失败")
	}

	var msg = new(kafka.Message)
	if err := json.Unmarshal(content, &msg); err != nil {
		return nil, fmt.Errorf("demo.json 文件反序列化 mesaage 失败: %v", err)
	}

	return msg, nil
}
