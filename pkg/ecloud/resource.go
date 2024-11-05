package ecloud

import (
	"net/http"
)

type ComputerOperation string
type ResourceStatus string

const (
	ComputerOperationRestart   ComputerOperation = "restart"   // 重启
	ComputerOperationShutdown  ComputerOperation = "shutdown"  // 关机
	ComputerOperationReload    ComputerOperation = "reload"    // 重装
	ComputerOperationAvailable ComputerOperation = "available" // 开机

	ResourceStatusAvailable   ResourceStatus = "available"   // 运行中
	ResourceStatusOnAvailable ResourceStatus = "onAvailable" // 开机中
	ResourceStatusShutdown    ResourceStatus = "shutdown"    // 已关机
	ResourceStatusOnShutdown  ResourceStatus = "onShutdown"  // 关机中
	ResourceStatusReload      ResourceStatus = "onReload"    // 重装中
	ResourceStatusRestart     ResourceStatus = "onRestart"   // 重启中
)

func GetComputerStatus(status string) ResourceStatus {
	switch ResourceStatus(status) {
	case ResourceStatusAvailable, ResourceStatusOnAvailable,
		ResourceStatusShutdown, ResourceStatusOnShutdown,
		ResourceStatusReload, ResourceStatusRestart:
		return ResourceStatus(status)
	default:
		return ""
	}
}

type ComputerInfo struct {
	MachineID   string         `json:"machineId"`
	MachineName string         `json:"machineName"`
	CompanyCode string         `json:"companyCode"`
	Status      ResourceStatus `json:"resourceStatus"`
}

// OperateComputer 云电脑操作
func (c *Client) OperateComputer(info ComputerInfo, operate ComputerOperation) (*Response, error) {
	request, err := c.request(http.MethodPost, "/resource/operate", map[string]any{
		"machineId":   info.MachineID,
		"machineName": info.MachineName,
		"operate":     operate,
		"accessToken": c.session.Token,
	})
	if err != nil {
		return nil, err
	}

	return request, nil
}
