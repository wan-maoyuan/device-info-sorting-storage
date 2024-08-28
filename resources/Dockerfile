# 程序编译镜像
FROM golang:1.22 as build

ENV GO111MODULE=on                               \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ./dist/device-info-sorting-storage ./cmd/device-info-sorting-storage

# 构建镜像
FROM alpine:3.20 as prod

ENV MQ_URI="amqp://rabbit:123456@127.0.0.1:5672/"   \
    MQ_QUEUE="test"                                 \
    LOG_FILE=""                                     \
    LOG_LEVEL="info"                                \
    LOG_SIZE=100                                    \
    LOG_AGE=100                                     \
    CONCURRENT=100                                  \
    TIME_OUT_SECOND=10                              \
    STATISTICS_INTERVAL_MINUTE=5

# 在build阶段复制时区到
COPY --from=build /go/release/dist/device-info-sorting-storage /device-info-sorting-storage

CMD [ "/device-info-sorting-storage"]