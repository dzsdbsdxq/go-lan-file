package request

type FileUploadRequest struct {
	FileName         string `json:"file_name" form:"file_name"`
	ChunkNumber      int    `json:"chunk_number" form:"chunk_number"`
	CurrentChunkSize int    `json:"current_chunk_size" form:"current_chunk_size"`
	FileSize         int    `json:"file_size" form:"file_size"`
	TotalChunks      int    `json:"total_chunks" form:"total_chunks"`
	Identifier       string `json:"identifier" form:"identifier"`
	UploadId         string `json:"upload_id" form:"upload_id"`
	ShareId          string `json:"share_id" form:"share_id"`
}
