package gin_server

import (
	"Book/MiddleWare"
	"Book/data"
	"Book/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type HelloReq struct {
	Id int `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

type HelloRsp struct {
	Rsp struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			Name string `json:"name"`
			Id   int    `json:"id"`
		} `json:"data"`
	}
}




type CheckOutRsp struct {
	Status string `json:"status"`
	Description string `json:"description"`
	Id string `json:"_id"`
	UpdateTime  string `json:"updatedAt"`
	CreatedAt  string `json:"createdAt"`
	UserName string `json:"userName"`
	Name string `json:"name"`
	IsAdmin bool `json:"isAdmin"`
}

func GetTest(c*gin.Context)  {
	req:= HelloReq{}
	rsp:= HelloRsp{}
	if err:=c.Bind(&req);err!=nil{
		log.Fatal("test") //服务断开
	}
	fmt.Println(req)
	rsp.Rsp.Data.Name = req.Name
	rsp.Rsp.Data.Id = req.Id
	c.JSON(200,rsp)
}

func GetMenu(c *gin.Context) {

	c.JSON(200, "")
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
	rsp.Data.Token = token
	rsp.Status = "200"
	rsp.Description = "success"
	c.JSON(200,rsp)
}


type ChildRen struct {
	Id int `json:"id"`
	Name string `json:"authName"`
	Path string `json:"path"`
}

type Menu struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Order string `json:"order"`
	Path string `json:"path"`
	ChildRen []ChildRen `json:"children"`
}


type GetMenusRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
		Menus  [] Menu `json:"menus"`
	} `json:"data"`
}


//根据具体的身份获取相应的菜单
func GetMenus(c *gin.Context)  {

	rsp:=GetMenusRsp{}
	rsp.Status = "200"
	rsp.Description = "success"
	menus :=make([]Menu,0)
	for i:=0;i<3;i++{
		menus = append(menus,Menu{
			Id: i,
			Name: "简易菜单",
			Order: "order",
			ChildRen: []ChildRen{
				{
					Id: i,
					Name: "mipha",
					Path: "user",
				},

			},
		})
	}
	rsp.Data.Menus = menus

	c.JSON(200,rsp)
}




//?验证
func CheckOut(c *gin.Context) {
	rsp := CheckOutRsp{}
	rsp.IsAdmin = true
	rsp.Name = "管理员"
	rsp.CreatedAt = "2021-04-07T02:51:16.008Z"
	rsp.UpdateTime = "2021-04-07T02:51:16.008Z"
	rsp.Id = "606d1e24eebc850a6c30d1c1"
	c.JSON(200, rsp)
}


type Users struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Status      int `json:"status"`
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
	 		Id: v.Id,
			Name:v.Name,
			Status: v.Status,
		})
	 }

	rsp.Status = "200"
	rsp.Data.Total = total
	c.JSON(200,rsp)
}