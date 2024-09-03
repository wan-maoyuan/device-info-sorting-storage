# device-info-sorting-storage
> 设备信息分类存储服务

## 测试环境准备
- `kafka` 
```shell
# 启动
docker-compose -f resources/compose/kafka-compose.yml up -d

# 停止
docker-compose -f resources/compose/kafka-compose.yml down
```

- `redis`
```shell
# 启动
docker-compose -f resources/compose/redis-compose.yml up -d

# 停止
docker-compose -f resources/compose/redis-compose.yml down
```

- `rabbit-mq`
```shell
# 启动
docker-compose -f resources/compose/rabbitmq-compose.yml up -d

# 停止
docker-compose -f resources/compose/rabbitmq-compose.yml down
```

## 配置
- 配置信息
```shell
# kafka 节点信息,多个节点用逗号分割: KAFKA_BROKERS="192.168.2.43:9094,192.168.2.43:9092"
export KAFKA_BROKERS="192.168.2.43:9094"
# kafka 主题
export KAFKA_TOPIC="test_topic"
# kafka 群组
export KAFKA_GROUP="test_group"

# 存放 安卓 ID的 redis 服务信息
export ANDROID_REDIS_URL="redis://default:redis123..@192.168.2.43:6379/0"
# 存放 IOS ID的 redis 服务信息
export IOS_REDIS_URL="redis://default:redis123..@192.168.2.43:6379/1"
# 文件分割时间间隔,单位为小时,需要整数
export FILE_INTERVAL_HOUR=1

# 是否需要发送到 rabbit mq中，0表示不发送，1表示发送
export IS_NEED_SEND_MQ=0
# rabbit_mq 链接 uri
export MQ_URI="amqp://rabbit:123456@192.168.2.43:5672/ads-vhost"
# 需要发送的 rabbit_mq 队列名称前缀
export MQ_QUEUE_PREFIX="demo-test"

# 是否需要保存日志，为空不保存，直接在控制台打印。可以填 "./logs/website-verification.log"
export LOG_FILE=""
# 日志等级：debug,info,warn,error,fatal,panic
export LOG_LEVEL="info"
# 每一个日志文件最大的 size ，不超过 100 M
export LOG_SIZE=100
# 日志文件最大的个数，不超过100个
export LOG_AGE=100
```

- 激活配置
```shell
source resources/.env_dev.sh
```

## 程序运行
```shell
make
source resources/.env_dev.sh
./dist/device-info-sorting-storage
```

## `docker` 镜像打包
```shell
make container
```

## `docker` 镜像启动
```yml
version: '3.8'

services:
  device-info-sorting-storage:
    image: device-info-sorting-storage:v0.0.1
    container_name: device-info-sorting-storage
    environment:
      KAFKA_BROKERS: "192.168.2.43:9094"
      KAFKA_TOPIC: "test_topic"
      KAFKA_GROUP: "test_group"
      ANDROID_REDIS_URL: "redis://default:redis123..@192.168.2.43:6379/0"
      IOS_REDIS_URL: "redis://default:redis123..@192.168.2.43:6379/1"
      FILE_INTERVAL_HOUR: 1
      IS_NEED_SEND_MQ: 0                                                   
      MQ_URI: "amqp://rabbit:123456@192.168.2.43:5672/ads-vhost"         
      MQ_QUEUE_PREFIX: "demo-test"
      LOG_FILE: ""
      LOG_LEVEL: "info"
      LOG_SIZE: 100
      LOG_AGE: 100
    volumes:
      - /etc/localtime:/etc/localtime
      - ./device_info:/device_info
```

- 启动 `docker` 镜像
```shell
docker-compose -f resources/compose/docker-compose.yml up -d
```

- 停止 `docker` 镜像
```shell
docker-compose -f resources/compose/docker-compose.yml down
```

## 代码测试
> 根据需求,修改测试代码的参数

```shell
go run test/send_message/main.go
```