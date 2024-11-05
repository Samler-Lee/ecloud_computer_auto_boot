package ecloud

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func getPublicKey() string {
	return "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCqisJL7YvdPC/gJA7fLrr1G+t6J0arJr0sVfieVJTXTclm/2afP/fjNYY/CFcg1MUx8KPmPC2CqsUHRMZq6Ev1/UNXE74I1TfJC/2b8aexcdZ+Lokj7AwzrM9yPy2qfV6vXtxyRrTs+JcFHVXtV6phNkorNyIahyfy46+iNB+FSQIDAQAB\n-----END PUBLIC KEY-----"
}

func getPrivateKey() string {
	return "-----BEGIN PRIVATE KEY-----\nMIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAKqKwkvti908L+AkDt8uuvUb63onRqsmvSxV+J5UlNdNyWb/Zp8/9+M1hj8IVyDUxTHwo+Y8LYKqxQdExmroS/X9Q1cTvgjVN8kL/Zvxp7Fx1n4uiSPsDDOsz3I/Lap9Xq9e3HJGtOz4lwUdVe1XqmE2Sis3IhqHJ/Ljr6I0H4VJAgMBAAECgYBD6lx0BlajtRtPxKxTfvWfNQ4yqD+BWz0M0fPfgcmAcI7bQKyqkLv0NNWQdo7UGUeqmq16u85X8g/i1CW8X2QYHOSYNBUWsK3k5gFT1wdk+bwuIMZqgjEc48TXzM4pidcplJLyD1tnNiubzcXIsZCIIuQ/GmWcuxn7ULHnXDsQMQJBANMl4V97be6fkd1beGqYZWIx3XNnL96AQsapBrEbbORTu/JnwTCRbsRWRBHU11FZuK85dBDXrH8reoAsgepmsF0CQQDOxL99OFjozj8g1weFGwI/otMKcPhkaslU2tj3QF44zT1TZiOZ710I8GQLPlKeu1yGWvVUwgH4bCY0M8M1/gndAkB9sU4RTeOqKjllwT7UjbXEl5SRTzrSxB18L0B5i67t2N7INXVumRSMMiJBTyeCGNv1C0mJgSoBZft9c4E+7TRNAkB+7Azza7Q/6+KaYQRPs32U3HkZbrE6ysYdXV1ToOJ1kZ60Y/00j9cXFqECudXzc+Ve39S6m4CkIpbs8l1A9ljNAkBy6Rp19R5wWMr/3feIMZ18akWXT5mgRvZpkT5MgmrjVu1lRv8bHsEsAzRYvdPSjzp0nCkUbOWUITxWp7d//Fwc\n-----END PRIVATE KEY-----"
}

func encryptData(data []byte, key string) (string, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not a valid RSA public key")
	}

	// 分块加密（公钥模长 - 11）
	chunkSize := pub.N.BitLen()/8 - 11
	encryptedData := make([]byte, 0)
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]
		encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, pub, chunk)
		if err != nil {
			return "", err
		}

		encryptedData = append(encryptedData, encryptedChunk...)
	}

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func decryptData(data string, key string) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	var privateKey *rsa.PrivateKey
	var err error
	parsedKey, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
	if parseErr == nil {
		var ok bool
		privateKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not a valid RSA private key")
		}
	} else {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	decryptedData := make([]byte, 0)
	chunkSize := 128
	for i := 0; i < len(encryptedBytes); i += chunkSize {
		end := i + chunkSize
		if end > len(encryptedBytes) {
			end = len(encryptedBytes)
		}

		chunk := encryptedBytes[i:end]
		decryptedChunk, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, chunk)
		if err != nil {
			return nil, err
		}

		decryptedData = append(decryptedData, decryptedChunk...)
	}

	return decryptedData, nil
}
