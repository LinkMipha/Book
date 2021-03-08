package gin_server

import (
	"Book/basic"
	"Book/data"
	"Book/model"
	"github.com/gin-gonic/gin"
	"log"
)

//书籍操作接口

type GetBookByIdRsp struct {
	Rsp struct{
		Status string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			Name string `json:"name"`
			Id   int    `json:"id"`
		} `json:"data"`
	}
}

type GetBookByIdReq struct {
	Id int64 `json:"id" form:"id"`
}

func GetBookById(c *gin.Context)  {
	req:=GetBookByIdReq{}
	rsp:=GetBookByIdRsp{}

	if err:=c.Bind(req);err!=nil{
		log.Printf("GetBookById param err: %v",err)
		basic.ResponseError(c,rsp,"param error")
		return
	}
	book,err:= new(model.Book).GetBook(data.Db,req.Id)
	if err!=nil{
		log.Printf("GetBookById GetBook error :%v",err)
		basic.ResponseError(c,rsp,"GetBookById GetBook error")
		return
	}
	rsp.Rsp.Data.Name = book.BookName
	basic.ResponseOk(c,rsp,"")
}