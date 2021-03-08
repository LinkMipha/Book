package data

import (
	"context"
	"github.com/go-redis/redis"
	"go-server/conf"
	"log"
)
var Client *redis.Client


func initRedis(config conf.Config){
	Client = redis.NewClient(&redis.Options{
		Addr: config.Redis.Address,
		Password: config.Redis.Password,
		DB: 0,
	})

	_,err:=FlushAll()
	if err!=nil{
		log.Fatal("flushAll",err.Error())
	}

}

func FlushAll()(string,error)  {
	return Client.FlushAll(context.TODO()).Result()
}

//加载lua脚本
func LoadScript(script string)string  {
	scriptExists,err:= Client.ScriptExists(context.TODO(),script).Result()
	if err!=nil{
		panic("exist failed err: %v"+err.Error())
	}
	if !scriptExists[0]{
		scriptSHA,err:=Client.ScriptLoad(context.TODO(),script).Result()
		if err!=nil{
			panic("load script error"+err.Error())
		}
		return scriptSHA
	}
	log.Println("script exists")
	return ""
}

func EvalSHA(SHA string,args[]string)(interface{},error)  {
	val,err:=Client.Eval(context.TODO(),SHA,args).Result()
	if err!=nil{
		log.Println("eval failed err",err.Error())
		return nil, err
	}
	return val,err
}

//set time forever
func SetTime(key string,values interface{})(string,error)  {
	val,err:=Client.Set(context.TODO(),key,values,0).Result()
	return val,err
}

//设置hash值 返回bool值
func SetHash(key string,field map[string]interface{})(bool,error) {
	return Client.HMSet(context.TODO(),key,field).Result()
}

//获取hash表多个字段值hget 只能获取一个字段
func GetMap(key string,fields ...string) ([]interface{},error) {
	return Client.HMGet(context.TODO(),key,fields...).Result()
}


func SetAdd(key string,field string)(int64,error)  {
	return  Client.SAdd(context.TODO(),key,field).Result()
}

//判断是否是集合的值
func SetIsMember(key string,field string)(bool,error)  {
	return Client.SIsMember(context.TODO(),key,field).Result()
}

//获取集合所有成员
func GetMembers(key string)([]string,error)  {
	return Client.SMembers(context.TODO(),key).Result()
}


