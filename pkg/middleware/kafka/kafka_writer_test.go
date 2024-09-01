package kafka

import (
	"device-info-sorting-storage/pkg/conf"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewKafkaWriter(t *testing.T) {
	conf.Get().KafkaBrokers = []string{"192.168.2.43:9094"}
	conf.Get().KafkaGroup = "example"
	conf.Get().KafkaTopic = "example"

	writer, err := NewKafkaWriter()
	if err != nil {
		t.Errorf("NewKafkaWriter: %v", err)
		return
	}
	defer writer.Close()
}

func TestKafkaWriterWriteMessage(t *testing.T) {
	conf.Get().KafkaBrokers = []string{"192.168.2.43:9094"}
	conf.Get().KafkaGroup = "example"
	conf.Get().KafkaTopic = "example"

	writer, err := NewKafkaWriter()
	if err != nil {
		t.Errorf("NewKafkaWriter: %v", err)
		return
	}
	defer writer.Close()

	writer.WriteMessage("", &Message{
		Timestamp: time.Now().Format(time.DateTime),
		ADX:       "media",
		APP: struct {
			Cat    []string "json:\"cat\""
			Bundle string   "json:\"bundle\""
			Name   string   "json:\"name\""
			Ver    string   "json:\"ver\""
		}{
			Cat:    []string{"IAB1", "IAB9"},
			Bundle: "1617391485",
			Name:   "Block Blast！",
			Ver:    "4.1.4",
		},
		Device: struct {
			Country  string "json:\"country\""
			IP       string "json:\"ip\""
			Language string "json:\"language\""
			OS       string "json:\"os\""
			OSV      string "json:\"osv\""
			UA       string "json:\"ua\""
			IFA      string "json:\"ifa\""
		}{
			Country:  "Algeria",
			IP:       "85.140.24.103",
			Language: "ru",
			OS:       "IOS",
			OSV:      "17.0.3",
			UA:       "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
			IFA:      uuid.New().String(),
		},
	})
}
