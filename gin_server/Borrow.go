package gin_server

import (
	"Book/data"
	"Book/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)


type ReNewBorrowReq struct {
	UserId string `json:"user_id"`
	BookId string `json:"book_id"`
}

type ReNewBorrowRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {

	} `json:"data"`

}



//续借
func ReNewBorrow(c*gin.Context)  {
	req:=ReNewBorrowReq{}
	rsp:=ReNewBorrowRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("ReNewBorrow bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}

	var bor model.Borrow
	//软删除
	err := bor.ReNewBorrow(data.Db,req.UserId,req.BookId)
	if err!=nil{
		fmt.Println("ReNewBorrow model ReNewBorrow err"+err.Error())
		rsp.Description = err.Error()
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "200"
	rsp.Description = "success"
	c.JSON(200,rsp)
	return
}



type AddBorrowReq struct {
     Borrow model.Borrow `json:"borrow" form:"borrow"`
}


type AddBorrowRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {

	} `json:"data"`

}

//借书增加
func AddBorrow(c*gin.Context) {
	req:=AddBorrowReq{}
	rsp:=AddBorrowRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("AddBorrow bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}
	//判断是否借阅已经存在，且未删除未归还

	var borrow model.Borrow
	bors,err:=borrow.GetBorrowByUserIdBookId(data.Db,req.Borrow.UserId,req.Borrow.BookId)
	if err!=nil{
		fmt.Println("AddBorrow model GetBorrowByUserIdBookId err"+err.Error())
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}
	if len(bors)>0{
		fmt.Println("AddBorrow len bors >0 ")
		rsp.Status = "401"
		rsp.Description = "借阅已存在，请归还后再借"
		c.JSON(200,rsp)
		return
	}

	req.Borrow.BrrowTime = time.Now()
	//归还时间增加一个月
	req.Borrow.RetTime = time.Now().AddDate(0,1,0)
	req.Borrow.RealTime = time.Now().AddDate(0,0,-1)
	//使用事务会出释放问题
	//data.Db.Begin()

	//减少库存数量
	var book model.Book
	err =book.SubBookStockByIsbn(data.Db,req.Borrow.BookId)
	if err!=nil{
		//data.Db.Rollback()
		fmt.Println("AddBorrow model SubBookStockByIsbn err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "库存不足"
		c.JSON(200,rsp)
		return
	}


	//借阅图书
	err =borrow.AddBorrow(data.Db,req.Borrow)
	if err!=nil{
		//data.Db.Rollback()
		//data.Db.Rollback()
		fmt.Println("AddBorrow model AddBorrow err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "model AddBorrow error"
		c.JSON(200,rsp)
		return
	}


	//提交事务
	//data.Db.Commit()
	rsp.Status = "200"
	rsp.Description = "借阅成功"
	c.JSON(200,rsp)
	return
}



//借阅管理
type GetBorrowRecordsReq struct {
	Isbn string `json:"isbn" form:"isbn"`
	Name  string `json:"query" form:"query"` //user_id 用户名
	PageNum int `json:"pagenum" form:"pagenum"`
	PageSize int `json:"pagesize" form:"pagesize"`
}

type Borrows struct {
	UserId string `json:"user_id" gorm:"user_id"`
	BookId string `json:"book_id" gorm:"book_id"`
	BookName string `json:"book_name" gorm:"book_name"`
	BorFreq int `json:"bor_freq" gorm:"bor_freq"` //用不到
	BrrowTime time.Time`json:"brrow_time" gorm:"brrow_time"`
	RetTime time.Time `json:"ret_time" gorm:"ret_time"`
	IsOver int `json:"is_over" gorm:"is_over"`
	Status int `json:"status" gorm:"status"`//审核状态0 未通过  1 已通过 2 已归还
}

type GetBorrowRecordsRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		Total int `json:"total"`
		Book []Borrows `json:"book"`
	} `json:"data"`

}

//获取借阅记录
func GetBorrowRecords(c*gin.Context)  {
	req:=GetBorrowRecordsReq{}
	rsp:=GetBorrowRecordsRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("GetBorrowRecords bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}
	var bor model.Borrow
	var bookRecords []Borrows
	borrows,err:=bor.GetBorrowByType(data.Db,req.Name,req.Isbn,req.PageNum,req.PageSize)
	if err!=nil{
		fmt.Println("GetBorrowRecords  GetBorrowByType model err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "查询失败"
		c.JSON(200,rsp)
		return
	}

	//查书名
	var book model.Book
	for _,v:=range borrows{
		ret,err := book.GetBookByIsbn(data.Db,v.BookId)
		if err!=nil{
			fmt.Println("GetBorrowRecords GetBookByIsbn error")
			return
		}
		bookRecords = append(bookRecords,Borrows{
			UserId :v.UserId,
			BookId:v.BookId,
			BookName :ret.BookName,
			BorFreq :v.BorFreq,
			BrrowTime :v.BrrowTime,
			IsOver :v.IsOver,
			Status :v.Status,
			RetTime:v.RetTime,
		})
	}

	total,err:=bor.GetBorrowTotalByType(data.Db,req.Name,req.Isbn)
	if err!=nil{
		fmt.Println("GetBorrowRecords  GetBorrowTotal model err"+err.Error())
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}

	rsp.Data.Book = bookRecords
	rsp.Data.Total = total
	rsp.Status = "200"
	c.JSON(200,rsp)
	return
}



//获取个人借阅记录

type GetUserBorrowRecordReq struct {
	Isbn string `json:"isbn" form:"isbn"`
	UserName  string `json:"user_name" form:"user_name"` //用户名
	PageNum int `json:"pagenum" form:"pagenum"`
	PageSize int `json:"pagesize" form:"pagesize"`
}

type GetUserBorrowRecordRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		Total int `json:"total"`
		Book []Borrows `json:"book"`
	} `json:"data"`

}


func GetUserBorrowRecord(c*gin.Context)  {
	req:=GetUserBorrowRecordReq{}
	rsp:=GetUserBorrowRecordRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("GetUserBorrowRecord bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}
	fmt.Println("GetUserBorrowRecordReq GetUserBorrowRecordReq ",req)
	var bor model.Borrow
	var bookRecords []Borrows
	borrows,err:=bor.GetUserBorrowRecord(data.Db,req.UserName,req.Isbn,req.PageNum,req.PageSize)
	if err!=nil{
		fmt.Println("GetUserBorrowRecord  GetUserBorrowRecord model err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "查询失败"
		c.JSON(200,rsp)
		return
	}

	//查书名
	var book model.Book
	for _,v:=range borrows{
		ret,err := book.GetBookByIsbn(data.Db,v.BookId)
		if err!=nil{
			fmt.Println("GetUserBorrowRecord GetBookByIsbn error")
			return
		}
		bookRecords = append(bookRecords,Borrows{
			UserId :v.UserId,
			BookId:v.BookId,
			BookName :ret.BookName,
			BorFreq :v.BorFreq,
			BrrowTime :v.BrrowTime,
			IsOver :v.IsOver,
			Status :v.Status,
			RetTime:v.RetTime,
		})
	}

	total,err:=bor.GetBorrowTotalByUserId(data.Db,req.UserName,req.Isbn)
	if err!=nil{
		fmt.Println("GetUserBorrowRecord  GetBorrowTotal model err"+err.Error())
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}

	rsp.Data.Book = bookRecords
	rsp.Data.Total = total
	rsp.Status = "200"
	c.JSON(200,rsp)
	return
}





type DelBorRecordReq struct {
	UserId string `json:"user_id"`
	BookId string `json:"book_id"`
}


type DelBorRecordRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {

	} `json:"data"`

}

//删除记录
func DelBorRecord(c*gin.Context)  {
	req:=DelBorRecordReq{}
	rsp:=DelBorRecordRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("DelBorRecord bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}

	var bor model.Borrow
	//软删除
	err := bor.DelBorrowByIdSoft(data.Db,req.UserId,req.BookId)
	if err!=nil{
		fmt.Println("DelBorRecord model DelBorrowById err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "200"
	rsp.Description = "success"
	c.JSON(200,rsp)
	return
}


type RevertBookReq struct {
	UserId string `json:"user_id"`
	BookId string `json:"book_id"`
}


type RevertBookRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {

	} `json:"data"`

}



//管理员借出图书 (更改状态)
func BorrowAddRecord(c*gin.Context)  {
	//使用上面一样的数据结构
	req:=RevertBookReq{}
	rsp:=RevertBookRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("BorrowAddRecord bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}
	var bor model.Borrow

	//更改状态
	err :=bor.BorrowBook(data.Db,req.UserId,req.BookId)
	if err!=nil{
		fmt.Println("BorrowAddRecord BorrowBook model  err"+err.Error())
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "200"
	rsp.Description = "success"
	c.JSON(200,rsp)
	return
}


//管理员还书
func RevertBook(c*gin.Context)  {
	req:=RevertBookReq{}
	rsp:=RevertBookRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("RevertBook bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}

	var book model.Book

	var bor model.Borrow
	//增加库存
	err :=book.AddBookStockByIsbn(data.Db,req.BookId)
	if err!=nil{
		//data.Db.Rollback()
		fmt.Println("RevertBook model AddBookStockByIsbn err"+err.Error())
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}

	err =bor.RevertBook(data.Db,req.UserId,req.BookId)
	if err!=nil{
		fmt.Println("RevertBook RevertBook model  err"+err.Error())
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "200"
	rsp.Description = "success"
	c.JSON(200,rsp)
	return
}