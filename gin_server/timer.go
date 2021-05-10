package gin_server

import (
	"Book/data"
	"Book/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)


//定时任务相关定时器 每天0点计算逾期
func startZeroTimer(f func()) {
	go func() {
		for {
			//执行函数
			f()
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			//下一个零点距离现在的时长
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}

//手动设置启动时间，间隔24小时
func StartTimerBegin(begin time.Time,f func( int, int))  {
	go func(){
		if begin.IsZero(){
			begin = time.Now()
		}

		for{
			go f(0,0)
			//定时24小时
			next:=begin.Add(time.Hour*24)
			t:=time.NewTimer(next.Sub(begin))
			<-t.C
			//重新begin赋值
			begin = next
		}
	}()

}



//定时任务计算借阅是否逾期 准确时间
func CountBookTime(lastId int, Redis int) {
	//redis获取lastId，默认走redis
	begin:=time.Now()
	for {
		var borrow model.Borrow
		record, err := borrow.GetCountRecord(data.Db, lastId)
		if err != nil {
			fmt.Println("CountBookTime GetCountRecord err and lastId", lastId, err)
		}
		//逾期精准计算
		for _, v := range record {
			//现在时间在应归还日期之后,归还日期在应还日期之前
			if time.Now().After(v.RetTime) && v.IsDel == 0 && v.RealTime.Before(v.RetTime) {
				ma := make(map[string]interface{})
				ma["is_over"] = 1
				err = borrow.UpdateBorrow(data.Db, v.Id, ma)
				if err != nil {
					fmt.Println("CountBookTime UpdateBorrow err and lastId", lastId, err)
				}
			}
		}
		//上次的最后一条数据id
		lastId = record[len(record)-1].Id
		if time.Now().Sub(begin)>time.Hour*24{
			break
		}
	}
}


type GinCountBookTimeReq struct {
	LastId int `json:"last_id"`
}
type GinCountBookTimeRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        struct {
	} `json:"data"`
}

//手动调用逾期函数
func GinCountBookTime(c*gin.Context)  {
	req:=GinCountBookTimeReq{}
	rsp:=GinCountBookTimeRsp{}
	if err:=c.Bind(&req);err!=nil{
		fmt.Println("GinCountBookTime err bind ",err)
		return
	}
	//手动调用计算逾期
	go CountBookTime(req.LastId,0)

	rsp.Status = "200"
	rsp.Description = "启动成功"
	c.JSON(200,rsp)
	return
}




// MailboxConf 邮箱配置
type MailboxConf struct {
	// 邮件标题
	Title string
	// 邮件内容
	Body string
	// 收件人列表
	RecipientList []string
	// 发件人账号
	Sender string
	// 发件人密码，QQ邮箱这里配置授权码
	SPassword string
	// SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPAddr string
	// SMTP端口 QQ邮箱是25
	SMTPPort int
}


//邮箱发送逾期提醒邮件
func SendEmail()  {
	var mailConf MailboxConf
	mailConf.Title = "图书管理提醒"
	mailConf.Body = "借阅图书即将逾期，请及时续借或还书～"
	mailConf.RecipientList = []string{`2454546080@qq.com`}
	mailConf.Sender = `2454546080@qq.com`
	mailConf.SPassword = "plzlqoslwvmxdjci"
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 25

	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender)
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, mailConf.Body)
	m.Attach("/Users/link/Desktop/1.jpg")   //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		log.Fatalf("Send Email Fail, %s", err.Error())
		return
	}
	log.Printf("Send Email Success")
}

//发送逾期邮件，定时任务启动
func SendEmailTimer()  {

	//从用户开始寻找借阅逾期的记录
	var user model.User
	users,err:=user.GetAllUsers(data.Db)
	if err!=nil{
		fmt.Println("SendEmailTimer GetAllUsers err ",err)
	}
	for _,v:=range users{
		var borrow model.Borrow
		//发送邮件
	}
}