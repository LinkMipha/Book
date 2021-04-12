package model

import "github.com/jinzhu/gorm"

//书籍操作
type Book struct {
	Id          int64   `json:"id" gorm:"id"`
	Isbn        string  `json:"isbn" gorm:"isbn"`
	BookName    string  `json:"book_name" gorm:"book_name"`
	Author      string  `json:"author" gorm:"author"`
	Publish     string  `json:"publish" gorm:"publish"`
	PublishTime string  `json:"publish_time" gorm:"publish_time"`
	Price       float32 `json:"price" gorm:"price"`
	BookType    string  `json:"book_type" gorm:"book_type"` //类型
	Stock       int     `json:"stock" gorm:"stock"`
	UploadTime  string  `json:"upload_time" gorm:"upload_time"`
	CreateTime  string  `json:"create_time" gorm:"create_time"`
	UpdateTime  string  `json:"update_time" gorm:"update_time"`
}

func (b *Book) TableName() string {
	return "biz_book"
}

//创建
func (b *Book) CreateBook(db *gorm.DB) error {
	return db.Table(b.TableName()).Create(b).Error
}

//删除
func (b *Book) DeleteBook(db *gorm.DB, id int64) error {
	return db.Table(b.TableName()).Where("id = ?", id).Delete(b).Error
}

//更新信息
func (b *Book) UpdateBook(db *gorm.DB, id int64, message map[string]interface{}) error {
	return db.Table(b.TableName()).Where("id = ?", id).Update(message).Error
}

//根据id查询书籍信息
func (b *Book) GetBook(db *gorm.DB, Id int64) (book Book, err error) {
	err = db.Table(b.TableName()).Where("id = ?", Id).First(&book).Error
	return book, err
}

//根据名字模糊查询
func (b *Book) GetBookByName(db *gorm.DB, name string) (book Book, err error) {
	sqlStr := "select *from biz_book where book_name like ?"
	err = db.Table(b.TableName()).Raw(sqlStr, name).Find(&book).Error
	return book, err
}

//根据类型分页获取书籍 一页 20
func (b *Book) GetBookByType(db *gorm.DB, bookType string, pageIndex int, pageSize int) ([]Book, error) {
	var books []Book
	var err error
	sql := db.Table(b.TableName()).Where("book_type = ?", bookType)
	//跳过条数
	if pageIndex > 0 && pageSize > 0 {
		err = sql.Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&books).Error
	} else {
		//默认20条数据
		err = sql.Limit(20).Find(&books).Error
	}
	return books, err
}
