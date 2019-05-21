package redisClient

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
	"log"
	"time"
	"github.com/astaxie/beego/logs"
)

var (
	redisC *redis.Pool
)

func init () {
	redisHost := beego.AppConfig.String("redis.host")
	redisPort := beego.AppConfig.String("redis.port")
	redisPass := beego.AppConfig.String("redis.pass")

	redisC = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     beego.AppConfig.DefaultInt("redis.maxidle", 10),
		MaxActive:   beego.AppConfig.DefaultInt("redis.maxactive", 1000),
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost+":"+redisPort , redis.DialPassword(redisPass))
			if err != nil {
				log.Println(err)
			}
			return c, err
		},
	}
}



func SetKey (key string , value string ) {
	rc := redisC.Get()
	defer rc.Close()
	_, err := rc.Do("set", key, value)
	if err != nil {
		fmt.Println(err)
	}
	logs.Debug("set key: %v value %v err %v", key , value , err)
}

func GetKey (key string) (re string) {
	rc := redisC.Get()
	defer rc.Close()
	value, err := redis.String(rc.Do("get", key))
	if err != nil {
		fmt.Println(err)
		return "nil"
	}
	return value
}

//func (c *CacheClient)GetKey (key string) (re string) {
//	rc := c.redisCli.Get()
//	defer rc.Close()
//	value, err := redis.String(rc.Do("get", key))
//	if err != nil {
//		fmt.Println(err)
//		return "nil"
//	}
//	return value
//}

func Setexpire (key string, period int){
	rc := redisC.Get()
	n, _ := rc.Do("EXPIRE", key, period)
	logs.Debug("set exporpe key %s period %v return %v",key ,period, n)
}