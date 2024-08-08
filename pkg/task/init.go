package task

import (
	"ecloud_computer_auto_boot/pkg/conf"
	"ecloud_computer_auto_boot/pkg/util"
	"fmt"
	"github.com/robfig/cron/v3"
	"gitlab.ecloud.com/ecloud/ecloudsdkcomputer"
	"gitlab.ecloud.com/ecloud/ecloudsdkcore/config"
)

var (
	cronInstance *cron.Cron
	cloudClient  *ecloudsdkcomputer.Client
)

func Init() {
	util.Log().Info("[定时任务] 初始化中")

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(fmt.Sprintf("@every %ds", conf.Cron.Duration), machineMonitor)
	if err != nil {
		util.Log().Error("[定时任务] 任务创建失败")
		return
	}

	c.Start()
	client := ecloudsdkcomputer.NewClient(&config.Config{
		AccessKey: &conf.Secret.AccessKey,
		SecretKey: &conf.Secret.SecretKey,
		PoolId:    &conf.Secret.PoolId,
	})

	cronInstance = c
	cloudClient = client
	util.Log().Info("[定时任务] 初始化完毕")

	// 立即执行一次
	machineMonitor()
}

func Destroy() {
	cronInstance.Stop()
}
