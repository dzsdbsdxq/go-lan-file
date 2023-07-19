package response

type FileUploadResponse struct {
	HttpCode        int    `json:"http_code"`
	IsUploaded      int    `json:"isUploaded"`      // 是否已经完成上传了 0:否  1:是 - 秒传
	HasBeenUploaded string `json:"hasBeenUploaded"` // 曾经上传过的分片chunkNumber - 断点续传
	Merge           int    `json:"merge"`           // 是否可以合并了   0：否  1:是
	Status          int    `json:"status"`          // 是否可以合并了   0：成功  1:失败
	Msg             string `json:"msg"`             // 其他信息
	URL             string `json:"url"`             // 后端保存的url
}
