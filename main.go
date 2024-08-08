package main

import (
	"ecloud_computer_auto_boot/bootstrap"
	"ecloud_computer_auto_boot/pkg/task"
	"ecloud_computer_auto_boot/pkg/util"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	bootstrap.Init()
}

func main() {
	// 收到信号后关闭服务器
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-sigChan
	util.Log().Info("收到信号 %s, 开始关闭进程", sig)
	task.Destroy()
}
