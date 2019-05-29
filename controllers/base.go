package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/gomodule/redigo/redis"
	"net/http"
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

	beego.ErrorHandler("404", Error404)
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

func Error404(rw http.ResponseWriter, r *http.Request){
	data := &Response{Code:404, Msg:"Page not found", Data:""}
	content, _ := json.Marshal(data)
	_, err := rw.Write(content)
	if err != nil {
		beego.Error(err.Error())
	}
}