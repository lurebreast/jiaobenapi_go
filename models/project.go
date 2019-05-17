package models

import "github.com/astaxie/beego/orm"

type Project struct {
	Typeid int `orm:"pk"`
	Typename string
	Createtime string
	Updatetime string
	IsDelete uint8
}

func (d *Project) TableName() string {
	return "type"
}

func init() {
	orm.RegisterModel(new(Project))
}