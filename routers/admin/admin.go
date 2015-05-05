package admin

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
)

func Dashboard(ctx *macaron.Context, auth *auth.Auth, sess session.Store) {
	//TODO
	if !auth.IsGranted("admin.dash", sess) {
		ctx.HTML(500, "admin/dashboard")
	}
	ctx.Data["admin_dashboard"] = true
	ctx.HTML(200, "admin/dashboard")
}
