package basic

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ResponseOK = "OK"
	ResponseErr = "Error"
)

type HttpResponse struct {
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

//返回值进行封装
func ResponseOk(c *gin.Context, data interface{}, desc string) {
	c.JSON(http.StatusOK, HttpResponse{
		Status:      ResponseOK,
		Description: desc,
		Data:        data,
	})
}

func ResponseError(c *gin.Context, data interface{}, desc string) {
	c.JSON(http.StatusOK, HttpResponse{
		Status:      ResponseErr,
		Description: desc,
		Data:        data,
	})
}
