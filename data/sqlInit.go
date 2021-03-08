package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-server/conf"
	"log"
	"time"
)

var Db *gorm.DB

func InitMysql(config conf.Config)  {
	fmt.Println("Load dbService config...")

	//设置连接参数
	dbType :=config.Database.Type
	usr := config.Database.User
	pwd := config.Database.Password
	address := config.Database.Address
	dbName := config.Database.DbName
	dbLink := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		usr, pwd, address, dbName)
	//失败重试
	var err error
	for Db, err = gorm.Open(dbType, dbLink); err != nil; Db, err = gorm.Open(dbType, dbLink) {
		log.Println("Failed to connect database: ", err.Error())
		log.Println("Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	Db.DB().SetMaxIdleConns(config.Database.MaxIdle)
	Db.DB().SetMaxOpenConns(config.Database.MaxOpen)
}