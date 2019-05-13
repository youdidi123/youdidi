package redisClient

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
	"log"
	"time"
)

type CacheClient struct {
	redisCli   *redis.Pool
	redis_host string
	redis_port string
}

func (c *CacheClient)GetConnet() {
	fmt.Println("here to connet redis")
	// 从配置文件获取redis的ip以及db
	c.redis_host = beego.AppConfig.String("redis.host")
	c.redis_port = beego.AppConfig.String("redis.port")
	// 建立连接池
	c.redisCli = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     beego.AppConfig.DefaultInt("redis.maxidle", 1),
		MaxActive:   beego.AppConfig.DefaultInt("redis.maxactive", 10),
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", c.redis_host+":"+c.redis_port)
			if err != nil {
				log.Println(err)
			}
			return c, err
		},
	}
}

func (c *CacheClient)SetKey (key string , value string ) {
	rc := c.redisCli.Get()
	defer rc.Close()
	_, err := rc.Do("set", key, value)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *CacheClient)GetKey (key string) (re string) {
	rc := c.redisCli.Get()
	defer rc.Close()
	value, err := redis.String(rc.Do("get", key))
	if err != nil {
		fmt.Println(err)
		return "nil"
	}
	return value
}