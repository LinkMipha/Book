package model

import "github.com/jinzhu/gorm"

type Menus struct {
	Id int `json:"id" gorm:"id"`
	ParentId int `json:"parent_id" gorm:"column:parent_id"`
	Name string `json:"name" gorm:"name"`
	Path string `json:"path" gorm:"path"`
	Status int `json:"status" gorm:"status"`
}

func (m*Menus)TableNmae() string {
	return "menus"
}

//根据父类id进行查询
func (m*Menus)GetByParentId(db*gorm.DB,parentId int)([]Menus,error)  {
	var menus []Menus
	var err error
	err  = db.Table(m.TableNmae()).Where("parent_id = ?",parentId).Find(&menus).Error
	return menus,err
}

//查询parent_id 为0的
func (m*Menus)GetParentMenus(db*gorm.DB)([]Menus,error) {
	var menus []Menus
	var err error
	err  = db.Table(m.TableNmae()).Where("parent_id = ?",0).Find(&menus).Error
	return menus,err
}