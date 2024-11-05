# ecloud computer auto boot
移动云电脑自动开机

# 配置
默认读取运行目录中的config.yml作为配置，如果无配置文件，首次运行将会生成一个默认的配置文件

## 示例配置
```yaml
cron:
    # 任务执行间隔（秒）
    duration: 60
    # 需要监控的实例 machine id, 如果为空，则监控所有实例
    # machines: []
    machines:
        - machine_id_1
        - machine_id_2
secret:
    # 客户端类型, public: 公众版, business: 政企版
    type: public
    # [公众版专用] 登录账号
    username: ""
    # [公众版专用] 登录密码
    password: ""
    # [政企版专用] 移动云 Access Key
    access-key: ""
    # [政企版专用] 移动云 Secret Key
    secret-key: ""
    # [政企版专用] 资源池ID
    pool-id: "CIDC-CORE-00"
server:
    # 是否开启 Debug 模式
    debug: false
    # 日志输出等级，等级划分: debug、info、warning、error，开启 Debug 模式时，此配置无效
    log-level: info

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

## 信任设备
首次运行时，请先运行该命令，在本地运行或服务器上执行均可

**注意：你需要先在配置文件中配置好您的账号密码信息**

Windows
```shell
./ecloud_computer_auto_boot.exe trust
```

Linux
```shell
./ecloud_computer_auto_boot trust
```

## 获取设备列表信息
如果您不知道哪里获取设备ID，可以运行该命令获取设备列表信息

Windows
```shell
./ecloud_computer_auto_boot.exe list-machines
```

Linux
```shell
./ecloud_computer_auto_boot list-machines
```

## 开始监控
Windows
```shell
./ecloud_computer_auto_boot.exe run
```

Linux
```shell
./ecloud_computer_auto_boot run
```

# 在Docker中使用（推荐）
在使用镜像运行之前，请您确认您的账号是否能够通过设备信任检查，如果不能，可以先在本地执行信任设备命令。

**注意：请您先准备好配置文件，以便在运行时挂载到容器中**

## 使用预构建的镜像
```shell
docker run -itd --restart=always -v /path_to_config/config.yml:/app/config.yml --name ecloud_computer_auto_boot samlerlee/ecloud_computer_auto_boot:latest
```

## 自行构建
### 1、构建
```shell
docker build -t ecloud_computer_auto_boot:latest .
```

### 2、运行
```shell
docker run -itd --restart=always -v /path_to_config/config.yml:/app/config.yml --name ecloud_computer_auto_boot ecloud_computer_auto_boot:latest
```