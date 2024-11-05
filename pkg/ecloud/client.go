package ecloud

import (
	"bytes"
	"ecloud_computer_auto_boot/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Session struct {
	Ticket                string `json:"ticket"`
	TicketExpired         bool   `json:"ticketExpired"`
	Token                 string `json:"token"`
	TokenExpired          bool   `json:"tokenExpired"`
	LoginUID              string `json:"loginUid"`
	UpdateSessionInterval int64  `json:"updateSessionInterval"`
	UpdateSessionTicker   *time.Ticker
	Mobile                string        `json:"mobile"`
	Username              string        `json:"username"`
	ConnectList           []ConnectItem `json:"connectList"`
}

type Client struct {
	Username            string
	Password            string
	version             string
	userAgent           string
	device              map[string]any
	client              *http.Client
	publicKey           string
	privateKey          string
	isTrustDevice       bool
	session             *Session
	enableAutoReLogin   bool
	reLoginMutex        sync.Mutex
	reLoginWg           sync.WaitGroup
	reVerifyTicketMutex sync.Mutex
	reVerifyTicketWg    sync.WaitGroup
}

func NewClient(username string, password string) (*Client, error) {
	c := &Client{
		Username:  username,
		Password:  password,
		version:   getClientVersion(),
		userAgent: getUserAgent(),
		device:    getDeviceInfo(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		publicKey:           getPublicKey(),
		privateKey:          getPrivateKey(),
		isTrustDevice:       false,
		session:             &Session{},
		enableAutoReLogin:   false,
		reLoginMutex:        sync.Mutex{},
		reLoginWg:           sync.WaitGroup{},
		reVerifyTicketMutex: sync.Mutex{},
		reVerifyTicketWg:    sync.WaitGroup{},
	}

	return c, nil
}

func (c *Client) GetSession() Session {
	return *c.session
}

func (c *Client) request(method string, cmd string, data map[string]any) (*Response, error) {
	params := orderedmap.New()
	params.Set("AccessKey", "53bb79015a3f47c4be166d9371f68f14")
	params.Set("SignatureMethod", "HmacSHA1")
	params.Set("SignatureNonce", getUUID())
	params.Set("SignatureVersion", "V2.0")
	params.Set("Timestamp", getTimestamp())
	params.Set("Signature", getSignature(method, cmd, params))

	query := url.Values{}
	for _, key := range params.Keys() {
		val, _ := params.Get(key)
		query.Add(key, fmt.Sprintf("%v", val))
	}

	var requestBody io.Reader
	if data != nil {
		for key, val := range c.device {
			data[key] = val
		}

		requestData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		encryptedParams, err := encryptData(requestData, c.publicKey)
		if err != nil {
			return nil, err
		}

		requestBytes, err := json.Marshal(map[string]any{
			"params": encryptedParams,
		})
		if err != nil {
			return nil, err
		}

		requestBody = bytes.NewBuffer(requestBytes)
	}

	request, err := http.NewRequest(method, fmt.Sprintf("%s?%s", "https://ecloud.10086.cn/api/cem/gateway/outer/cem-webapi"+cmd, query.Encode()), requestBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed, status code: %d", resp.StatusCode)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var encryptedResp map[string]any
	err = json.Unmarshal(body, &encryptedResp)
	if err != nil {
		return nil, err
	}

	if encryptedData, exists := encryptedResp["params"]; exists {
		decryptedData, err := decryptData(encryptedData.(string), c.privateKey)
		if err != nil {
			return nil, err
		}

		var result Response
		err = json.Unmarshal(decryptedData, &result)
		if err != nil {
			return nil, err
		}

		util.Log().Debug("response: %+v", result)
		// Token失效检测
		if result.ErrorCode == "401" {
			util.Log().Debug("401: %s, cmd: %s, data: %+v", result.ErrorMessage, cmd, data)
			if result.ErrorMessage == "ticket失效" {
				if c.reLoginMutex.TryLock() {
					c.reLoginWg.Add(1)
					c.session.TicketExpired = true

					c.ReLogin()
					c.reLoginMutex.Unlock()
				} else {
					c.reLoginWg.Wait()
				}

				if c.session.TicketExpired {
					return nil, fmt.Errorf("登录会话失效")
				}

				if _, exists := data["accessTicket"]; exists {
					data["accessTicket"] = c.session.Ticket
				}

				if _, exists := data["accessToken"]; exists {
					data["accessToken"] = c.session.Token
				}

				return c.request(method, cmd, data)
			}

			if result.ErrorMessage == "token失效" && cmd != "/login/verifyAccessTicket" {
				if c.reVerifyTicketMutex.TryLock() {
					c.reVerifyTicketWg.Add(1)
					c.session.TokenExpired = true

					c.UpdateAccessToken()
					c.reVerifyTicketMutex.Unlock()
				} else {
					c.reVerifyTicketWg.Wait()
				}

				if c.session.TokenExpired {
					return nil, fmt.Errorf("尝试更新 AccessToken 失败")
				}

				if _, exists := data["accessTicket"]; exists {
					data["accessTicket"] = c.session.Ticket
				}

				if _, exists := data["accessToken"]; exists {
					data["accessToken"] = c.session.Token
				}

				return c.request(method, cmd, data)
			}
		}

		return &result, nil
	}

	return nil, errors.New("params not found in response")
}

func (c *Client) ReLogin() {
	defer c.reLoginWg.Done()
	if !c.session.TicketExpired {
		return
	}

	for range 3 {
		if _, err := c.Login(); err != nil {
			util.Log().Error("尝试重新登录失败: %s", err)
			time.Sleep(3 * time.Second)
		} else {
			_, _ = c.VerifyAccessTicket()
			_, _ = c.RecordDeviceInfo()
			break
		}
	}
}

func (c *Client) UpdateAccessToken() {
	defer c.reVerifyTicketWg.Done()
	if !c.session.TokenExpired {
		return
	}

	for range 3 {
		if _, err := c.VerifyAccessTicket(); err != nil {
			util.Log().Error("尝试重新获取AccessToken失败: %s", err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
}

// HasTrustDeviceRecord 是否有信任设备记录
func (c *Client) HasTrustDeviceRecord() bool {
	return c.isTrustDevice
}
