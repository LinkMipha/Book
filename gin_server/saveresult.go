package gin_server

import (
	"go-server/data"
	"go-server/model"
	"log"
)

type secKilMessage struct {
	username string
	coupon  model.Coupon
}

const  maxMessageNum = 20000

var secKillChannel = make(chan secKilMessage,maxMessageNum)

//异步存储
func Consumer()  {

	for{
		message := <-secKillChannel
		log.Println("get message: ",message)
		//用户优惠添加

		var err error
		coupon:=model.Coupon{}
		coupon.UserName = message.username
		coupon.CouponName = message.coupon.CouponName

		err = coupon.Insert(data.Db)//用户优惠券数+1
		if err != nil {
			println("Error when inserting user's coupon. " + err.Error())
		}

		//优惠券总库存减少
		coupon,err = coupon.GetCoupon(data.Db)
		ma:=make(map[string]interface{})
		ma["remain"] = coupon.Remain-1
		err = coupon.Update(data.Db,ma) //优惠券库存自减1
		if err != nil {
			println("Error when decreasing coupon left. " + err.Error())
		}

	}
}

var cousumerRun = false

func RunConsumer()  {
	if !cousumerRun{
		go Consumer()
		cousumerRun = true
	}
}
