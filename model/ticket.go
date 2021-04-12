package model

import (
	"github.com/jinzhu/gorm"
)

//逾期表
type Ticket struct {
	Id         int    `json:"id" gorm:"id"`
	UserId     string `json:"user_id" gorm:"user_id"`
	BookId     string `json:"book_id" gorm:"book_id"`
	OverId     int    `json:"over_id" gorm:"over_id"`
	TicketFee  string `json:"ticket_fee" gorm:"ticket_fee"`
	CreateTime string `json:"create_time" gorm:"create_time"`
}

func (t *Ticket) TableName() string {
	return "ticket"
}

//根据userId获取逾期信息
func (t *Ticket) GetTicketsByUserId(db *gorm.DB, userId string) ([]Ticket, error) {
	var ticket []Ticket
	var err error
	err = db.Table(t.TableName()).Where("user_id = ?", userId).Find(&ticket).Error
	return ticket, err
}

//分页显示逾期信息
func (t *Ticket) GetTicketsList(db *gorm.DB, pageIndex int, pageSize int) ([]Ticket, error) {
	var ticket []Ticket
	var err error
	if pageSize > 0 && pageIndex > 0 {
		err = db.Table(t.TableName()).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&ticket).Error
	} else {
		//默认20
		err = db.Table(t.TableName()).Limit(20).Offset(0).Find(&ticket).Error
	}
	return ticket, err
}
