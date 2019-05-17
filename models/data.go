package models

import "github.com/astaxie/beego/orm"

type Data struct {
	Id   int  //`orm:"pk"`
	//Data string
	Tid int
	Orderid int
	Status int
	Creattime string
	Updatetime string
	Img string
	Img1 string
	Mobile string
	Account string
	Password string
	Ip string
	IpAttribution string
	Imei string
	DeviceMode string
	DeviceVersion string
	Imsi string
	ImsiId string
	Name string
	IdCard string
}

func (d *Data) TableName() string {
	return "typedata"
}

func init() {
	orm.RegisterModel(new(Data))
}

