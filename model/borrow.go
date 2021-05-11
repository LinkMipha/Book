package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type Borrow struct {
	Id int `json:"id" gorm:"id"`
	UserId string `json:"user_id" gorm:"user_id"`
	BookId string `json:"book_id" gorm:"book_id"`
	BorFreq int `json:"bor_freq" gorm:"bor_freq"` //用不到
	BrrowTime time.Time `json:"brrow_time" gorm:"brrow_time"`
	RetTime time.Time `json:"ret_time" gorm:"ret_time"`
	RealTime time.Time `json:"real_time" gorm:"real_time"`
	IsOver int `json:"is_over" gorm:"is_over"`
	IsDel int `json:"is_del" gorm:"is_del"`
	Status int `json:"status" gorm:"status"`//审核状态0 未通过  1 已通过 2已经归还
}

func (b*Borrow)TableName() string {
	return "borrow"
}


//增加数据
func (b*Borrow)AddBorrow(db*gorm.DB,data Borrow) error{
	err:=db.Table(b.TableName()).Create(&data).Error
	return err
}


//根据userId bookId 查询是否借阅过 ErrRecordNotFound
func (b*Borrow)GetBorrowByUserAndBookId(db*gorm.DB,userId string,bookId string) (Borrow,error){
	var bor Borrow
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).Where("is_del = 0").First(&bor).Error
	return bor,err
}


//类型获取数据量 模糊搜索用户 用户名
func (b*Borrow)GetBorrowTotalByType(db*gorm.DB,userId string,bookId string)   (int ,error){
	var total int
	sql:=db.Table(b.TableName()).Where("is_del = ?",0)

	if userId!=""{
		sql = sql.Where("user_id like ?","%"+userId+"%")
	}
	if bookId!=""{
		sql = sql.Where("book_id = ?",bookId)
	}
	err:=sql.Count(&total).Error
	return total,err
}


//类型获取数据量 精确搜索用户 userId
func (b*Borrow)GetBorrowTotalByUserId(db*gorm.DB,userId string,bookId string) (int ,error) {
	var total int
	sql:=db.Table(b.TableName()).Where("is_del = ?",0)

	if userId!=""{
		sql = sql.Where("user_id = ?",userId)
	}
	if bookId!=""{
		sql = sql.Where("book_id = ?",bookId)
	}
	err:=sql.Count(&total).Error
	return total,err

}

//根据Userid分页查询借阅记录  （增加status查询已经借出去的书本）
func  (b*Borrow)GetBorrowByType(db*gorm.DB,userId string,bookId string,pageIndex int,pageSize int)([]Borrow,error)  {
	var borrow []Borrow
	var err error
	sql:=db.Table(b.TableName()).Where("is_del = ?",0)

	if userId!=""{
		sql = sql.Where("user_id like ?","%"+userId+"%")
	}
	if bookId!=""{
		sql = sql.Where("book_id = ?",bookId)
	}
	//跳过条数
	if pageIndex > 0 && pageSize > 0 {
		err = sql.Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&borrow).Error
	} else {
		//默认20条数据
		err = sql.Limit(20).Find(&borrow).Error
	}
	return borrow, err
}


//个人借阅记录分页获取
func (b*Borrow)GetUserBorrowRecord(db*gorm.DB,userId string,bookId string,pageIndex int,pageSize int) ([]Borrow,error) {
	var borrow []Borrow
	var err error
	sql:=db.Table(b.TableName()).Where("is_del = ?",0)

	if userId!=""{
		sql = sql.Where("user_id = ?",userId)
	}
	if bookId!=""{
		sql = sql.Where("book_id = ?",bookId)
	}
	//跳过条数
	if pageIndex > 0 && pageSize > 0 {
		err = sql.Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&borrow).Error
	} else {
		//默认20条数据
		err = sql.Limit(20).Find(&borrow).Error
	}
	return borrow, err
}



//获取数据总数量
func (b*Borrow)GetBorrowTotal(db*gorm.DB)(int,error)  {
	total:=0
	err:=db.Table(b.TableName()).Count(&total).Error
	return total,err
}






//删除借阅记录(软删除保留记录)
func (b*Borrow)DelBorrowByIdSoft(db*gorm.DB,userId string,bookId string) error{
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).Where("is_del = 0").Update("is_del",1).Error
	return err
}


//删除借阅记录
func (b*Borrow)DelBorrowByIdHard(db*gorm.DB,userId string,bookId string) error{
	var bor Borrow
	db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).First(&bor)
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).Delete(&bor).Error
	return err
}

//同意借阅
func (b*Borrow)AgreeBorrowById(db*gorm.DB,userId string,bookId string) error{
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).Where("is_del = 0").Update("status",1).Error
	return err
}


//根据userId查询借阅列表 逾期，未逾期 通过未通过都算上
func (b*Borrow)GetBorrowByUserId(db*gorm.DB,userId string)([]Borrow,error){
	var borrows []Borrow
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("is_del = 0").Find(&borrows).Error
	return borrows,err
}

//根据user_id查询逾期的书籍
func (b*Borrow)GetOverBorrowByUserId(db*gorm.DB,userId string)([]Borrow,error){
	var borrows []Borrow
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("is_over = 1").Find(&borrows).Error
	return borrows,err
}

//还书 更改 status = 1改为status = 2
func (b*Borrow)RevertBook (db*gorm.DB,userId string, bookId string)error {
	status:=make(map[string]interface{})
	status["status"] = 2
	status["is_over"] = 0
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).Where("status = ?",1).Update(status).Error
	return err
}

//借出 更改status = 0 改为 status = 1
func (b*Borrow)BorrowBook (db*gorm.DB,userId string, bookId string)error {
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).Where("status = ?",0).Update("status",1).Error
	return err
}



//用户续借
func (b*Borrow)ReNewBorrow(db*gorm.DB,userId string, bookId string)error  {
	//更改时间
	var book Borrow
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).Where("is_del = ?",0).Where("status = ?",1).First(&book).Error
	if err!=nil{
		return err
	}
	if book.BorFreq>=2{
		return  errors.New("续借次数上限，归还后再次借阅")
	}
	freq:=book.BorFreq+1

	renewTiem:=book.RetTime.AddDate(0,1,0)
	err=db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).Where("is_del = ?",0).Where("status = ?",1).Update("ret_time",renewTiem).Update("bor_freq",freq).Error
	if err!=nil{
		return err
	}
	return err
}


//查找借阅的书
//根据userId查询借阅列表 逾期，未逾期 通过未通过都算上
func (b*Borrow)GetBorrowByUserIdBookId(db*gorm.DB,userId string,bookId string)([]Borrow,error){
	var borrows []Borrow
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("book_id = ?",bookId).Where("is_del = 0").Find(&borrows).Error
	return borrows,err
}

//用户借阅的书籍，通过借阅且未归还
func (b*Borrow)GetBorRecordByUserId(db*gorm.DB,userId string)([]Borrow,error){
	var borrows []Borrow
	err:=db.Table(b.TableName()).Where("user_id = ?",userId).Where("is_del = ?",0).Where("status = ?",1).Find(&borrows).Error
	return borrows,err
}


//批量获取未删除借阅数据  只计算已经借出去的
func (b*Borrow)GetCountRecord(db*gorm.DB,lastId int)([]Borrow,error){
	var borrows []Borrow
	err:=db.Table(b.TableName()).Where("id > ?",lastId).Where("is_del = ?",0).Where("status = ?",1).Order("id asc").Limit(100).Find(&borrows).Error
	return borrows,err
}


//更新数据
func (b*Borrow)UpdateBorrow(db*gorm.DB, id int,ma map[string]interface{})error {
	err := db.Table(b.TableName()).Where("id = ?",id).Update(ma).Error
	return  err
}

//获取借阅最多次数的四本图书
//func (b*Borrow)GetTopBooks(db*gorm.DB) (error,Borrow) {
//
//}