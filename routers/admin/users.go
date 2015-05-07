package admin

import (
	"log"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
)

func verificationAuth(ctx *macaron.Context, auth *auth.Auth, sess session.Store) {
	ctx.Data["admin_users"] = true
	if !auth.IsGranted("admin.users", sess) {
		ctx.Data["message_categorie"] = "negative"
		ctx.Data["message_icon"] = "warning sign"
		ctx.Data["message_header"] = "Access forbidden"
		ctx.Data["message_text"] = "It's seem you don't have the right to be there"
		ctx.HTML(403, "other/message")
	}
}

// Users generate the admin page for users management
func Users(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	verificationAuth(ctx, auth, sess)
	ctx.Data["users_count"], ctx.Data["Users"] = db.GetUsers()
	ctx.HTML(200, "admin/users")
}

// UserAdd generate the admin page for adding a user
func UserAdd(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	verificationAuth(ctx, auth, sess)
	ctx.Data["users_count"] = db.NbUsers()
	ctx.Data["roles"] = auth.GetRoles()
	log.Printf("%v", auth.GetRoles())
	ctx.HTML(200, "admin/add_user")
}

// UserAddPost handle the adding of a user
func UserAddPost(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	verificationAuth(ctx, auth, sess)
	err := db.CreateUser(ctx.Query("username"), ctx.Query("password"), ctx.Query("email"), ctx.Query("role"), auth.GetRoles())
	if err != nil {
		log.Println("User add failed !")
		ctx.Data["users_count"] = db.NbUsers()
		ctx.Data["roles"] = auth.GetRoles()
		ctx.Data["UserAddError"] = true
		ctx.Data["UserAddErrorText"] = err.Error()
		ctx.HTML(200, "admin/add_user")
		return
	}
	//TODO verify and add the user to the db
	ctx.Data["message_categorie"] = "positive"
	ctx.Data["message_icon"] = "add user"
	ctx.Data["message_header"] = "User added !"
	ctx.Data["message_text"] = "The user has been added to the database and can login right now."
	ctx.Data["message_redirect"] = "/admin/users"
	ctx.HTML(403, "other/message")
}
