package file

import "os"

const (
	BaseDir = "device_info"
)

func init() {
	if _, err := os.Stat(BaseDir); err != nil {
		os.Mkdir(BaseDir, 0777)
	}
}
