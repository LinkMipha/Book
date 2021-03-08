package gin_server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go-server/basic"
	"go-server/data"
	"go-server/middleware"
	"go-server/model"
	"log"
	"net/http"
	"time"
)

type LoginAuthReq struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Kind int `json:"kind" form:"kind"`
}

type LoginAuthRsp struct {
	Rsp struct{
		Status string `json:"status"`
		Description string `json:"description"`
		Data struct{

		} `json:"data"`
	}
}

//登陆
func LoginAuth(c *gin.Context)  {
	req:=LoginAuthReq{}
	rsp:=LoginAuthRsp{}.Rsp.Data
	if err:=c.Bind(&req);err!=nil{
		log.Println("param error",err.Error())
		basic.ResponseError(c,rsp,"param err"+err.Error())
		return
	}

	user,err:=new(model.User).GetUserByName(data.Db,req.Username)
	if err==gorm.ErrRecordNotFound{
		log.Println("user not exist")
		basic.ResponseError(c,rsp,"userName not exist")
		return
	}
	if err!=nil{
		log.Println("GetUserByName err: ",err.Error())
		basic.ResponseError(c,rsp,"GetUserByName error"+err.Error())
		return
	}
	
	if user.Password!=model.GetMd5(req.Password){
		log.Println("password error")
		basic.ResponseError(c,rsp,"password error")
		return
	}
	
	//生成令牌token
	GetToken(c,req.Username,req.Password,req.Kind)
}

func GetToken(c*gin.Context,userName string,passWord string,kind int)  {
	j:=middleware.NewJwt()
	claims:=middleware.CustomClaims{
		Username: userName,
		Password: passWord,
		Kind: kind,
		StandardClaims:jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix()-1000),
			ExpiresAt: int64(time.Now().Unix()+3600),
			Issuer: middleware.Issuer,
		},
	}
	token,err:=j.CreateToken(claims)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"kind":kind,
			"errMsg":err,
		})
		return
	}
	c.Header("Authorization",token)
	c.JSON(http.StatusOK,gin.H{
		"kind":kind,
		"errMsg":"",
	})
	return
}

