package ecloud

import "encoding/json"

type DeviceInfo struct {
	ClientType         string  `json:"clientType" mapstructure:"clientType"`
	DeviceUid          string  `json:"deviceUid" mapstructure:"deviceUid"`
	DeviceName         string  `json:"deviceName" mapstructure:"deviceName"`
	CompanyCode        string  `json:"companyCode" mapstructure:"companyCode"`
	ClientVersion      string  `json:"clientVersion" mapstructure:"clientVersion"`
	DeviceType         string  `json:"deviceType" mapstructure:"deviceType"`
	OperatingVersion   string  `json:"operatingVersion" mapstructure:"operatingVersion"`
	OperatingSystem    string  `json:"operatingSystem" mapstructure:"operatingSystem"`
	DeviceSystem       string  `json:"deviceSystem" mapstructure:"deviceSystem"`
	DiskTotal          string  `json:"diskTotal" mapstructure:"diskTotal"`
	DiskUsed           string  `json:"diskUsed" mapstructure:"diskUsed"`
	Ram                float64 `json:"ram" mapstructure:"ram"`
	Cores              int     `json:"cores" mapstructure:"cores"`
	IpAddress          string  `json:"ipAddress" mapstructure:"ipAddress"`
	MacAddress         string  `json:"macAddress" mapstructure:"macAddress"`
	DeviceCompany      string  `json:"deviceCompany" mapstructure:"deviceCompany"`
	DeviceModel        string  `json:"deviceModel" mapstructure:"deviceModel"`
	Processor          string  `json:"processor" mapstructure:"processor"`
	SystemArchitecture string  `json:"systemArchitecture" mapstructure:"systemArchitecture"`
}

func getDeviceInfo() map[string]any {
	device := &DeviceInfo{
		ClientType:         "pc_windows",
		DeviceUid:          "Default string",
		DeviceName:         "VM-12",
		CompanyCode:        "ECloud",
		ClientVersion:      getClientVersion(),
		DeviceType:         "pc",
		OperatingVersion:   "10.0.22631",
		OperatingSystem:    "10.0.22631",
		DeviceSystem:       "Windows 11",
		DiskTotal:          "6676",
		DiskUsed:           "3061",
		Ram:                63.83799362182617,
		Cores:              28,
		IpAddress:          "10.0.0.2",
		MacAddress:         "00:00:00:00:00:00",
		DeviceCompany:      "Micro-Star International Co., Ltd.",
		DeviceModel:        "MS-7E07",
		Processor:          "Intel(R) Core(TM) i7-14700KF",
		SystemArchitecture: "ia32",
	}

	bytes, _ := json.Marshal(device)
	var result map[string]any
	_ = json.Unmarshal(bytes, &result)

	return result
}
