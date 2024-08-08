# ecloud computer auto boot
移动云电脑自动开机

# 二进制构建

## 1、安装移动云电脑SDK
```shell
go env -w GOPROXY=https://ecloud.10086.cn/api/query/developer/nexus/repository/go-sdk/
go env -w GONOSUMDB=gitlab.ecloud.com
go get -u gitlab.ecloud.com/ecloud/ecloudsdkcomputer
```

## 2、安装其它依赖
```shell
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GONOSUMDB=
go mod tidy
```

## 3、运行构建命令
```shell
go build .
```

# 运行
Windows
```shell
./ecloud_computer_auto_boot.exe
```

Linux
```shell
./ecloud_computer_auto_boot
```

# Docker镜像构建
```shell
docker build -t ecloud_computer_auto_boot:latest .
```