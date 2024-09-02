package middleware

import (
	"context"
	"device-info-sorting-storage/pkg/middleware/kafka"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	IOS_OS           string = "IOS"
	ANDROID_OS       string = "android"
	ConcurrencyCount int    = 100
)

// 从 kafka 中读取设备信息，保存到 redis (7天内已经有的就不存)
func HandleDeviceInfo(ctx context.Context) error {
	msgCh := middle.kafkaReader.ReadMessage(ctx)
	wg := &sync.WaitGroup{}
	ch := make(chan struct{}, ConcurrencyCount)
	defer close(ch)

	for {
		select {
		case <-ctx.Done():
			logrus.Info("kafka 消息处理任务停止")
			wg.Wait()
			return nil
		case msg := <-msgCh:
			if msg == nil {
				continue
			}

			handleMessage(msg, wg, ch)
		}
	}
}

func handleMessage(msg *kafka.Message, wg *sync.WaitGroup, ch chan struct{}) {
	wg.Add(1)
	ch <- struct{}{}

	go checkAndSave(msg, wg, ch)
}

// 去 redis 检查 7 天内是否已经存在该设备信息，去重
func checkAndSave(msg *kafka.Message, wg *sync.WaitGroup, ch chan struct{}) {
	defer func() {
		wg.Done()
		<-ch
	}()

	switch msg.Device.OS {
	case IOS_OS:
		if middle.iosRedis.SetDeviceID(context.TODO(), msg.Device.IFA) {
			logrus.Debugf("收到一个 iOS 设备信息: %s", msg.Device.IFA)
			middle.fileSave.WriteMessage2File(msg)

			if middle.rabbitChan != nil {
				middle.rabbitChan <- msg
			}
		}
	case ANDROID_OS:
		if middle.androidRedis.SetDeviceID(context.TODO(), msg.Device.IFA) {
			logrus.Debugf("收到一个 Android 设备信息: %s", msg.Device.IFA)
			middle.fileSave.WriteMessage2File(msg)

			if middle.rabbitChan != nil {
				middle.rabbitChan <- msg
			}
		}
	default:
		logrus.Errorf("设备系统信息: %s 错误, 需要: %s %s json: %+v", msg.Device.OS, IOS_OS, ANDROID_OS, msg)
	}
}
