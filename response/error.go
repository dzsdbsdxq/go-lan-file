package response

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FileError struct {
	Name   string `xml:"name"`
	Type   string `xml:"type"`
	Sign   string `xml:"sign"`
	Reason string `xml:"reason"`
}

func (fe *FileError) HandleFileError(c *gin.Context) {
	xmlData, err := xml.MarshalIndent(fe, "", "  ")
	if err != nil {
		http.Error(c.Writer, "Error encoding XML", http.StatusInternalServerError)
		return
	}
	// 设置Content-Type为XML
	c.Writer.Header().Set("Content-Type", "application/xml")

	// 返回XML数据
	_, err = c.Writer.Write(xmlData)
	if err != nil {
		http.Error(c.Writer, "Error encoding XML", http.StatusInternalServerError)
		return
	}
}
