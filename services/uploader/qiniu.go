package uploader

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"share.ac.cn/common"
	"share.ac.cn/config"
	"time"
)

var qBoxMac *qbox.Mac

func InitQiNiuMac() {
	qBoxMac = qbox.NewMac(config.Conf.QinNiu.AccessKey, config.Conf.QinNiu.SecretKey)
}

func GetQiNiuAccessToken() string {
	putPolicy := storage.PutPolicy{
		Scope:            config.Conf.QinNiu.Bucket,
		CallbackURL:      config.Conf.QinNiu.CallbackUrl,
		CallbackBody:     `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)","uid":"$(x:uid)"}`,
		CallbackBodyType: "application/json",
	}
	putPolicy.Expires = 600 //10分钟有效期
	upToken := putPolicy.UploadToken(qBoxMac)
	return upToken
}
func GetQiNiuDownloadUrl(key string) string {
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(qBoxMac, config.Conf.QinNiu.DownUrl, key, deadline)
	return privateAccessURL
}
func bucketManager() *storage.BucketManager {
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Region=&storage.ZoneHuabei
	return storage.NewBucketManager(qBoxMac, &cfg)
}

func DeleteObject(key string) error {
	err := bucketManager().Delete(config.Conf.QinNiu.Bucket, key)
	if err != nil {
		common.Log.Errorf("删除七牛云文件%s出错，错误原因：%s\n", key, err.Error())
		return err
	}
	return nil
}

func DeleteAfterDays(key string) error {
	err := bucketManager().DeleteAfterDays(config.Conf.QinNiu.Bucket, key, 1)
	if err != nil {
		common.Log.Errorf("设置七牛云文件%s过期时间出错，错误原因：%s\n", key, err.Error())
		return err
	}
	return nil
}
