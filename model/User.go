package model

import "github.com/jinzhu/gorm"

//用户信息
type User struct {
	Id         int64  `json:"id" gorm:"id"`
	UserId     string `json:"user_id" gorm:"user_id"`
	UserName   string `json:"user_name" gorm:"user_name"`
	Mobile     string `json:"mobile" gorm:"mobile"`
	PassWord   string `json:"password" gorm:"password"`
	CreateTime string `json:"create_time" gorm:"create_time"`
	IsDelete   int    `json:"is_delete" gorm:"is_delete"` //软删除
}

func (u *User) TableName() string {
	return "user"
}

//查询用户
func (u *User) GetUserIdByUserId(db *gorm.DB, userId string) (User, error) {
	var user User
	var err error
	err = db.Table(u.TableName()).Where("user_id = ?", userId).Where("is_delete = 0").Find(&user).Error
	return user, err
}

//标记删除用户
func (u *User) DeleteUserById(db *gorm.DB, userId string) error {
	err := db.Table(u.TableName()).Where("user_id = ?", userId).Update("is_delete", 1).Error
	return err
}

//根据名字模糊查询
func (u *User) GetUsersByName(db *gorm.DB, name string) (user []User, err error) {
	sqlStr := "select *from biz_book where book_name like ? and is_delete = 0"
	err = db.Table(u.TableName()).Raw(sqlStr, name).Find(&user).Error
	return user, err
}
