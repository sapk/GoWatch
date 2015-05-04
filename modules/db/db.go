package db

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Db struct {
	Orm *orm.Ormer
}

func InitDb() *Db {
	orm.RegisterDriver("sqlite3", orm.DR_Sqlite)
	orm.RegisterDataBase("default", "sqlite3", "gowatch.db")

	orm.RegisterModel(new(User))
	orm.Debug = true

	o := orm.NewOrm()
	o.Using("default") // Using default, you can use other database

	//Generate table if not exist
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println(err)
	}
	return &Db{Orm: &o}
}
