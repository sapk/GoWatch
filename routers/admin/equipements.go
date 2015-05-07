package admin

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
)

func Equipements(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {

	if err := verificationAuth(ctx, auth, sess); err != nil {
		return
	}
	fillGlobal(ctx, db)
	ctx.Data["admin_equipements"] = true
	ctx.HTML(200, "admin/equipements")
}
