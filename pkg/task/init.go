package task

import (
	"ecloud_computer_auto_boot/pkg/conf"
	"ecloud_computer_auto_boot/pkg/ecloud"
	"ecloud_computer_auto_boot/pkg/util"
	"fmt"
	"github.com/robfig/cron/v3"
	"gitlab.ecloud.com/ecloud/ecloudsdkcomputer"
	"gitlab.ecloud.com/ecloud/ecloudsdkcore/config"
	"os"
)

var (
	cronInstance *cron.Cron
	apiClient    *ecloudsdkcomputer.Client
	publicClient *ecloud.Client
)

func Init() {
	util.Log().Info("[定时任务] 初始化中")

	if conf.Secret.Type == "public" {
		client, err := ecloud.NewClient(conf.Secret.Username, conf.Secret.Password)
		if err != nil {
			util.Log().Error("客户端创建失败: %s", err)
			os.Exit(1)
		}

		publicClient = client

		if _, err = publicClient.Login(); err != nil {
			util.Log().Error("登录失败: %s", err)
			os.Exit(1)
		}

		if !publicClient.HasTrustDeviceRecord() {
			util.Log().Error("登录账号未受信任, 无法进行监控, 请先运行 trust 命令进行信任")
			os.Exit(1)
		}

		_, err = publicClient.VerifyAccessTicket()
		if err != nil {
			util.Log().Error("验证访问票据失败: %s", err)
			os.Exit(1)
		}

		_, err = publicClient.RecordDeviceInfo()
		if err != nil {
			util.Log().Error("记录登录设备信息失败: %s", err)
			os.Exit(1)
		}
	} else {
		client := ecloudsdkcomputer.NewClient(&config.Config{
			AccessKey: &conf.Secret.AccessKey,
			SecretKey: &conf.Secret.SecretKey,
			PoolId:    &conf.Secret.PoolId,
		})

		apiClient = client
	}

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(fmt.Sprintf("@every %ds", conf.Cron.Duration), launchMonitor)
	if err != nil {
		util.Log().Error("[定时任务] 任务创建失败")
		os.Exit(1)
	}

	c.Start()

	cronInstance = c
	util.Log().Info("[定时任务] 初始化完毕")

	// 立即执行一次
	launchMonitor()
}

func launchMonitor() {
	if conf.Secret.Type == "public" {
		startMachineMonitorOnPublic()
	} else {
		startMachineMonitorOnOpenAPI()
	}
}

func Destroy() {
	cronInstance.Stop()
}
