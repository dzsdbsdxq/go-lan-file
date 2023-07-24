package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

// RSAReadKeyFromFile 从文件中读取RSA key
func RSAReadKeyFromFile(filename string) []byte {
	f, err := os.Open(filename)
	var b []byte

	if err != nil {
		return b
	}
	defer f.Close()
	fileInfo, _ := f.Stat()
	b = make([]byte, fileInfo.Size())
	f.Read(b)
	return b
}

// RSAEncrypt RSA加密
func RSAEncrypt(data, publicBytes []byte) ([]byte, error) {
	var res []byte
	// 解析公钥
	block, _ := pem.Decode(publicBytes)

	if block == nil {
		return nil, fmt.Errorf("无法加密, 公钥可能不正确")
	}
	//http://localhost:10000/api/file/DjMG

	// 使用X509将解码之后的数据 解析出来
	// x509.MarshalPKCS1PublicKey(block):解析之后无法用，所以采用以下方法：ParsePKIXPublicKey
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("无法加密, 公钥可能不正确, %v", err)
	}

	// 使用公钥加密数据
	pubKey := keyInit.(*rsa.PublicKey)

	keySize := pubKey.Size()

	for i := 0; i < len(data); i += keySize - 11 {
		end := i + keySize - 11
		if end > len(data) {
			end = len(data)
		}
		part := data[i:end]
		chunk, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, part)
		if err != nil {
			return nil, err
		}
		res = append(res, chunk...)
	}
	// 将数据加密为base64格式
	return []byte(EncodeStr2Base64(string(res))), nil
}

// RSADecrypt 对数据进行解密操作
func RSADecrypt(base64Data string, privateBytes []byte) ([]byte, error) {
	var res []byte
	// 将base64数据解析
	data := []byte(DecodeStrFromBase64(base64Data))
	// 解析私钥
	block, _ := pem.Decode(privateBytes)
	if block == nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确")
	}
	// 还原数据
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确, %v", err)
	}
	keySize := privateKey.PublicKey.Size()

	for i := 0; i < len(data); i += keySize {
		end := i + keySize
		if end > len(data) {
			end = len(data)
		}
		part := data[i:end]
		chunk, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, part)
		if err != nil {
			return []byte{}, err
		}
		res = append(res, chunk...)
	}

	return res, nil
}

// EncodeStr2Base64 加密base64字符串
func EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// DecodeStrFromBase64 解密base64字符串
func DecodeStrFromBase64(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}
