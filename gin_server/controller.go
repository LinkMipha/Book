package gin_server

import (
	"Book/MiddleWare"
	"Book/data"
	"Book/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
)


type RegisterReq struct {
	UserName string `json:"username" form:"username"`
	PassWrod string `json:"password" from:"password"`
}


type RegisterRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {

	} `json:"data"`
}

func Register(c*gin.Context)  {
	req:=RegisterReq{}
	rsp:=RegisterRsp{}
	if err:=c.BindJSON(&req);err!=nil{
		log.Printf("Register bind error:%v",err)
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}

	var user model.User
	_,err:=user.GetUserByUserName(data.Db,req.UserName)
	if err!=nil&&err!=gorm.ErrRecordNotFound{
		log.Printf("Register GetUserByUserName error:%v",err)
		rsp.Description = "注册失败"
		c.JSON(200,rsp)
		return
	}
	if err!=gorm.ErrRecordNotFound{
		log.Printf("Register GetUserByUserName error:%v",err)
		rsp.Description = "用户名已存在"
		c.JSON(200,rsp)
		return
	}
    user.UserName = req.UserName
    user.PassWord = req.PassWrod
	err = user.AddUserByMessage(data.Db,user)
	if err!=nil{
		log.Printf("Register AddUserByMessage error:%v",err)
		rsp.Description = "注册失败"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "200"
	rsp.Description = "注册成功"
	c.JSON(200,rsp)
}




type LoginInReq struct {
	UserName string `json:"username" form:"username"`
	PassWrod string `json:"password" from:"password"`
}

type LoginInRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
		Token string `json:"token"`
		IsAdmin int `json:"is_admin"`
	} `json:"data"`
}

//登陆验证
func LoginIn(c *gin.Context)  {
	req:= LoginInReq{}
	rsp:=LoginInRsp{}

	if err:=c.BindJSON(&req);err!=nil{
		log.Printf("LoginIn bind error:%v",err)
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}

	user:=model.User{}
	user,err:=user.GetUserByUserName(data.Db,req.UserName)
	if err!=nil{
		fmt.Printf("LoginIn error:%v\n",err)
		rsp.Description = "GetUserByUserName error"
		c.JSON(200,rsp)
		return
	}
	fmt.Println(user)
	if user.PassWord!=req.PassWrod{
		rsp.Status = "401"
		rsp.Description = "password"
		c.JSON(200,rsp)
		return
	}

	//设置token并且返回
	var jwt = MiddleWare.NewJwt()

	//荷载信息
	claim:=MiddleWare.CustomClaims{
		Username: req.UserName,
		Password: req.PassWrod,
		Kind:12,
	}

	token,err:=jwt.CreateToken(claim)
	if err!=nil{
		fmt.Println("CreateToken error",err)
		return
	}
	rsp.Data.IsAdmin = user.IsAdmin
	rsp.Data.Token = token
	rsp.Status = "200"
	rsp.Description = "success"
	c.JSON(200,rsp)
}




type Users struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	UserName string `json:"userName"`
	//PhoneNumber string `json:"phone_number"`
	Sex  int `json:"sex"`
	Status      int `json:"status"`
	Phone string `json:"phone" form:"phone"`
}




type GetUserListReq struct {
	Name  string `json:"query" form:"query"`
	PageNum int `json:"pagenum" form:"pagenum"`
	PageSize int `json:"pagesize" form:"pagesize"`
}

type GetUserListRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		Total int `json:"total"`
		User []Users `json:"users"`

	} `json:"data"`

}
//用户相关代码

func GetUserList(c*gin.Context)  {
	req:=GetUserListReq{}
	rsp:=GetUserListRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("GetUserList err"+err.Error())
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}

	var user model.User
	 st,err:=user.GetUsersList(data.Db,req.PageNum,req.PageSize,req.Name)
	 if err!=nil{
	 	fmt.Println("GetUsersList err"+err.Error())
	 }

	 total,err:=user.GetUserTotal(data.Db,req.Name)
	 if err!=nil{
		 fmt.Println("GetUserTotal err"+err.Error())
	 }


	 for _,v:=range st{
	 	rsp.Data.User = append(rsp.Data.User,Users{
	 		UserName: v.UserName,
	 		Id: v.Id,
			Name:v.Name,
			Sex: v.Sex,
			Status: v.Status,
			Phone: v.Phone,
		})
	 }

	rsp.Status = "200"
	rsp.Data.Total = total
	c.JSON(200,rsp)
}


//添加用户 json会影响数据的解析
type AddUserReq struct {
	UserName string `json:"username" form:"username"`
	PassWord string `json:"password" form:"password"`
	Name string `json:"name" form:"name"`
	Sex string `json:"sex" form:"sex"` //被迫使用string
	Phone string `json:"phone" form:"phone"`
}

type AddUserRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {


	} `json:"data"`

}

func AddUser(c*gin.Context)  {
	req:=AddUserReq{}
	rsp:=AddUserRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("AddUser err"+err.Error())
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}
	var user model.User
	total,err:=user.GetUserIdByUserName(data.Db,req.UserName)
	if err!=nil{
		fmt.Println(" AddUser GetUserIdByUserName err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "添加失败"
		c.JSON(200,rsp)
		return
	}
	if total>0{
		rsp.Status = "401"
		rsp.Description = "用户名已存在"
		c.JSON(200,rsp)
	}


	var addUser model.User
	addUser.UserName = req.UserName
	addUser.PassWord = req.PassWord
	addUser.Name = req.Name
	addUser.Phone = req.Phone
	if req.Sex=="男"{
		addUser.Sex = 0
	}else{
		addUser.Sex = 1
	}
	//添加用户
	err=user.AddUserByMessage(data.Db,addUser)
	if err!=nil{
		fmt.Println(" AddUser AddUserByMessage err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "添加失败"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "201"
	rsp.Description = "添加成功"
	c.JSON(200,rsp)
}

type GetUserByUserNameReq struct {
	UserName string `json:"userName" form:"userName"`
}


type GetUserByUserNameRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		User Users `json:"user"`
	} `json:"data"`
}


func GetUserByUserName(c*gin.Context)  {
    req:=GetUserByUserNameReq{}
    rsp:=GetUserByUserNameRsp{}
	//if err:=c.Bind(&req);err!=nil{
	//	fmt.Println("GetUserByUserName bind err"+err.Error())
	//	rsp.Status = "401"
	//	c.JSON(200,rsp)
	//	return
	//}

	//直接请求路径获取参数
	req.UserName = c.Param("userName")

	var user model.User
	user,err:=user.GetUserByUserName(data.Db,req.UserName)
	if err!=nil{
		fmt.Println("GetUserByUserName model err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "查询用户失败"
		c.JSON(200,rsp)
		return
	}
	rsp.Data.User.UserName = user.UserName
	rsp.Data.User.Name = user.Name
	rsp.Data.User.Sex = user.Sex
	rsp.Data.User.Phone = user.Phone
    rsp.Status = "200"
    c.JSON(200,rsp)
	return
}

type EditUserByUserNameReq struct {
	UserName string `json:"userName" form:"userName"`
	Name string `json:"name" form:"name"`
	Sex string `json:"sex" form:"sex"`
	Phone string `json:"phone" form:"phone"`
}


type EditUserByUserNameRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
	} `json:"data"`
}


//修改用户信息
func EditUserByUserName(c*gin.Context)  {
	req:=EditUserByUserNameReq{}
	rsp:=EditUserByUserNameRsp{}
	//url中获取数据
	req.UserName = c.Param("userName")
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("EditUserByUserName bind err"+err.Error())
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}
	editUser:=make(map[string]interface{},0)
	editUser["Name"] = req.Name
	editUser["phone"] = req.Phone
	if req.Sex=="男"{
		editUser["sex"] = 0
	}else{
		editUser["sex"] = 1
	}


	var user model.User
	err:=user.UpdatesUser(data.Db,req.UserName,editUser)
	if err!=nil{
		fmt.Println("EditUserByUserName UpdatesUser error ",err)
		rsp.Status = "401"
		c.JSON(200,rsp)
		return
	}

	rsp.Status = "200"
	rsp.Description = "修改成功"
	c.JSON(200,rsp)
	return
}



type DeleteUserByUserNameReq struct {
	UserName string `json:"userName" form:"userName"`
}

type DeleteUserByUserNameRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
	} `json:"data"`
}

//删除信息
func DeleteUserByUserName(c*gin.Context)  {
	req:=DeleteUserByUserNameReq{}
	rsp:=DeleteUserByUserNameRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("DeleteUserByUserName bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}
	var user model.User
	err:=user.DeleteUserById(data.Db,req.UserName)
	if err!=nil{
		fmt.Println("DeleteUserById error")
		rsp.Status = "401"
		rsp.Description = "DeleteUserById error"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "200"
	rsp.Description = "删除成功"
	c.JSON(200,rsp)
	return

}

type ResetUserPasswordReq struct {
	UserName string `json:"user_name"`
}


type ResetUserPasswordRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
	} `json:"data"`
}

func ResetUserPassword(c*gin.Context)  {
	req:=ResetUserPasswordReq{}
	rsp:=ResetUserPasswordRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("ResetUserPassword bind err"+err.Error())
		rsp.Status = "401"
		rsp.Description = "bind error"
		c.JSON(200,rsp)
		return
	}

	var user model.User
	err:=user.ResetPassword(data.Db,req.UserName)
	if err!=nil{
		fmt.Println("ResetPassword error")
		rsp.Status = "401"
		rsp.Description = "ResetPassword error"
		c.JSON(200,rsp)
		return
	}
	rsp.Status = "200"
	rsp.Description = "重置成功"
	c.JSON(200,rsp)
	return
}