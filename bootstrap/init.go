package bootstrap

import (
	"ecloud_computer_auto_boot/pkg/conf"
	"ecloud_computer_auto_boot/pkg/task"
)

func Init() {
	InitApplication()
	conf.Init()

	task.Init()
}
