package main

import (
	"device-info-sorting-storage/pkg/conf"
	"device-info-sorting-storage/pkg/middleware/kafka"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	JSON_PATH = "./resources/demo.json"
)

var countryList = []string{
	"Afghanistan",
	"Albania",
	"Algeria",
	"Andorra",
	"Angola",
	"Antigua and Barbuda",
	"Argentina",
	"Armenia",
	"Australia",
	"Austria",
	"Azerbaijan",
	"Bahamas",
	"Bahrain",
	"Bangladesh",
	"Barbados",
	"Belarus",
	"Belgium",
	"Belize",
	"Benin",
	"Bhutan",
	"Bolivia",
	"Bosnia and Herzegovina",
	"Botswana",
	"Brazil",
	"Brunei",
	"Bulgaria",
	"Burkina Faso",
	"Burundi",
	"Cabo Verde",
	"Cambodia",
	"Cameroon",
	"Canada",
	"Central African Republic",
	"Chad",
	"Chile",
	"China",
	"Colombia",
	"Comoros",
	"Congo",
	"Costa Rica",
	"Croatia",
	"Cuba",
	"Cyprus",
	"Czech Republic",
	"Denmark",
	"Djibouti",
	"Dominica",
	"Dominican Republic",
	"East Timor (Timor-Leste)",
	"Ecuador",
}

func init() {
	c := conf.New()
	c.Show()
}

func main() {
	writer, err := kafka.NewKafkaWriter()
	if err != nil {
		logrus.Error(err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100000000; i++ {
		for _, country := range countryList {
			msg := getIOSMessage(country)
			logrus.Debugf("generate a message: %v", msg)

			if err := writer.WriteMessage("", msg); err != nil {
				logrus.Error(err)
				continue
			}
		}
	}
}

func getIOSMessage(country string) *kafka.Message {
	var msg = &kafka.Message{
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
			Country:  country,
			IP:       "85.140.24.103",
			Language: "ru",
			OS:       "IOS",
			OSV:      "17.0.3",
			UA:       "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
			IFA:      uuid.New().String(),
		},
	}

	return msg
}
