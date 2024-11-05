package ecloud

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

// Login 登录
func (c *Client) Login() (*Response, error) {
	resp, err := c.request(http.MethodPost, "/login/verify", map[string]any{
		"username": c.Username,
		"password": c.Password,
	})
	if err != nil {
		return nil, err
	}

	if resp.State != "OK" {
		return nil, errors.New(resp.ErrorMessage)
	}

	if resp.ErrorCode != "200" && resp.ErrorCode != "30002009" {
		return nil, errors.New(resp.ErrorMessage)
	}

	if body, ok := resp.Body.(map[string]any); ok {
		c.enableAutoReLogin = true
		c.session.TicketExpired = false
		c.isTrustDevice = resp.ErrorCode != "30002009"

		if val, ok := body["accessTicket"]; ok && val != nil {
			c.session.Ticket = val.(string)
		}

		if val, ok := body["mobile"]; ok && val != nil {
			c.session.Mobile = val.(string)
		}

		if val, ok := body["username"]; ok && val != nil {
			c.session.Username = val.(string)
		}
	}

	return resp, nil
}

// SendTrustDeviceVerifySms 发送信任设备验证短信
func (c *Client) SendTrustDeviceVerifySms() (*Response, error) {
	return c.request(http.MethodPost, "/login/sendVerifySms", map[string]any{
		"mobile":   c.session.Mobile,
		"codeType": "trust",
	})
}

// TrustDevice 信任设备
func (c *Client) TrustDevice(code string) (*Response, error) {
	return c.request(http.MethodPost, "/login/trustDevice", map[string]any{
		"mobile":           c.session.Mobile,
		"verificationCode": code,
	})
}

// VerifyAccessTicket 验证 AccessTicket 并获取 AccessToken
func (c *Client) VerifyAccessTicket() (*Response, error) {
	resp, err := c.request(http.MethodPost, "/login/verifyAccessTicket", map[string]any{
		"accessTicket": c.session.Ticket,
	})
	if err != nil {
		return nil, err
	}

	if resp.State != "OK" {
		return resp, errors.New(resp.ErrorMessage)
	}

	if resp.ErrorCode != "200" {
		return resp, errors.New(resp.ErrorMessage)
	}

	if body, ok := resp.Body.(map[string]any); ok {
		c.session.TokenExpired = false

		if val, ok := body["accessToken"]; ok && val != nil {
			c.session.Token = val.(string)
		}

		if val, ok := body["userName"]; ok && val != nil {
			c.session.Username = val.(string)
		}

		if val, ok := body["mobile"]; ok && val != nil {
			c.session.Mobile = val.(string)
		}
	}

	return resp, nil
}

// RecordDeviceInfo 添加信任设备信息
func (c *Client) RecordDeviceInfo() (*Response, error) {
	resp, err := c.request(http.MethodPost, "/login/recordDeviceInfo", map[string]any{
		"accessToken": c.session.Token,
	})
	if err != nil {
		return nil, err
	}

	if resp.State != "OK" {
		return resp, errors.New(resp.ErrorMessage)
	}

	if resp.ErrorCode != "200" {
		return resp, errors.New(resp.ErrorMessage)
	}

	if body, ok := resp.Body.(map[string]any); ok {
		if val, ok := body["loginUid"]; ok && val != nil {
			c.session.LoginUID = val.(string)
		}

		if val, ok := body["intervalTime"]; ok && val != nil {
			intervalTime, _ := strconv.ParseInt(val.(string), 10, 64)
			if c.session.UpdateSessionTicker == nil {
				c.session.UpdateSessionTicker = time.NewTicker(time.Duration(intervalTime) * time.Second)
				go c.AutoUpdateSession()
			} else if intervalTime != c.session.UpdateSessionInterval {
				c.session.UpdateSessionTicker.Reset(time.Duration(intervalTime) * time.Second)
			}

			c.session.UpdateSessionInterval = intervalTime
		}
	}

	return resp, nil
}
