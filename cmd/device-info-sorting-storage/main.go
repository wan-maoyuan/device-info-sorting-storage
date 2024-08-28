package main

import (
	"context"
	"device-info-sorting-storage/pkg/conf"
	"device-info-sorting-storage/pkg/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func init() {
	c := conf.New()
	c.Show()
}

func main() {
	if err := BeforeStartFunc(); err != nil {
		logrus.Errorf("服务初始化失败: %v", err)
		return
	}

	server, err := server.NewServer()
	if err != nil {
		logrus.Errorf("服务创建失败: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.TODO())
	finishCh := make(chan os.Signal, 1)

	go func() {
		if err := server.Run(ctx); err != nil {
			logrus.Errorf("服务启动失败: %v", err)
			finishCh <- syscall.SIGQUIT
		}
	}()

	signal.Notify(finishCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-finishCh

	cancel()
	AfterStopFunc()
	logrus.Info("服务停止")
}

func BeforeStartFunc() error {

	return nil
}

func AfterStopFunc() error {

	return nil
}
