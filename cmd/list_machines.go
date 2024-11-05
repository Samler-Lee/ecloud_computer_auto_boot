package cmd

import (
	"ecloud_computer_auto_boot/pkg/conf"
	"ecloud_computer_auto_boot/pkg/ecloud"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.ecloud.com/ecloud/ecloudsdkcomputer"
	"gitlab.ecloud.com/ecloud/ecloudsdkcomputer/model"
	"gitlab.ecloud.com/ecloud/ecloudsdkcore/config"
)

// listMachinesCmd represents the listMachine command
var listMachinesCmd = &cobra.Command{
	Use:   "list-machines",
	Short: "List all machines in your account",
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init()
		if conf.Secret.Type == "public" {
			client, err := ecloud.NewClient(conf.Secret.Username, conf.Secret.Password)
			if err != nil {
				fmt.Printf("客户端初始化失败: %s\n", err)
				return
			}

			_, err = client.Login()
			if err != nil {
				fmt.Printf("登录错误: %s\n", err)
				return
			}

			if !client.HasTrustDeviceRecord() {
				fmt.Printf("你还没有进行设备信任，请先执行 trust 参数命令\n")
			}

			_, err = client.VerifyAccessTicket()
			if err != nil {
				fmt.Printf("验证登录凭证失败: %s\n", err)
				return
			}

			resp, err := client.GetDeviceInfo()
			if err != nil {
				fmt.Printf("获取资源列表时错误: %s\n", err)
				return
			}

			if !resp.Success() {
				fmt.Printf("获取资源列表时错误: %s\n", resp.ErrorMessage)
				return
			}

			fmt.Printf("您的账户下有以下资源:\n")
			list := resp.GetBody()["machineList"].([]any)
			for i, machine := range list {
				info := machine.(map[string]any)
				fmt.Printf("[%d] MachineID: %s, MachineName: %s, MachineStatus: %s (%s)\n", i, info["machineId"], info["machineName"], info["machineStatus"], info["resourceStatusCn"].(string))
			}
		} else {
			machines, err := findAllMachine()
			if err != nil {
				fmt.Printf("获取资源列表时错误: %s\n", err)
				return
			}

			fmt.Printf("您的账户下有以下资源:\n")
			for i, machine := range machines {
				fmt.Printf("[%d] MachineID: %s, MachineName: %s, MachineStatus: %s\n", i, *machine.MachineId, *machine.MachineName, *machine.MachineStatus)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listMachinesCmd)
}

func findAllMachine() ([]model.GetResourceListResponseData, error) {
	var result []model.GetResourceListResponseData
	client := ecloudsdkcomputer.NewClient(&config.Config{
		AccessKey: &conf.Secret.AccessKey,
		SecretKey: &conf.Secret.SecretKey,
		PoolId:    &conf.Secret.PoolId,
	})

	maxPage := 1
	for i := 0; i < maxPage; i++ {
		body := &model.GetResourceListBody{}
		body.SetPage(int32(i + 1)).SetPageSize(50)

		request := &model.GetResourceListRequest{
			GetResourceListBody: body,
		}

		resp, err := client.GetResourceList(request)
		if err != nil {
			return nil, err
		}

		maxPage = int(*resp.Body.TotalSize / 50)
		for _, data := range *resp.Body.Data {
			result = append(result, data)
		}
	}

	return result, nil
}
