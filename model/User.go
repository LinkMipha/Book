package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

//用户信息
type User struct {
	Id         int  `json:"id" gorm:"id"`
	UserName string `json:"userName" gorm:"userName"`
	PassWord string `json:"password" gorm:"password"`
	Sex int `json:"sex" gorm:"sex"`
	Name string `json:"name" gorm:"name"`
	IsAdmin int `json:"isAdmin" gorm:"isAdmin"`
	Status int `json:"status" gorm:"status"` //1删除
}

func (u *User) TableName() string {
	return "users"
}


//获取所有用户数量
func (u *User)GetUserTotal(db *gorm.DB,name string)(total int, err error){
	var sqlStr string
	 if name==""{
		 err = db.Table(u.TableName()).Where("status  = ?",0).Count(&total).Error

	 }else{
		 sqlStr = "select count(*) total from users where name like ? and status = 0"
		 err = db.Table(u.TableName()).Raw(sqlStr, name).Count(&total).Error
	 }
	 return total,err
}

//分页模糊获取用户
func (u *User) GetUsersList(db *gorm.DB, pageIndex int, pageSize int,name string) ([]User, error) {
	var users []User
	var err error
	if pageSize > 0 && pageIndex > 0 {
		var sqlStr string
		if name==""{
			sqlStr = "select *from users where  status = 0 limit ?,?"
			err = db.Table(u.TableName()).Raw(sqlStr, (pageIndex - 1) * pageSize,pageSize).Find(&users).Error
		}else{
			sqlStr = "select *from users where name like ? and status = 0 limit ?,?"
			err = db.Table(u.TableName()).Raw(sqlStr, name,(pageIndex - 1) * pageSize,pageSize).Find(&users).Error
		}
		//分页
		//err = db.Table(u.TableName()).Where(fmt.Sprintf(" userName like '%%%s' ", name)).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&users).Error
	} else {
		//默认20
		err = db.Table(u.TableName()).Limit(20).Offset(0).Find(&users).Error
	}
	fmt.Println(users)
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
