package middleware

import (
	"context"
	"device-info-sorting-storage/pkg/middleware/kafka"

	"github.com/sirupsen/logrus"
)

const (
	IOS_OS     string = "IOS"
	ANDROID_OS string = "android"
)

// 从 kafka 中读取设备信息，保存到 redis (7天内已经有的就不存)
func HandleDeviceInfo(ctx context.Context) error {
	msgCh := middle.kafkaReader.ReadMessage(ctx)

	for {
		select {
		case <-ctx.Done():
			logrus.Info("kafka 消息处理任务停止")
			return nil
		case msg := <-msgCh:
			if msg == nil {
				continue
			}

			handleMessage(msg)
		}
	}
}

func handleMessage(msg *kafka.Message) {
	switch msg.Device.OS {
	case IOS_OS:
		if middle.iosRedis.SetDeviceID(context.TODO(), msg.Device.IFA) {
			// 设备ID写入成功保存到文件中
			logrus.Debugf("收到一个 iOS 设备信息: %s", msg.Device.IFA)
			middle.fileSave.WriteMessage2File(msg)
		}
	case ANDROID_OS:
		if middle.androidRedis.SetDeviceID(context.TODO(), msg.Device.IFA) {
			// 设备ID写入成功保存到文件中
			logrus.Debugf("收到一个 Android 设备信息: %s", msg.Device.IFA)
			// middle.fileSave.WriteMessage2File(msg)
		}
	default:
		logrus.Errorf("设备系统信息: %s 错误, 需要: %s %s json: %+v", msg.Device.OS, IOS_OS, ANDROID_OS, msg)
	}
}
