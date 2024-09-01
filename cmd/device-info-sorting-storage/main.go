package main

import (
	"context"
	"device-info-sorting-storage/pkg/conf"
	"device-info-sorting-storage/pkg/middleware"
	"device-info-sorting-storage/pkg/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	c := conf.New()
	c.Show()
}

func main() {
	if err := BeforeStartFunc(); err != nil {
		logrus.Errorf("服务资源初始化失败: %v", err)
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
	time.Sleep(time.Second * 5)

	if err := AfterStopFunc(); err != nil {
		logrus.Info("服务资源清理失败")
	}
}

func BeforeStartFunc() error {
	if err := middleware.InitMiddleware(); err != nil {
		return fmt.Errorf("InitMiddleware: %v", err)
	}

	return nil
}

func AfterStopFunc() error {
	middleware.CloseMiddleware()

	return nil
}
