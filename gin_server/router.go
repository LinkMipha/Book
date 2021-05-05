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


	//注册
	router.POST("api/register",Register)

	//登陆
	router.POST("api/login",LoginIn)


	//获取菜单操作  。。。。。。。。。。。。。。。
	router.GET("api/menus",GetMenus)





	//用户相关操作 之后移动到其他服务......................
	router.GET("api/users",GetUserList)

	//增加
	router.POST("api/adduser",AddUser)

	//用户名搜索
	router.GET("api/get_user_by_username/:userName",GetUserByUserName)

	//更新
	router.PUT("api/edituser/:userName",EditUserByUserName)

    //删除用户
	router.POST("api/deleteUser",DeleteUserByUserName)

    //重置密码
    router.POST("api/reset_password",ResetUserPassword)

	//图书相关操作之后移动到其他服务 。。。。。。。。。。。。。。


	router.GET("api/books",GetBookList)

	//增加
	router.POST("api/addBook",AddNewBook)

	//图书修改时查询信息
	router.GET("api/get_book_by_isbn/:isbn",GetBookByIsbn)

	//更新
	router.PUT("api/editbook/:isbn",EditBookByIsbn)


	router.POST("api/deleteBook",DeleteUserByIsbn)


	//图书续借  用户点击
	router.POST("api/renew_borrow",ReNewBorrow)

	//借阅相关接口
	router.POST("api/add_borrow",AddBorrow)

	//获取总的借阅记录
	router.GET("api/get_borrow_records",GetBorrowRecords)

	//获取个人借阅记录
	router.GET("api/get_user_borrow_record",GetUserBorrowRecord)



	//删除借阅记录 管理员功能
	router.POST("api/del_borrow_record",DelBorRecord)

	//借出图书
	router.POST("api/verifyRecord",BorrowAddRecord)

    //还书 管理员功能
    router.POST("api/revert_record",RevertBook)



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