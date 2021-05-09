package model

import "github.com/jinzhu/gorm"

type Menus struct {
	Id int `json:"id" gorm:"id"`
	ParentId int `json:"parent_id" gorm:"column:parent_id"`
	Name string `json:"name" gorm:"name"`
	Path string `json:"path" gorm:"path"`
	Status int `json:"status" gorm:"status"`
	IsAdmin int `json:"is_admin" gorm:"is_admin"`
}

func (m*Menus)TableNmae() string {
	return "menus"
}

//根据父类id进行查询
func (m*Menus)GetByParentId(db*gorm.DB,parentId int,isAdmin int)([]Menus,error)  {
	var menus []Menus
	var err error
	err  = db.Table(m.TableNmae()).Where("parent_id = ?",parentId).Find(&menus).Error
	return menus,err
}

//查询parent_id 为0的
func (m*Menus)GetParentMenus(db*gorm.DB,isAdmin int)([]Menus,error) {
	var menus []Menus
	var err error

	//不是管理员要限制
	sql:=db.Table(m.TableNmae())
	if isAdmin==0{
		sql = sql.Where("is_admin = ?",isAdmin)
	}
	err = sql.Where("parent_id = ?",0).Find(&menus).Error
	return menus,err
}