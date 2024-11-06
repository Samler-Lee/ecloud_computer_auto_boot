package cmd

import (
	"ecloud_computer_auto_boot/pkg/conf"
	"ecloud_computer_auto_boot/pkg/ecloud"
	"ecloud_computer_auto_boot/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
)

// trustCmd represents the trust command
var trustCmd = &cobra.Command{
	Use:   "trust",
	Short: "Trust your device",
	Long:  `Trust your device when you first use this program.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init()
		if conf.Secret.Type != "public" {
			fmt.Printf("仅支持公众版账号进行设备信任\n如果您是公众版账号，请将配置文件中的 secret.type 字段修改为 public，并配置好 secret.username 和 secret.password\n")
			return
		}

		client, err := ecloud.NewClient(conf.Secret.Username, conf.Secret.Password)
		if err != nil {
			util.Log().Error("客户端创建失败: %s", err)
			return
		}

		if _, err = client.Login(); err != nil {
			util.Log().Error("登录失败: %s", err)
			return
		}

		if client.HasTrustDeviceRecord() {
			util.Log().Info("该设备已在登录账号的信任列表中，无需再次运行该命令")
			return
		}

		resp, err := client.SendTrustDeviceVerifySms()
		if err != nil {
			util.Log().Error("发送验证码失败: %s", err)
			return
		}

		util.Log().Info("已向 %s 发送验证码, 请在 %d 秒内输入验证码，并按回车键确认: ", client.GetSession().Mobile, int(resp.Body.(map[string]interface{})["expireTime"].(float64)))

		var code string
		_, err = fmt.Scanf("%s", &code)
		if err != nil {
			util.Log().Error("输入验证码失败: %s", err)
			return
		}

		if resp, err = client.TrustDevice(code); err != nil {
			util.Log().Error("信任设备失败: %s", err)
			return
		}

		if !resp.Success() {
			util.Log().Error("信任设备失败: %s", resp.ErrorMessage)
			return
		}

		util.Log().Info("信任设备成功")
	},
}

func init() {
	rootCmd.AddCommand(trustCmd)
}
