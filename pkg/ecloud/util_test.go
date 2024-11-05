package ecloud

import (
	"github.com/iancoleman/orderedmap"
	"net/http"
	"testing"
)

func TestGetTimestamp(t *testing.T) {
	timestamp := getTimestamp()
	t.Log(timestamp)
}

func TestGetUUID(t *testing.T) {
	uuid := getUUID()
	t.Log(uuid)
}

func TestGetSignature(t *testing.T) {
	params := orderedmap.New()
	params.Set("AccessKey", "53bb79015a3f47c4be166d9371f68f14")
	params.Set("SignatureMethod", "HmacSHA1")
	params.Set("SignatureNonce", getUUID())
	params.Set("SignatureVersion", "V2.0")
	params.Set("Timestamp", getTimestamp())
	params.Set("Signature", getSignature(http.MethodPost, "/login/verify", params))

	t.Log(params.Values())
}

func TestDecrypt(t *testing.T) {
	encrypted := ""
	data, err := decryptData(encrypted, getPrivateKey())
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(data))
}
