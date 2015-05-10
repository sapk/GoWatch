package admin

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/csrf"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
)

// Equipements generate the admin page for Equipement management
func Equipements(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {

	if err := verificationAuth(ctx, auth, sess); err != nil {
		return
	}
	fillGlobalPage(ctx, db, "admin_equipements")
	ctx.Data["users_count"], ctx.Data["Users"] = db.GetEquipements()
	ctx.HTML(200, "admin/equipements")
	//TODO representation in tmeplate
}

// EquipementDel handle deletion of one user
func EquipementDel(ctx *macaron.Context, auth *auth.Auth, sess session.Store, dbb *db.Db, x csrf.CSRF) {
	//TODO
}

// EquipementAdd generate the admin page for adding a user
func EquipementAdd(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if err := verificationAuth(ctx, auth, sess); err != nil {
		return
	}
	fillGlobalPage(ctx, db, "admin_equipements")
	//TODO ?	ctx.Data["organizations"] = auth.GetOrganizations()
	//TODO ?	ctx.Data["locations"] = auth.GetLocations()
	ctx.Data["types"] = db.GetEquipementTypes()
	ctx.HTML(200, "admin/add_equipement")
}

// EquipementAddPost handle the adding of a user
func EquipementAddPost(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	//TODO add support for adding wildcard dns (find host in DNS) && support ip scan of range
	//TODO convert hostname to IP
	//TODO use config for DNS resolver
	//TODO
}
