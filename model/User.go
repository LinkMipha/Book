package model

import "github.com/jinzhu/gorm"

//用户信息
type User struct {
	Id         int  `json:"id" gorm:"id"`
	UserName string `json:"userName" gorm:"userName"`
	PassWord string `json:"password" gorm:"password"`
	Sex int `json:"sex" gorm:"sex"`
	Name string `json:"name" gorm:"name"`
	IsAdmin int `json:"isAdmin" gorm:"isAdmin"`
	Status int `json:"status" gorm:"Status"` //1删除
}

func (u *User) TableName() string {
	return "users"
}


//分页获取用户
func (u *User) GetUsersList(db *gorm.DB, pageIndex int, pageSize int) ([]User, error) {
	var users []User
	var err error
	if pageSize > 0 && pageIndex > 0 {
		err = db.Table(u.TableName()).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&users).Error
	} else {
		//默认20
		err = db.Table(u.TableName()).Limit(20).Offset(0).Find(&users).Error
	}
	return users, err
}

//查询用户
func (u *User) GetUserIdByUserId(db *gorm.DB, Name string) (User, error) {
	var user User
	var err error
	err = db.Table(u.TableName()).Where("name = ?", Name).Where("status = 0").Find(&user).Error
	return user, err
}

//标记删除用户
func (u *User) DeleteUserById(db *gorm.DB, name string) error {
	err := db.Table(u.TableName()).Where("name = ?", name).Update("status", 1).Error
	return err
}


//根据名字模糊查询
func (u *User) GetUsersByName(db *gorm.DB, name string) (user []User, err error) {
	sqlStr := "select *from biz_book where name like ? and status = 0"
	err = db.Table(u.TableName()).Raw(sqlStr, name).Find(&user).Error
	return user, err
}
