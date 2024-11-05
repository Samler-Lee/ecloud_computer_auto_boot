package ecloud

import (
	"ecloud_computer_auto_boot/pkg/util"
	"net/http"
)

type ConnectItem struct {
	ConnectID     string `json:"connectId"`
	ConnectStatus bool   `json:"connectStatus"`
	MachineID     string `json:"machineId"`
	CompanyCode   string `json:"companyCode"`
}

// MachineConnect 连接云电脑
func (c *Client) MachineConnect(computer ComputerInfo) (*Response, error) {
	return c.request(http.MethodPost, "/session/machineConnect", map[string]any{
		"ticket":      c.session.Ticket,
		"accessToken": c.session.Token,
		"machineId":   computer.MachineID,
		"machineName": computer.MachineName,
		"status":      "success",
		"flag":        true,
	})
}

// UpdateSession 更新会话状态
func (c *Client) UpdateSession() (*Response, error) {
	params := map[string]any{
		"loginUid":    c.session.LoginUID,
		"loginStatus": "0",
	}

	if c.session.ConnectList != nil {
		params["connectList"] = c.session.ConnectList
	}

	util.Log().Debug("UpdateSession params: %+v", params)
	return c.request(http.MethodPost, "/session/updateSessionStatus", params)
}

// AutoUpdateSession 自动更新会话状态
func (c *Client) AutoUpdateSession() {
	_, err := c.UpdateSession()
	if err != nil {
		util.Log().Error("更新会话失败：%s", err)
	}

	for {
		select {
		case <-c.session.UpdateSessionTicker.C:
			_, err := c.UpdateSession()
			if err != nil {
				util.Log().Error("更新会话失败：%s", err)
			}
		}
	}
}
