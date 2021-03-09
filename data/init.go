package data

import (
	"Book/conf"
)

//初始化配置可以使用etcd读取配置
func init()  {

	config,err:=conf.GetConfig()
	if err!=nil{
		panic("failed to load data config"+err.Error())
	}
	InitMysql(config)
	initRedis(config)
}
