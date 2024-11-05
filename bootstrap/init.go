package bootstrap

import (
	"ecloud_computer_auto_boot/pkg/conf"
	"ecloud_computer_auto_boot/pkg/task"
	"ecloud_computer_auto_boot/pkg/util"
	"net/http"
	_ "net/http/pprof"
)

func Init() {
	InitApplication()
	conf.Init()

	if conf.Server.Debug {
		go func() {
			if err := http.ListenAndServe(":8080", nil); err != nil {
				util.Log().Error("性能分析工具启动失败: %s", err)
				return
			}

			util.Log().Info("性能分析工具已开启")
		}()
	}

	task.Init()
}
