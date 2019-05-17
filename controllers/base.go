package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisClient *redis.Pool

func init() {
	RedisClient = &redis.Pool{
		MaxIdle :  5,
		MaxActive : 10,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", beego.AppConfig.String("redisdsn"))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}
type BaseController struct {
	beego.Controller
}

type Response struct {
	Code int
	Msg string
	Data interface{}
}

func (b *BaseController) success(data interface{}) {
	r := &Response{Code:200, Msg:"", Data:data}
	b.Data["json"] = r
	b.ServeJSON()
}

func (b *BaseController) Error(err string)  {
	r := &Response{Code:500, Msg:err, Data:""}
	b.Data["json"] = r
	b.ServeJSON()
}
