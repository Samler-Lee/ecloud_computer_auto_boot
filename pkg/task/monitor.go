package task

import (
	"ecloud_computer_auto_boot/pkg/conf"
	"ecloud_computer_auto_boot/pkg/util"
	"gitlab.ecloud.com/ecloud/ecloudsdkcomputer/model"
)

func machineMonitor() {
	page := 1
	failedCnt := 0

	for {
		if failedCnt >= 3 {
			return
		}

		response, err := findMachine(int32(page))
		if err != nil {
			failedCnt++
			util.Log().Error("[实例状态监控] 请求失败: %s", err)
			continue
		}

		if *response.ErrorCode != "" {
			failedCnt++
			util.Log().Error("[实例状态监控] 请求失败: %s", *response.ErrorMessage)
			continue
		}

		respData := *response.Body.Data
		for _, instance := range respData {
			machineId := *instance.MachineId
			if util.InArray(conf.Cron.MachineList, machineId) {
				if *instance.MachineStatus == "shutdown" {
					go func() {
						util.Log().Info("[实例状态监控] 检测到实例 %s 已关机, 请求启动", machineId)
						resp, err := startupMachine(machineId)
						if err != nil {
							util.Log().Error("[启动实例] 实例 %s 启动失败: %s", err)
							return
						}

						if *resp.ErrorCode != "" {
							util.Log().Error("[启动实例] 实例 %s 启动失败: %s", *response.ErrorMessage)
							return
						}
					}()
				}
			}
		}

		// 页尾判定
		if *response.Body.TotalSize <= int32(page*50) {
			return
		}

		page++
	}
}

func findMachine(page int32) (*model.GetResourceListResponse, error) {
	body := &model.GetResourceListBody{}
	body.SetPage(page).SetPageSize(50)

	request := &model.GetResourceListRequest{
		GetResourceListBody: body,
	}

	return cloudClient.GetResourceList(request)
}

func startupMachine(machineId string) (*model.OperateMachineByAvailableResponse, error) {
	request := &model.OperateMachineByAvailableRequest{}
	request.OperateMachineByAvailableQuery.SetMachineId(machineId)

	return cloudClient.OperateMachineByAvailable(request)
}
