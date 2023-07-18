package model

import "time"

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
