package ecloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func getTimestamp() string {
	return time.Now().Format("2006-1-2T15:04:05Z")
}

func getUUID() string {
	str := "0123456789abcdef"
	data := make([]byte, 32)
	for i := 0; i < 32; i++ {
		data[i] = str[rand.Intn(16)]
	}

	num, _ := strconv.ParseInt(string(data[19]), 16, 0)
	tmpIdx := (3 & num) | 8

	data[14] = '4'
	data[19] = str[tmpIdx]
	data[8] = data[23]
	data[13] = data[23]
	data[18] = data[23]

	return string(data)
}

func getSignature(method string, cmd string, params *orderedmap.OrderedMap) string {
	var paramsArray []string
	for _, key := range params.Keys() {
		val, _ := params.Get(key)
		paramsArray = append(paramsArray, fmt.Sprintf("%s=%s", key, url.QueryEscape(val.(string))))
	}

	hash256 := sha256.New()
	hash256.Write([]byte(strings.Join(paramsArray, "&")))
	msg := fmt.Sprintf("%s\n%s\n%s", strings.ToUpper(method), url.QueryEscape("/api/cem/gateway/outer/cem-webapi"+cmd), hex.EncodeToString(hash256.Sum(nil)))

	hash1 := hmac.New(sha1.New, []byte("BC_SIGNATURE&6b0d3b93f3aa4c7ea076c841bead1ddd"))
	hash1.Write([]byte(msg))
	return hex.EncodeToString(hash1.Sum(nil))
}
