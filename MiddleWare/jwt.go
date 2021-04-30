package MiddleWare

//引入路径名，而不是包名
import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const ErrMsg = "errMsg"

var (
	SignKey = "Secret Key"
	Auth    = "Authorization"
	Issuer  = "this is a issuer"
)

//接口调用时间
func GetUseTime(ctx *gin.Context) {
	begin := time.Now()
	ctx.Next()
	t := time.Since(begin)
	fmt.Printf("Use time :%v", t)
}

//鉴权
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(Auth)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{ErrMsg: "Not Authorized"})
			c.Abort()
			return
		}

		log.Println("JwtAuth token", token)

		j := NewJwt()
		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{ErrMsg: "failed Authorized"})
			c.Abort()
			return
		}
		c.Set("claims", claims)
	}
}

// Jwt 签名
type Jwt struct {
	SigningKey []byte
}

//载荷，可以加信息，透传下去
type CustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Kind     int    `json:"kind"`
	jwt.StandardClaims
}

func NewJwt() *Jwt {
	return &Jwt{
		[]byte(GetSignKey()),
	}
}
func GetSignKey() string {
	return SignKey
}

func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

//生成token
func (j *Jwt) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(j.SigningKey)
}

//解析token
func (j *Jwt) ParseToken(tokenString string) (*CustomClaims, error) {
	//func返回定义的key，检测token中的key是否正确
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	//错误暂时未分类处理
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	err = errors.New("Couldn't handle this token ")
	return nil, err

}

//更新token
func (j *Jwt) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(1) * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", errors.New("Couldn't handle this token ")
}
