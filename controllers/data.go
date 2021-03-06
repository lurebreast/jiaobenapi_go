package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"github.com/nfnt/resize"
	"image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"xiaozhang/models"
)

type DataController struct {
	BaseController
}

// @Title CreateUser
// @Description 上传数据
// @Param	Tid			post 	int 	true	"项目id"
// @Param	Orderid		post	int 	false	"序号 当有值时更新数据"
// @Param	Img			post 	string	false	"图片base64"
// @Param	Img1		post 	string	false	"图片1base64"
// @Param	File		post 	string	false	"文件base64"
// @Param	FileName	post 	string	false	"文件名字"
// @Param	Mobile		post	string	false	"手机号"
// @Param	Account		post	string	false	"用户名"
// @Param	Password	post 	string	false	"密码"
// @Param	Ip			post 	string	false	"ip"
// @Param	IpAttribution	post 	string 	false	"ip归属地"
// @Param	DeviceMode		post 	string 	false	"设备型号"
// @Param	DeviceVersion	post 	string 	false	"设备系统版本"
// @Param	Imsi			post	string	false	"imsi"
// @Param	SimId			post	string	false	"SimId"
// @Param	Name			post	string	false	"姓名"
// @Param	IdCard			post	string	false	"身份证号码"
// @Success 200 {object} models.Data
// @Failure 403 body is empty
// @router /post [post]
func (d *DataController) Post() {

	var isInsert = true
	Tid, _ := d.GetInt("Tid")
	if Tid == 0 {
		d.Error("项目id为空")
		return
	}
	o := orm.NewOrm()
	data := new(models.Data)

	var project models.Project
	err := o.QueryTable(new(models.Project)).Filter("Typeid", Tid).One(&project, "IsDelete")
	if err != nil {
		d.Error("项目id错误")
		return
	}

	if project.IsDelete == 1 {
		d.Error("此项目已经删除")
		return
	}

	OrderId, _ := d.GetInt("Orderid")
	if OrderId == 0 {
		OrderId = getOrderId(Tid, true)
	} else {
		isInsert = false
		var data1 models.Data

		err1 := orm.NewOrm().QueryTable(data).Filter("Tid", Tid).Filter("Orderid", OrderId).One(&data1)
		if err1 == nil {
			data = &data1
		}
	}

	data.Tid = Tid
	data.Orderid = OrderId
	if s := d.GetString("Mobile"); s != "" {
		data.Mobile = s
	}
	if s := d.GetString("Account"); s != "" {
		data.Account = s
	}
	if s := d.GetString("Password"); s != "" {
		data.Password = s
	}
	if s := d.GetString("Ip"); s != "" {
		data.Ip = s
	}
	if s := d.GetString("IpAttribution"); s != "" {
		data.IpAttribution = s
	}
	if s := d.GetString("Imei"); s != "" {
		data.Imei = s
	}
	if s := d.GetString("DeviceMode"); s != "" {
		data.DeviceMode = s
	}
	if s := d.GetString("DeviceVersion"); s != "" {
		data.DeviceVersion = s
	}
	if s := d.GetString("Imsi"); s != "" {
		data.Imsi = s
	}
	if s := d.GetString("SimId"); s != "" {
		data.SimId = s
	}
	if s := d.GetString("Name"); s != "" {
		data.Name = s
	}
	if s := d.GetString("IdCard"); s != "" {
		data.IdCard = s
	}

	t := time.Now()
	Timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	Timestamp = Timestamp[:10]

	data.Updatetime = Timestamp

	if img := d.GetString("Img"); img != "" {
		data.Img = saveImg(img, Tid, OrderId, "0")
	}
	if img1 := d.GetString("Img1"); img1 != "" {
		data.Img1 = saveImg(img1, Tid, OrderId, "1")
	}
	if file := d.GetString("File"); file != "" {
		if fileName := d.GetString("FileName"); fileName != "" {
			data.File = saveFile(file, fileName, Tid, OrderId)
		}
	}

	if isInsert {
		data.Creattime = Timestamp
		data.Status = 1
		_, err := o.Insert(data)
		if err != nil {
			beego.Error(err)
			d.Error(err.Error())
			return
		}
	} else {
		_, err := o.Update(data)
		if err != nil {
			beego.Error(err)
			d.Error(err.Error())
			return
		}
	}

	d.success(data)
}

// @Title 获取条数
// @Description 获取条数
// @Param	Tid		query 	int		true	"项目id"
// @Param	Day		query 	string	false	"获取多少天前零点到现在的数据"
// @Param	Status	query 	string	false	"是否获取已提取 1 已提取"
// @Success 200 {object} Response
// @Failure 403 body is empty
// @router /count [get]
func (d *DataController) Count() {
	var Total string
	var Typename string
	var result map[string]string

	Tid := d.GetString("Tid")
	Day, _ := d.GetInt("Day")
	Status := d.GetString("Status")

	if Tid == "" {
		d.Error("没有项目ID")
		return
	}

	rc := RedisClient.Get();
	defer rc.Close();

	key := "typeid_count2_" + Tid + "_" + Status + strconv.Itoa(Day)
	if  ok, _ := redis.Bool(rc.Do("setnx", key + "_lock", "1")); ok {
		rc.Do("expire", key + "_lock", 60)

		qs := orm.NewOrm().QueryTable(new(models.Data)).Filter("Tid", Tid)

		if Status == "1" {
			qs = qs.Filter("Status", 2)
		}
		if Day != 0 {
			timeStr := time.Now().Format("2006-01-02")
			t2, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
			updateTimestamp := t2.AddDate(0, 0, -Day).Unix()
			qs = qs.Filter("Updatetime__gt", updateTimestamp - 1)
		}

		t, err := qs.Count()
		if err != nil {
			beego.Error(err.Error())
		} else {
			Total = strconv.FormatInt(t, 10)
		}

		var project models.Project
		err = orm.NewOrm().QueryTable(new(models.Project)).Filter("Typeid", Tid).One(&project, "Typename")
		if err != nil {
			beego.Error(err.Error())
		} else {
			Typename = project.Typename
		}

		rc.Do("hset", key, "Total", Total)
		rc.Do("hset", key, "Typename", Typename)

		result = make(map[string]string)
		result["Total"] = Total
		result["Typename"] = Typename
	} else {
		result, _ = redis.StringMap(rc.Do("hgetall", key))
	}

	d.success(result)
}

// @Title 获取单条数据
// @Description 获取单条数据
// @Param	Tid		query 	int		true	"项目id"
// @Param	Orderid		query 	string	false	"数据id"
// @Param	Rand	query 	string	false	"是否随机 1 是"
// @Param	Update	query 	string	false	"是否更新为已提取 1 更新"
// @Success 200 {object} models.Data
// @Failure 403 body is empty
// @router /getone [get]
func (d *DataController) Getone() {

	Tid := d.GetString("Tid")
	Orderid := d.GetString("Orderid")
	Rand := d.GetString("Rand")
	Update := d.GetString("Update")

	if Tid == "" {
		d.Error("没有项目ID")
		return
	}

	TidInt, _ := strconv.Atoi(Tid)

	if getOrderId(TidInt, false) == 0 { // 通过序号检测项目是否存在
		d.Error("项目id错误")
		return
	}

	o := orm.NewOrm()

	rc := RedisClient.Get();
	defer rc.Close();

	TidKey := "tid_orderid_" + Tid
	if Update == "1" {
		if Rand == "1" {
			num := rand.Intn(2)
			if num == 1 {
				Orderid, _ = redis.String(rc.Do("RPOP", TidKey))
			} else {
				Orderid, _ = redis.String(rc.Do("LPOP", TidKey))
			}
		} else {
			if Orderid == "" {
				Orderid, _ = redis.String(rc.Do("RPOP", TidKey))
			}
		}

		if Orderid != "" {
			var data1 models.Data
			err := o.QueryTable(new(models.Data)).Filter("Tid", Tid).Filter("Orderid", Orderid).One(&data1)
			if err != nil {
				beego.Error(err.Error())
			} else {
				data := new(models.Data)

				data = &data1
				data.Status = 2
				o.Update(data, "Status")

				d.success(data)
				return
			}
		}
	} else {
		OrderidInt := getOrderId(TidInt, false)

		if Rand == "1" {
			//MinOrderId := getMinOrderId(TidInt)
			//Orderid1 := rand.Intn(OrderidInt - 1 - MinOrderId) + 1 + MinOrderId
			Orderid1 := rand.Intn(OrderidInt - 1) + 1
			Orderid = strconv.Itoa(Orderid1)
		} else {
			if Orderid == "" {
				Orderid = strconv.Itoa(OrderidInt)
			}
		}

		var data1 models.Data
		err := o.QueryTable(new(models.Data)).Filter("Tid", Tid).Filter("Orderid", Orderid).One(&data1)
		if err != nil {
			beego.Error(err.Error())
		} else {
			d.success(data1)
			return
		}
	}

	d.Error("没有可用数据")
}


func saveImg(s string, Tid int, Orderid int, imgId string) string {
	s = strings.Replace(s, "data:image/png;base64,", "", -1)
	s = strings.Replace(s, " ", "+", -1)
	decode, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		beego.Info(s)
		beego.Error(err.Error())
		return ""
	}

	img := "/images/" + strconv.Itoa(Tid) + "_" + strconv.Itoa(Orderid) + "_" + imgId + ".png"
	root_path := "/home/wwwroot/default/public"
	err = ioutil.WriteFile(root_path + img, decode, 0666)
	if err != nil {
		beego.Error(err.Error())
		return ""
	}

	file, err := os.Open(root_path + img)
	defer file.Close()
	if err != nil {
		beego.Error(err.Error())
		return ""
	}

	img_source, err := png.Decode(file)
	if err != nil {
		beego.Error(err.Error())
		return ""
	}

	b := img_source.Bounds()
	width := b.Max.X
	height := b.Max.Y
	width = int(math.Ceil(float64(width) * 0.6))
	height = int(math.Ceil(float64(height) * 0.6))
	w := uint(width);
	h := uint(height)

	m := resize.Resize(w, h, img_source, resize.NearestNeighbor)
	out, err := os.Create(root_path + img)
	defer out.Close()
	if err != nil {
		beego.Error(err.Error())
		return ""
	}

	png.Encode(out, m)

	return img
}

func saveFile(s, fileName string, Tid, Orderid int) string {
	s = strings.Replace(s, " ", "+", -1)
	decode, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		beego.Info(s)
		beego.Error(err.Error())
		return ""
	}

	path := "/images/" + strconv.Itoa(Tid) + "_" + strconv.Itoa(Orderid) + "_3" + fileName
	root_path := "/home/wwwroot/default/public"
	err = ioutil.WriteFile(root_path + path, decode, 0666)
	if err != nil {
		beego.Error(err.Error())
		return ""
	}

	return path
}

func getOrderId(t int, isIncr bool) int {
	typeid := strconv.Itoa(t)
	r := RedisClient.Get()
	defer r.Close()

	var typedataid = "0"
	key := "increment_order_id_" + typeid + "_2"

	if exists, _ := redis.Bool(r.Do("exists", key)); !exists {
		var data models.Data
		err  := orm.NewOrm().QueryTable(new(models.Data)).Filter("Tid", typeid).OrderBy("-Orderid").Limit(1).One(&data, "Orderid")
		if err == nil {
			r.Do("set", key, data.Orderid)
		}
	}

	if isIncr {
		incr, _ := redis.Int64(r.Do("incr", key))
		r.Do("rPush", "tid_orderid_" + typeid, incr)
		typedataid = strconv.FormatInt(incr, 10)
	} else {
		typedataid, _ = redis.String(r.Do("get", key))
	}

	orderid, _ :=  strconv.Atoi(typedataid)
	return orderid
}

// 获取tid最orderid
//func getMinOrderId(tid int) int {
//	typeid := strconv.Itoa(tid)
//	r := RedisClient.Get()
//	defer r.Close()
//
//	key := "min_orderid_" + typeid
//
//	if exists, _ := redis.Bool(r.Do("exists", key + "_lock")); !exists {
//		r.Do("setex", key + "_lock", 60, 1)
//
//		var data models.Data
//		err  := orm.NewOrm().QueryTable(new(models.Data)).Filter("Tid", typeid).OrderBy("Orderid").Limit(1).One(&data, "Orderid")
//		if err == nil {
//			r.Do("set", key, data.Orderid)
//		}
//	}
//
//	typedataid, _ := redis.Int(r.Do("get", key))
//	return typedataid
//}
