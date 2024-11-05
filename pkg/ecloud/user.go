package ecloud

import (
	"net/http"
)

// GetDeviceInfo 获取云电脑信息
func (c *Client) GetDeviceInfo() (*Response, error) {
	return c.request(http.MethodPost, "/user/getDeviceInfo", map[string]any{
		"accessToken": c.session.Token,
		"companyCode": "H3C",
		"allCompany":  true,
	})
}
