package model

import (
	"github.com/jinzhu/gorm"
	"share.ac.cn/common"
	"time"
)

type Files struct {
	common.Model
	Uid        string    `json:"uid" gorm:"index"`
	IsUploaded int       `json:"is_uploaded"`
	IsDel      int       `json:"is_del"`
	Views      int       `json:"views"`
	Downloads  int       `json:"downloads"`
	FileSize   int       `json:"file_size"`
	ExpiredAt  time.Time `json:"expired_at" gorm:"index"`
	FileName   string    `json:"file_name"`
	FileId     string    `json:"file_id"`
	ShareId    string    `json:"share_id" gorm:"index"`
	FilePath   string    `json:"file_path"`
	FileExt    string    `json:"file_ext"`
}

func (files *Files) BeforeCreate(scope *gorm.Scope) error {
	nowTime := time.Now()
	err := scope.SetColumn("CreatedAt", nowTime)
	if err != nil {
		return err
	}
	//设置文件一天的过期时间86400
	err = scope.SetColumn("ExpiredAt", nowTime.Add(24*time.Hour))
	if err != nil {
		return err
	}
	return nil
}

func (files *Files) GetFileByShareId(shareId string) (file *Files) {
	common.GetDb().Where("share_id = ? AND is_del = ?", shareId, 1).First(&files)
	return
}
func (files *Files) GetFileByFileId(fileId string) (*Files, error) {
	//var file *Files
	err := common.GetDb().Where("file_id = ?", fileId).First(&files).Error
	return files, err
}

func (files *Files) AddFiles() (int, error) {
	err := common.GetDb().Create(files).Error
	return files.ID, err
}
func (files *Files) DeleteFiles(id int) bool {
	common.GetDb().Where("id = ?", id).Delete(files)
	return true
}
func (files *Files) DeleteFilesByFileId(fileId string) bool {
	common.GetDb().Where("file_id = ?", fileId).Delete(Files{})
	return true
}
func (files *Files) ExistFileByShareId(shareId string) bool {
	common.GetDb().Select("id").Where("share_id = ?", shareId).First(&files)
	if files.ID > 0 {
		return true
	}
	return false
}

func (files *Files) ExistFileByFileId(fileId string) bool {
	common.GetDb().Select("id").Where("file_id = ?", fileId).First(&files)
	if files.ID > 0 {
		return true
	}
	return false
}

func (files *Files) UpdateViews() {
	common.GetDb().Model(&files).Where("id = ?", files.ID).UpdateColumn("views", gorm.Expr("views + ?", 1))
}
func (files *Files) UpdateDownload() {
	common.GetDb().Model(&files).Where("id = ?", files.ID).UpdateColumn("downloads", gorm.Expr("downloads + ?", 1))
}
func (files *Files) UpdateColumn(key string, value interface{}) error {
	err := common.GetDb().Model(&files).UpdateColumn(key, value).Error
	return err
}

type FileOnline struct {
	CreateTime   time.Time `json:"create_time"`
	ExpireTime   time.Time `json:"expire_time"`
	FileSize     uint64    `json:"file_size"`
	FileViews    uint64    `json:"views"`
	FileDowns    uint64    `json:"downloads"`
	ShareId      string    `json:"sid"`
	FileId       string    `json:"file_id"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	FileExt      string    `json:"file_ext"`
	FileHash     string    `json:"file_hash"`
	FileHashName string    `json:"file_hash_name"`
}
