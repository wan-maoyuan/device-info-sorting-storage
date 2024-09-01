# 程序编译镜像
FROM golang:1.22 as build

ENV GO111MODULE=on                               \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ./dist/device-info-sorting-storage ./cmd/device-info-sorting-storage

# 构建镜像
FROM alpine:3.20 as prod

ENV KAFKA_BROKERS="192.168.2.43:9094"                                   \
    KAFKA_TOPIC="test_topic"                                            \
    KAFKA_GROUP="test_group"                                            \
    KAFKA_PARTITION=0                                                   \
    ANDROID_REDIS_URL="redis://default:redis123..@192.168.2.43:6379/0"  \
    IOS_REDIS_URL="redis://default:redis123..@192.168.2.43:6379/1"      \
    FILE_INTERVAL_HOUR=1                                                \
    LOG_FILE=""                                                         \
    LOG_LEVEL="info"                                                    \
    LOG_SIZE=100                                                        \
    LOG_AGE=100                                     

# 在build阶段复制时区到
COPY --from=build /go/release/dist/device-info-sorting-storage /device-info-sorting-storage

CMD [ "/device-info-sorting-storage"]