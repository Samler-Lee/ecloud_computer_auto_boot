package task

import (
	"ecloud_computer_auto_boot/pkg/conf"
	"ecloud_computer_auto_boot/pkg/ecloud"
	"ecloud_computer_auto_boot/pkg/util"
)

func startMachineMonitorOnPublic() {
	resp, err := publicClient.GetDeviceInfo()
	if err != nil {
		util.Log().Error("获取设备信息失败: %s", err)
		return
	}

	if !resp.Success() {
		util.Log().Error("获取云电脑列表失败: %s", resp.ErrorMessage)
		return
	}

	monitorAll := len(conf.Cron.Machines) == 0

	body := resp.GetBody()
	machineList := body["machineList"].([]any)
	for _, machine := range machineList {
		info := machine.(map[string]any)

		computer := ecloud.ComputerInfo{
			MachineID:   info["machineId"].(string),
			MachineName: info["machineName"].(string),
			CompanyCode: info["companyCode"].(string),
			Status:      ecloud.GetComputerStatus(info["resourceStatus"].(string)),
		}

		util.Log().Debug("id: %s, name: %s, companyCode: %s, status: %s", computer.MachineID, computer.MachineName, computer.CompanyCode, computer.Status)
		if (monitorAll || util.InArray(conf.Cron.Machines, computer.MachineID)) && computer.Status == ecloud.ResourceStatusShutdown {
			go func() {
				util.Log().Info("[%s] 检测到机器已关机, 请求开机", computer.MachineID)
				_, err := publicClient.OperateComputer(computer, ecloud.ComputerOperationAvailable)
				if err != nil {
					util.Log().Error("[%s] 开机失败: %s", computer.MachineID, err)
				} else {
					util.Log().Info("[%s] 已完成开机操作", computer.MachineID)
				}
			}()
		}
	}
}
