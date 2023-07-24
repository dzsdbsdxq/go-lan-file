package common

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Resp struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetUserIdByRandom() string {
	userId := uuid.Must(uuid.NewV4(), nil).String()
	h := md5.New()
	h.Write([]byte(userId + GetRandomId(10)))
	return hex.EncodeToString(h.Sum(nil))[0:32]
}

func GetRandomId(length int) (orderId string) {
	// 定义随机字符集
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	// 生成随机ID
	id := make([]byte, length)
	for i := 0; i < length; i++ {
		id[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(id)
}
