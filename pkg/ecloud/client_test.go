package ecloud

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient(os.Getenv("PUBLIC_USERNAME"), os.Getenv("PUBLIC_PASSWORD"))
	if err != nil {
		t.Error(err)
		return
	}

	_, err = c.Login()
	if err != nil {
		t.Error(err)
		return
	}

	if !c.isTrustDevice {
		resp, err := c.SendTrustDeviceVerifySms()
		if err != nil {
			t.Error(err)
			return
		}

		t.Logf("已向 %s 发送验证码, 请在 %d 秒内输入", c.session.Mobile, int(resp.Body.(map[string]any)["expireTime"].(float64)))
		var code string
		_, err = fmt.Scanf("%s", &code)
		if err != nil {
			t.Error(err)
			return
		}

		resp, err = c.TrustDevice(code)
		if err != nil {
			t.Error(err)
			return
		}
	}

	resp, err := c.VerifyAccessTicket()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("response: %+v", resp)

	resp, err = c.GetDeviceInfo()
	if err != nil {
		t.Error(err)
		return
	}

	if !resp.Success() {
		t.Error(resp.ErrorMessage)
		return
	}

	body := resp.GetBody()
	machineList := body["machineList"].([]any)
	var computers []ComputerInfo
	for _, machine := range machineList {
		info := machine.(map[string]any)
		computers = append(computers, ComputerInfo{
			MachineID:   info["machineId"].(string),
			MachineName: info["machineName"].(string),
			CompanyCode: info["companyCode"].(string),
			Status:      GetComputerStatus(info["resourceStatus"].(string)),
		})
		t.Logf("id: %s, name: %s, companyCode: %s, status: %s (%s)", info["machineId"], info["machineName"], info["companyCode"], info["resourceStatus"], info["resourceStatusCn"])
	}

	if len(computers) > 0 {
		resp, err := c.OperateComputer(computers[0], ComputerOperationAvailable)
		if err != nil {
			t.Error(err)
			return
		}

		if !resp.Success() {
			t.Error(resp.ErrorMessage)
			return
		}

		t.Logf("response: %+v", resp)
		for i := 0; i < 30; i++ {
			resp, err := c.GetDeviceInfo()
			if err != nil {
				t.Error(err)
				return
			}

			if !resp.Success() {
				t.Error(resp.ErrorMessage)
				return
			}

			body := resp.GetBody()
			machineList := body["machineList"].([]any)
			for _, machine := range machineList {
				info := machine.(map[string]any)
				t.Logf("id: %s, name: %s, companyCode: %s, status: %s (%s)", info["machineId"], info["machineName"], info["companyCode"], info["resourceStatus"], info["resourceStatusCn"])
			}

			time.Sleep(1 * time.Second)
		}
	}
}
