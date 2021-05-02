package gin_server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	//后续增加日志
	_ "github.com/sirupsen/logrus"
	"io"
	_ "net/http/pprof"
	"os"
)

const port = 20080


// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		//指定端口跨域
		//c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func StartHttpServer(listen string)  {

	//控制台颜色
	gin.ForceConsoleColor()
	//日志
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	ginRouter :=gin.Default() //暂时使用default
	ginRouter.Use(Cors())
	//增加中间件
	//router.Use(middleware.GetUseTime)

	router:=ginRouter.Group("/")
	router.POST("renren-fast/sys/login", GetTest)
	router.GET("renren-fast/sys/menu/nav",GetMenu)
	router.GET("renren-fast/sys/user/info",GetMenu)

	//登陆
	router.POST("api/login",LoginIn)


	//获取菜单操作  。。。。。。。。。。。。。。。
	router.GET("api/menus",GetMenus)
	//检查登陆状态
	router.GET("api/login/checkcode",CheckOut)



	//用户相关操作 之后移动到其他服务......................
	router.GET("api/users",GetUserList)

	//增加
	router.POST("api/adduser",AddUser)

	//用户名搜索
	router.GET("api/get_user_by_username/:userName",GetUserByUserName)

	//更新
	router.PUT("api/edituser/:userName",EditUserByUserName)


	router.POST("api/deleteUser",DeleteUserByUserName)



	//图书相关操作之后移动到其他服务 。。。。。。。。。。。。。。


	router.GET("api/books",GetBookList)

	//增加
	router.POST("api/addBook",AddNewBook)

	//图书修改时查询信息
	router.GET("api/get_book_by_isbn/:isbn",GetBookByIsbn)

	//更新
	router.PUT("api/editbook/:isbn",EditBookByIsbn)


	router.POST("api/deleteBook",DeleteUserByIsbn)




	//性能测试
	go func() {
		fmt.Println("pprof start")
		fmt.Println(http.ListenAndServe(":9876",nil))
	}()

	//start
	if err := ginRouter.Run(fmt.Sprintf(":%s",listen)); err != nil {
		println("Error when running server. " + err.Error())
	}

}