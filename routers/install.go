package routers

import (
	"log"

	"github.com/Unknwon/macaron"
	"github.com/sapk/GoWatch/models/user"
)

//Install deliver install page
func Install(ctx *macaron.Context) {
	if user.ContainMaster() {
		ctx.Redirect("/")
		return
	}
	ctx.HTML(200, "install")
}

//InstallPost handle installation
func InstallPost(ctx *macaron.Context) {
	if user.ContainMaster() {
		ctx.Redirect("/")
		return
	}
	log.Println("Installing ...")
	err := user.CreateUser(ctx.Query("username"), ctx.Query("password"), ctx.Query("email"), "master", []string{"master"})
	if err != nil {
		log.Println("Install failed !")
		ctx.Data["InstallError"] = true
		ctx.Data["InstallErrorText"] = err.Error()
		ctx.HTML(200, "install")
		return
	}

	ctx.Redirect("/")
}
