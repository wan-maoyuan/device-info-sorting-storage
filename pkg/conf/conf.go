package conf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var config = new(Conf)

func Get() *Conf {
	return config
}

type Conf struct {
	KafkaBrokers     []string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic       string   `mapstructure:"KAFKA_TOPIC"`
	KafkaGroup       string   `mapstructure:"KAFKA_GROUP"`
	KafkaPartition   int      `mapstructure:"KAFKA_PARTITION"`
	AndroidRedisUrl  string   `mapstructure:"ANDROID_REDIS_URL"`
	IOSRedisUrl      string   `mapstructure:"IOS_REDIS_URL"`
	FileIntervalHour uint     `mapstructure:"FILE_INTERVAL_HOUR"`
	Log              Log      `mapstructure:"LOG"`
}

func New() *Conf {
	viper.AutomaticEnv()

	config.KafkaBrokers = strings.Split(viper.GetString("KAFKA_BROKERS"), ",")
	config.KafkaTopic = viper.GetString("KAFKA_TOPIC")
	config.KafkaGroup = viper.GetString("KAFKA_GROUP")
	config.KafkaPartition = viper.GetInt("KAFKA_PARTITION")

	config.AndroidRedisUrl = viper.GetString("ANDROID_REDIS_URL")
	config.IOSRedisUrl = viper.GetString("IOS_REDIS_URL")

	config.FileIntervalHour = viper.GetUint("FILE_INTERVAL_HOUR")

	config.Log = NewLog()

	config.Log.File = viper.GetString("LOG_FILE")
	config.Log.Level = viper.GetString("LOG_LEVEL")
	config.Log.MaxSize = viper.GetInt("LOG_SIZE")
	config.Log.MaxAge = viper.GetInt("LOG_AGE")

	config.Log.InitLog()

	return config
}

func (c *Conf) Show() {
	if b, err := yaml.Marshal(c); err != nil {
		return
	} else {
		fmt.Printf(`
-----------------------------------------------------------------------------------------
%s
-----------------------------------------------------------------------------------------
`, string(b))
	}
}
