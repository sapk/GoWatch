package admin

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
)

func Dashboard(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if !auth.IsGranted("admin.dashboard", sess) {
		ctx.Data["message_categorie"] = "negative"
		ctx.Data["message_icon"] = "warning sign"
		ctx.Data["message_header"] = "Access forbidden"
		ctx.Data["message_text"] = "It's seem you don't have the right to be there"
		ctx.HTML(403, "other/message")
		return
	}
	ctx.Data["users_count"] = db.NbUsers()
	ctx.Data["admin_dashboard"] = true
	ctx.HTML(200, "admin/dashboard")
}
