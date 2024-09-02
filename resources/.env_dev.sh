export KAFKA_BROKERS="192.168.2.43:9094"
export KAFKA_TOPIC="test_topic"
export KAFKA_GROUP="test_group"

# export KAFKA_BROKERS="47.92.72.122:31594"
# export KAFKA_TOPIC="test_topic"
# export KAFKA_GROUP="test_group"

export ANDROID_REDIS_URL="redis://default:redis123..@192.168.2.43:6379/0"
export IOS_REDIS_URL="redis://default:redis123..@192.168.2.43:6379/1"

export IS_NEED_SEND_MQ=0
export MQ_URI="amqp://rabbit:123456@192.168.2.43:5672/ads-vhost"
export MQ_QUEUE="demo-test"

export FILE_INTERVAL_HOUR=1

export LOG_FILE=""
export LOG_LEVEL="info"
export LOG_SIZE=100
export LOG_AGE=100