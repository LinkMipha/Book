package data

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"time"
)
var (
	config clientv3.Config
	client *clientv3.Client
	err error
)

func init()  {
	// 客户端配置
	config = clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
}

// 设置key值
func SetKey(key,value string) {
	ctx,cancel := context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()
	_, err = client.Put(ctx, key, value)
}

// 获取key值
func GetKey(key string)  {
	// 实例化一个用于操作ETCD的KV
	kv := clientv3.NewKV(client)
	if getResp, err := kv.Get(context.TODO(), "link"); err != nil {
		fmt.Println(err)
		return
	}else{
		// 输出本次的Revision
		fmt.Printf("Key is s %s \n Value is %s \n", getResp.Kvs[0].Key, getResp.Kvs[0].Value)
	}
}

// 监听key值的滨化
func WatchKey(key string)  {
	fmt.Println("\n...watch demo...")
	stopChan := make(chan interface{}) // 是否停止信号
	go func() {
		watchChan := client.Watch(context.TODO(), key, clientv3.WithPrefix())
		for {
			select {
			case result := <- watchChan:
				for _, event := range result.Events {
					fmt.Printf("%s %q : %q\n", event.Type, event.Kv.Key, event.Kv.Value)
				}
			case <-stopChan:
				fmt.Println("stop watching...")
				return
			}
		}
	}()

	//多次更改其中的值
	for i := 0; i < 5; i++ {
		value:=strconv.Itoa(i)
		client.Put(context.TODO(),key,value)
	}

	time.Sleep(time.Second * 1)
	stopChan <- 1 //停止watch，在插入就不会监听到了
}

func Close()  {
	if nil != client {
		client.Close()
	}
	log.Println("Stop")
}