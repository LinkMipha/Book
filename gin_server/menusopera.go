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
	menus :=make([]Menu,0)
 	var menusModel model.Menus
	parents,err:=menusModel.GetParentMenus(data.Db)
	if err!=nil{
		fmt.Println("GetMenus")
		rsp.Status= "400"
		c.JSON(2000,rsp)
		return
	}
	for _,v:=range parents{
		children,err:=menusModel.GetByParentId(data.Db,v.Id)
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

