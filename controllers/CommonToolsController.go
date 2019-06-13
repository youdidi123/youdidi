package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"youdidi/redisClient"
	"github.com/astaxie/beego/logs"
)

type CommonToolsController struct {
	beego.Controller
}

var (
	UserLockPrefix = "USER_LOCK_"
	OrderLockPrefix = "ORDER_LOCK_"
	OrderDetailLockPrefix = "OD_LOCK_"
	LockExpireTime = 2 * 60 //预计不会有哪个锁需要持续2min以上
)

var (
	//获取锁等待20s，如果获取不到，前端直接返回超时
	retryNum = 20
	retryTime = 1
)

func SetOrderLock (oid string) bool {
	i := 0
	for {
		if redisClient.Lock(OrderLockPrefix + oid) {
			redisClient.Setexpire(OrderLockPrefix + oid , LockExpireTime)
			logs.Debug("set order lock success oid=%v" , oid)
			return true
		}
		if (i >= retryNum) {
			logs.Emergency("FATAL:set order lock fail oid=%v" , oid)
			return false
		}
		i++
		time.Sleep(time.Duration(retryTime)*time.Second)
	}
}

func DelOrderLock (oid string) {
	i := 0
	for {
		if redisClient.UnLock(OrderLockPrefix + oid) {
			logs.Debug("unset order lock success oid=%v" , oid)
			break
		}
		if (i >= retryNum) {
			logs.Emergency("FATAL:unset order lock fail oid=%v" , oid)
			break
		}
		i++
		time.Sleep(time.Duration(retryTime)*time.Second)
	}
}

//sig:md5(所有字母必须大写) auth:base64
//短信验证码平台鉴权使用
func getSig (id string , token string) (string , string){
	ltime := time.Now().Format("20060102150405")
	fmt.Println(ltime)

	sig := md5.New()
	sig.Write([]byte(id+token+ltime))

	auth := base64.StdEncoding.EncodeToString([]byte(id+":"+ltime))

	return strings.ToUpper(hex.EncodeToString(sig.Sum(nil))),auth
}

//公共函数，获取一个以当前时间为sed的6位随机数
func GetRandomCode () string{
	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
	min := 1000
	code := r1.Intn(10000)
	if (code < min) {
		code += min
	}
	return strconv.Itoa(code)
}