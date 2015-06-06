package admin

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
)

// Dashboard genrate the home admin page
func Dashboard(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if err := auth.VerificationAuth(ctx, sess, []string{"admin.dashboard"}); err != nil {
		return
	}
	fillGlobalPage(ctx, db, "admin_dashboard")
	ctx.HTML(200, "admin/dashboard")
}

func fillGlobalPage(ctx *macaron.Context, db *db.Db, page string) {
	ctx.Data["page_admin"] = true
	if page != "" {
		ctx.Data[page] = true
	}
	ctx.Data["users_count"] = db.NbUsers()
	ctx.Data["equipements_count"] = db.NbEquipements()
}
