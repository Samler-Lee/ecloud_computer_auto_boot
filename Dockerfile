FROM golang:1.23.2-alpine AS builder

COPY . /app
WORKDIR /app

ENV GO111MODULE=on
RUN go env -w GOPROXY=https://ecloud.10086.cn/api/query/developer/nexus/repository/go-sdk/ && \
    go env -w GONOSUMDB=gitlab.ecloud.com && \
    go get -u gitlab.ecloud.com/ecloud/ecloudsdkcomputer

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go build -o ecloud_computer_auto_boot .

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

ENV TZ=Asia/Shanghai
RUN apk update \
	&& apk add tzdata \
	&& echo "${TZ}" > /etc/timezone \
	&& ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
	&& rm /var/cache/apk/*

WORKDIR /app
COPY --from=builder /app/ecloud_computer_auto_boot ./

RUN chmod u+x ./ecloud_computer_auto_boot

ENTRYPOINT ["./ecloud_computer_auto_boot", "run"]