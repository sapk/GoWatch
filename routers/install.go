package routers

import (
	"github.com/Unknwon/macaron"
	"github.com/sapk/GoWatch/modules/db"
	"log"
)

func Install(ctx *macaron.Context) {
	ctx.HTML(200, "install")
}

func InstallPost(ctx *macaron.Context, db *db.Db) {
	log.Println("Installing ...")
	err := db.CreateUser(ctx.Query("username"), ctx.Query("password"), ctx.Query("email"), "master")
	if err != nil {
		log.Println("Install failed !")
		ctx.Data["InstallError"] = true
		ctx.Data["InstallErrorText"] = err.Error()
		ctx.HTML(200, "install")
		return
	}

	ctx.Redirect("/")
}
