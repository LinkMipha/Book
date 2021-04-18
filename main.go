package main

import (
	"Book/gin_server"
)
 var listen string = "7001"

func main(){

	//data.SetKey("link","world",uint64(3600))
	//data.SetKey("link","mipha")
	//data.GetKey("link")
	//data.WatchKey("link")
	//启动服务
	gin_server.StartHttpServer(listen)

}