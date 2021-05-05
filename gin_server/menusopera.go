package gin_server

import (
	"Book/data"
	"Book/model"
	"fmt"
	"github.com/gin-gonic/gin"
)



type ChildRen struct {
	Id int `json:"id"`
	Name string `json:"authName"`
	Path string `json:"path"`
}


type Menu struct {
	Id int `json:"id"`
	Name string `json:"name"`
	//Order string `json:"order"`
	Path string `json:"path"`
	ChildRen []ChildRen `json:"children"`
}

type GetMenusReq struct {
	UserName string `json:"user_name" form:"user_name"`
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
	req:=GetMenusReq{}
	rsp:=GetMenusRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("GetMenus err"+err.Error())
		rsp.Status = "400"
		c.JSON(200,rsp)
		return
	}


	//查询用户是否为管理员
	var user model.User
	userData,err :=user.GetUserByUserName(data.Db,req.UserName)
	if err!=nil{
		rsp.Status= "400"
		c.JSON(2000,rsp)
		return
	}
	menus :=make([]Menu,0)
 	var menusModel model.Menus
    fmt.Println("用户权限",userData)
	parents,err:=menusModel.GetParentMenus(data.Db,userData.IsAdmin)
	fmt.Println("用户权限菜单",parents)
	if err!=nil{
		fmt.Println("GetMenus")
		rsp.Status= "400"
		c.JSON(2000,rsp)
		return
	}
	for _,v:=range parents{
		children,err:=menusModel.GetByParentId(data.Db,v.Id,user.IsAdmin)
		if err!=nil{
			fmt.Println("GetByParentId err",err)
		}
		//children装入数据结构
		var childMenus []ChildRen
		for _,v:=range children{
			childMenus = append(childMenus,ChildRen{
				Id: v.Id,
				Name: v.Name,
				Path: v.Path,
			})
		}
		menus = append(menus,Menu{
			Id: v.Id,
			Name: v.Name,
			Path: v.Path,
			ChildRen: childMenus,
		})
	}

	rsp.Data.Menus = menus
	rsp.Status = "200"
	rsp.Description = "success"
	c.JSON(200,rsp)
}

