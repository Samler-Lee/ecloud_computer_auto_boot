package ecloud

import "fmt"

func getClientVersion() string {
	return "3.1.2"
}

func getUserAgent() string {
	return fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Ecloud-Cloud-Computer-Application/%s Chrome/108.0.5359.215 Electron/22.3.27 Safari/537.36", getClientVersion())
}
