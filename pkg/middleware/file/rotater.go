package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// 文件轮换器
type FileRotater struct {
	country     string
	interval    time.Duration
	currentTime time.Time
	file        *os.File
	mu          sync.Mutex
}

func NewFileRotater(country string, interval time.Duration) *FileRotater {
	currentTime, _ := time.Parse("2006010215", time.Now().Format("2006010215"))

	return &FileRotater{
		country:     country,
		interval:    interval,
		currentTime: currentTime,
		mu:          sync.Mutex{},
	}
}

func (rotate *FileRotater) Write(line []byte) {
	rotate.mu.Lock()
	defer rotate.mu.Unlock()

	if rotate.file == nil {
		if err := rotate.openNew(); err != nil {
			logrus.Errorf("rotate 创建文件失败: %v", err)
			return
		}
	}

	rotate.updateCurrentTime()

	if _, err := rotate.file.Write(line); err != nil {
		logrus.Errorf(
			"country: %s time: %v write line error: %v",
			rotate.country, rotate.currentTime, err,
		)
	}

	rotate.file.Write([]byte("\n"))
}

func (rotate *FileRotater) updateCurrentTime() {
	now := time.Now().Local()
	changeFlag := false

	for rotate.currentTime.Local().Add(rotate.interval).Before(now) {
		rotate.currentTime = rotate.currentTime.Add(rotate.interval)
		changeFlag = true
	}

	// 时间更新过,需要切换文件
	if changeFlag {
		rotate.rotate()
	}
}

func (rotate *FileRotater) rotate() error {
	if err := rotate.file.Close(); err != nil {
		return fmt.Errorf("旧的缓存文件关闭失败: %v", err)
	}

	if err := rotate.openNew(); err != nil {
		return err
	}

	return nil
}

func (rotate *FileRotater) openNew() (err error) {
	now := time.Now()

	dir := filepath.Join(BaseDir, now.Local().Format("20060102"))
	if _, err = os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, os.FileMode(0755)); err != nil {
			return fmt.Errorf("无法为缓存文件创建文件夹: %s 错误: %v", dir, err)
		}
	}

	fileName := fmt.Sprintf("%s_%s.txt", rotate.country, rotate.currentTime.Format("2006010215"))
	filePath := filepath.Join(dir, fileName)

	rotate.file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("缓存文件: %s 创建失败: %v", filePath, err)
	}

	return nil
}

func (rotate *FileRotater) Close() {
	rotate.mu.Lock()
	defer rotate.mu.Unlock()

	rotate.file.Close()
}
