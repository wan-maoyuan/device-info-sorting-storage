package file

import (
	"device-info-sorting-storage/pkg/conf"
	"device-info-sorting-storage/pkg/middleware/kafka"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	BaseDir = "device_info"
)

type Country string

type FileSaver struct {
	fileMap  sync.Map
	interval time.Duration
}

func NewFileSaver() (*FileSaver, error) {
	if _, err := os.Stat(BaseDir); err != nil {
		os.Mkdir(BaseDir, os.FileMode(0755))
	}

	interval := time.Duration(conf.Get().FileIntervalHour) * time.Hour
	logrus.Debugf("file saver interval: %v", interval)

	return &FileSaver{
		fileMap:  sync.Map{},
		interval: interval,
	}, nil
}

func (f *FileSaver) WriteMessage2File(msg *kafka.Message) {
	var rotater *FileRotater

	value, ok := f.fileMap.Load(msg.Device.Country)
	if ok {
		rotater = value.(*FileRotater)
	} else {
		rotater = NewFileRotater(msg.Device.Country, f.interval)
		f.fileMap.Store(msg.Device.Country, rotater)
	}

	lineBytes, err := json.Marshal(&msg)
	if err != nil {
		return
	}

	rotater.Write(lineBytes)
}

func (f *FileSaver) Close() {
	f.fileMap.Range(func(key, value any) bool {
		rotater := value.(*FileRotater)
		rotater.Close()
		return true
	})
}
