package main

import (
	"Book/gin_server"
	"time"
)
 var listen string = "7001"

func main(){

	//data.SetKey("link","world",uint64(3600))
	//data.SetKey("link","mipha")
	//data.GetKey("link")
	//data.WatchKey("link")
	//启动服务
	//go gin_server.SendEmail()

	//定时任务计算是否逾期
	go gin_server.StartTimerBegin(time.Now(),gin_server.CountBookTime(0,0))
	gin_server.StartHttpServer(listen)

}