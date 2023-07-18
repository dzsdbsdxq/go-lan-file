package repository

import (
	"share.ac.cn/cache"
	"share.ac.cn/model"
	"time"
)

type IFileRepository interface {
	GetFileInfoByShareId(shareId string) (*model.FileOnline, error)
	DeleteFilePath(filePath string) error
}
type FileRepository struct{}

func (f *FileRepository) GetFileInfoByShareId(shareId string) (*model.FileOnline, error) {
	//根据shareId，从redis中获取文件信息
	info, err := cache.GetFileOnlineInfo(shareId)
	if err != nil {
		return nil, err
	}
	//判断文件信息是否过期,过期时间小于当前时间
	if info.ExpireTime.Unix() <= time.Now().Unix() {
		//如果过期了，删除文件和清除redis中的文件信息，返回文件过期提示
	}
	//如果文件没有过期，返回文件信息
	return nil, nil
}

func (f *FileRepository) DeleteFilePath(filePath string) error {

	return nil
}

func NewFileRepository() IFileRepository {
	return &FileRepository{}
}
