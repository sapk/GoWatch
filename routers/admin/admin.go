package admin

import (
	"errors"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
)

// Dashboard genrate the home admin page
func Dashboard(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if err := verificationAuth(ctx, auth, sess); err != nil {
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

func verificationAuth(ctx *macaron.Context, auth *auth.Auth, sess session.Store) error {
	if !auth.IsGranted("admin.users", sess) {
		ctx.Data["message_categorie"] = "negative"
		ctx.Data["message_icon"] = "warning sign"
		ctx.Data["message_header"] = "Access forbidden"
		ctx.Data["message_text"] = "It's seem you don't have the right to be there"
		ctx.HTML(403, "other/message")
		return errors.New("Not allowed")
	}
	return nil
}
