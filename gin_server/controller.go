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



type CheckOutRsp struct {
	Status string `json:"status"`
	Description string `json:"description"`
	Id string `json:"_id"`
	UpdateTime  string `json:"updatedAt"`
	CreatedAt  string `json:"createdAt"`
	UserName string `json:"userName"`
	Name string `json:"name"`
	IsAdmin bool `json:"isAdmin"`
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

func GetMenu(c *gin.Context) {

	c.JSON(200, "")
}


func LoginIn(c *gin.Context)  {
	c.JSON(200,"stset")
}

func CheckOut(c *gin.Context) {
	rsp := CheckOutRsp{}
	rsp.IsAdmin = true
	rsp.Name = "管理员"
	rsp.CreatedAt = "2021-04-07T02:51:16.008Z"
	rsp.UpdateTime = "2021-04-07T02:51:16.008Z"
	rsp.Id = "606d1e24eebc850a6c30d1c1"
	c.JSON(200, rsp)
}
