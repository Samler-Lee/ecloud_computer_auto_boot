# ecloud computer auto boot
移动云电脑自动开机

# 配置
默认读取运行目录中的config.yml作为配置，如果无配置文件，首次运行将会生成一个默认的配置文件

## 示例配置
```yaml
cron:
    # 任务执行间隔（秒）
    duration: 60
    # 需要监控的实例 machine id
    machine_list:
        - machine_id_1
        - machine_id_2
secret:
    # 移动云 Access Key
    access_key: ""
    # 移动云 Secret Key
    secret_key: ""
    # 资源池ID
    pool_id: "CIDC-CORE-00"
server:
    # 是否开启 Debug 模式
    debug: false
    # 日志输出等级，等级划分: debug、info、warning、error，开启 Debug 模式时，此配置无效
    log_level: info

```

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

# 在Docker中使用（推荐）
## 1、构建
```shell
docker build -t ecloud_computer_auto_boot:latest .
```

## 2、运行
```shell
docker run -itd --restart=always -v /path_to_config/config.yml:/app/config.yml ecloud_computer_auto_boot:latest
```