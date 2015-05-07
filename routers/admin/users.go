package admin

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
)

func Users(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if !auth.IsGranted("admin.users", sess) {
		ctx.Data["message_categorie"] = "negative"
		ctx.Data["message_icon"] = "warning sign"
		ctx.Data["message_header"] = "Access forbidden"
		ctx.Data["message_text"] = "It's seem you don't have the right to be there"
		ctx.HTML(403, "other/message")
		return
	}
	ctx.Data["users_count"], ctx.Data["Users"] = db.GetUsers()
	//ctx.Data["users_count"] = db.NbUsers()
	ctx.Data["admin_users"] = true
	ctx.HTML(200, "admin/users")
}

func UserAdd(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if !auth.IsGranted("admin.users", sess) {
		ctx.Data["message_categorie"] = "negative"
		ctx.Data["message_icon"] = "warning sign"
		ctx.Data["message_header"] = "Access forbidden"
		ctx.Data["message_text"] = "It's seem you don't have the right to be there"
		ctx.HTML(403, "other/message")
		return
	}
	//ctx.Data["users_count"], ctx.Data["Users"] = db.GetUsers()
	ctx.Data["users_count"] = db.NbUsers()
	ctx.Data["roles"] = auth.GetRoles()
	ctx.Data["admin_users"] = true
	ctx.HTML(200, "admin/add_user")
}
