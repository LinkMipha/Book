package main

import (
	"Book/data"
	"Book/gin_server"
)
 var listen string = "9528"

func main(){

	//data.SetKey("link","world",uint64(3600))
	data.SetKey("link","mipha")
	data.GetKey("link")
	data.WatchKey("link")
	//启动服务
	gin_server.StartHttpServer(listen)

}