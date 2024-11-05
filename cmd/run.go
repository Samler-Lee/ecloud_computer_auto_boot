package cmd

import (
	"ecloud_computer_auto_boot/bootstrap"
	"ecloud_computer_auto_boot/pkg/task"
	"ecloud_computer_auto_boot/pkg/util"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the automatic boot cron",
	Long: `Run the automatic boot cron. 

If you run it on this device for the first time, you must first run the trust command.`,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.Init()

		// 收到信号后关闭服务器
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		sig := <-sigChan
		util.Log().Info("收到信号 %s, 开始关闭进程", sig)
		task.Destroy()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
