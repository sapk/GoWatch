package db

import (
	"fmt"
	"log"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3" // import your used driver
)

//Db represent the database
type Db struct {
	Orm *orm.Ormer
}

var db Db
var isInit = false

//Get get ref to db
func Get() *Db {
	return &db
}

//init init the database
func init() {
	if !isInit {
		log.Println("Initialisation of Database...")
		orm.RegisterDriver("sqlite3", orm.DR_Sqlite)
		//orm.RegisterDataBase("default", "sqlite3", "gowatch.db")
		orm.RegisterDataBase("default", "sqlite3", "gowatch.db?cache=shared&mode=memory")
		//orm.RegisterDataBase("default", "sqlite3", ":memory:")

		orm.RegisterModel(new(User))
		orm.RegisterModel(new(Equipement))
		orm.Debug = true

		o := orm.NewOrm()
		o.Using("default") // Using default, you can use other database

		//Generate table if not exist
		err := orm.RunSyncdb("default", false, true)
		if err != nil {
			fmt.Println(err)
		}
		db = Db{Orm: &o}
		isInit = true
	}
}
