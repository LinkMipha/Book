package gin_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type HelloReq struct {
	Id int `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

type HelloRsp struct {
	Rsp struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			Name string `json:"name"`
			Id   int    `json:"id"`
		} `json:"data"`
	}
}

func GetTest(c*gin.Context)  {
	req:= HelloReq{}
	rsp:= HelloRsp{}
	if err:=c.Bind(&req);err!=nil{
		log.Fatal("test") //服务断开
	}
	fmt.Println(req)
	rsp.Rsp.Data.Name = req.Name
	rsp.Rsp.Data.Id = req.Id
	c.JSON(200,rsp)
}

func GetMenu(c*gin.Context)  {

	c.JSON(200,"")
}

func LoginIn(c *gin.Context)  {
	c.JSON(200,"")
}


func CheckOut(c *gin.Context)  {
	c.JSON(200,"")
}
