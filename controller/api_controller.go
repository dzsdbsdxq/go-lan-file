package controller

import "github.com/gin-gonic/gin"

type IApiController interface {
	NewConnect(c *gin.Context)
}
type ApiController struct {
}

func (api *ApiController) NewConnect(c *gin.Context) {

}

func NewApiController() IApiController {
	return &ApiController{}
}
