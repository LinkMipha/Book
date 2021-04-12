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

		c.Header("Access-Control-Allow-Origin", "http://localhost:8001")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
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
	router:=gin.Default()//暂时使用default
	router.Use(Cors())
	//增加中间件
	//router.Use(middleware.GetUseTime)

	st:=router.Group("/")
	st.POST("renren-fast/sys/login", GetTest)
	st.GET("renren-fast/sys/menu/nav",GetMenu)
	st.GET("/get_book_by_id",GetBookById)
	st.GET("renren-fast/sys/user/info",GetMenu)

	//性能测试
	go func() {
		fmt.Println("pprof start")
		fmt.Println(http.ListenAndServe(":9876",nil))
	}()

	//start
	if err := router.Run(fmt.Sprintf(":%s",listen)); err != nil {
		println("Error when running server. " + err.Error())
	}

}