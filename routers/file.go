package routers

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/controller"
)

func InitFileRouters(r *gin.RouterGroup) gin.IRoutes {
	fileController := controller.NewFileController()
	router := r.Group("/file")
	{
		router.Any("/upload", fileController.UploadFile)
		router.GET("/download/:shareId", fileController.DownloadFile)
		//router.PATCH("/update/:addressId", addressController.UpdateAddressById)
		//router.DELETE("/delete/batch", addressController.BatchDeleteAddressByIds)
		//router.GET("/getCredential", addressController.GetCredential)
	}
	return r
}
