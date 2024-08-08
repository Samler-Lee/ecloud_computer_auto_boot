FROM golang:1.20-alpine as builder

COPY . /app
WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN go build -o ecloud_computer_auto_boot .

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
EXPOSE 10839

ENTRYPOINT ["./ecloud_computer_auto_boot"]