package gin_server

import (
	"Book/data"
	"Book/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)


type GetBookListReq struct {
	Isbn string `json:"isbn" form:"isbn"`
	Type  string `json:"type" from:"type"`

	Name  string `json:"query" form:"query"`
	PageNum int `json:"pagenum" form:"pagenum"`
	PageSize int `json:"pagesize" form:"pagesize"`
}


type Books struct {
	Isbn  string `json:"isbn"`
	BookName string `json:"book_name"`
    Author string `json:"author"`
	Publish string `json:"publish"`
	Price  float64 `json:"price"`
	BookType string `json:"book_type"`
	Stock int `json:"stock"`
	ImgUrl string `json:"img_url"`
}

type GetBookListRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		Total int `json:"total"`
		Book []Books `json:"book"`
	} `json:"data"`

}


//查询图书列表
func GetBookList(c*gin.Context){
	req:=GetBookListReq{}
	rsp:=GetBookListRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("GetBookList err"+err.Error())
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}

   var book model.Book
	total,err:=book.GetBookTotal(data.Db)
	if err!=nil{
		fmt.Println("GetBookTotal err"+err.Error())
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}
	rsp.Data.Total  = total

    books,err:=book.GetBookByItem(data.Db,req.Isbn,req.Name,req.Type,req.PageNum,req.PageSize)
    if err!=nil{
		fmt.Println("GetBookByItem err"+err.Error())
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}

	for _,v:=range books{
		rsp.Data.Book= append(rsp.Data.Book,Books{
			Isbn: v.Isbn,
			BookName: v.BookName,
			Author:v.Author,
			Publish:v.Publish,
			Price  :v.Price,
			BookType:v.BookType,
			Stock :v.Stock,
			ImgUrl: v.ImgUrl,
		})
	}
	rsp.Status = "200"
	rsp.Description = "查询成功"
	c.JSON(200,rsp)
	return
}


type AddNewBookReq struct {
	Isbn  string `json:"isbn" form:"Isbn"`
	BookName string `json:"book_name" form:"book_name"`
	Author string `json:"author" form:"author"`
	Publish string `json:"publish" form:"publish"`
	Price  string `json:"price" form:"price"`
	BookType string `json:"book_type" form:"book_type"`
	Stock string `json:"stock" form:"stock"`
}


type AddNewBookRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {

	} `json:"data"`

}

//增加书籍信息
func AddNewBook(c*gin.Context)  {
	req:=AddNewBookReq{}
	rsp:=AddNewBookRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("AddNewBook err"+err.Error())
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}
	var book model.Book
	book.Isbn = req.Isbn
	book.BookName = req.BookName
	book.Author = req.Author
	book.Publish = req.Publish
	book.BookType  = req.BookType
	stock,atoiErr:= strconv.Atoi(req.Stock)
    if atoiErr!=nil{
		fmt.Println(" AddNewBook stock strconv.Atoi error",atoiErr)

	}
	book.Stock = stock
	price,parseErr := strconv.ParseFloat(req.Price, 64)
	if parseErr!=nil{
		fmt.Println(" AddNewBook price strconv.ParseFloat error",parseErr)
	}

	book.Price = price
	err:=book.CreateBook(data.Db)
	if err!=nil{
		fmt.Println("AddNewBook  CreateBook err"+err.Error())
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}

	rsp.Status = "200"
	rsp.Description = "添加图书成功"
	c.JSON(200,rsp)
	return
}


type GetBookByIsbnReq struct {
	Isbn string `json:"isbn" form:"isbn"`
}


type GetBookByIsbnRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		Total int `json:"total"`
		Book Books `json:"book"`
	} `json:"data"`

}

func GetBookByIsbn(c*gin.Context)  {
	req:=GetBookByIsbnReq{}
	rsp:=GetBookByIsbnRsp{}
	req.Isbn = c.Param("isbn")
	if req.Isbn==""{
		fmt.Println("GetBookByIsbn  Isbn err ")
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}
	var book model.Book
	bookData,err:=book.GetBookByIsbn(data.Db,req.Isbn)
	if err!=nil{
		fmt.Println("GetBookByIsbn Isbn model  err ")
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "200"
	rsp.Description = "获取成功"
	rsp.Data.Book = Books{
		Isbn :bookData.Isbn,
		BookName:bookData.BookName,
		Author :bookData.Author,
		Publish :bookData.Publish,
		Price :bookData.Price,
		BookType :bookData.BookType,
		Stock :bookData.Stock,
	}
	c.JSON(200,rsp)
	return
}


type EditBookByIsbnReq struct {
	Isbn  string `json:"isbn" form:"Isbn"`
	BookName string `json:"book_name" form:"book_name"`
	Author string `json:"author" form:"author"`
	Publish string `json:"publish" form:"publish"`
	Price  float64 `json:"price" form:"price"`
	BookType string `json:"book_type" form:"book_type"`
	Stock int `json:"stock" form:"stock"`
}


type EditBookByIsbnRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {

	} `json:"data"`

}

//更新
func EditBookByIsbn(c*gin.Context)  {
	req:=EditBookByIsbnReq{}
	rsp:=EditBookByIsbnRsp{}
	req.Isbn = c.Param("isbn")

	if err:=c.Bind(&req);err!=nil{
		fmt.Println("EditBookByIsbn err"+err.Error())
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}

	var book model.Book
	var bookData = make(map[string]interface{},0)
	bookData["book_name"] = req.BookName
	bookData["author"] = req.Author
	bookData["publish"]  =req.Publish
	bookData["book_type"] = req.BookType
	bookData["price"] = req.Price
	bookData["stock"] = req.Stock

   err:=book.UpdateBookByIsbn(data.Db,req.Isbn,bookData)
   if err!=nil{
	   fmt.Println("EditBookByIsbn UpdateBookByIsbn err"+err.Error())
	   rsp.Status = "400"
	   c.JSON(200,rsp)
	   return
   }

   rsp.Status = "200"
   rsp.Description = "修改成功"
   c.JSON(200,rsp)
   return

}


//删除图书信息
type DeleteUserByIsbnReq struct {
	Isbn string `json:"isbn" form:"isbn"`
}


type DeleteUserByIsbnRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {

	} `json:"data"`

}
func DeleteUserByIsbn(c*gin.Context)  {
	req:=DeleteUserByIsbnReq{}
	rsp:=DeleteUserByIsbnRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("DeleteUserByIsbn bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}
	var book model.Book
	err:=book.DeleteBookByIsbn(data.Db,req.Isbn)
	if err!=nil{
		fmt.Println("DeleteBookByIsbn error")
		rsp.Status = "401"
		rsp.Description = "DeleteBookByIsbn error"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "200"
	rsp.Description = "删除成功"
	c.JSON(200,rsp)
	return
}