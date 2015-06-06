package admin

import (
	"html/template"
	"log"
	"strconv"
	"strings"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/csrf"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
)

// Users generate the admin page for users management
func Users(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if err := auth.VerificationAuth(ctx, sess, []string{"admin.users"}); err != nil {
		return
	}
	fillGlobalPage(ctx, db, "admin_users")
	ctx.Data["users_count"], ctx.Data["Users"] = db.GetUsers()
	ctx.HTML(200, "admin/users")
}

// UserDel handle deletion of one user
func UserDel(ctx *macaron.Context, auth *auth.Auth, sess session.Store, dbb *db.Db, x csrf.CSRF) {
	if err := auth.VerificationAuth(ctx, sess, []string{"del.user"}); err != nil {
		return
	}
	id, _ := strconv.ParseUint(ctx.Params(":id"), 10, 64)
	user, err := dbb.GetUser(db.User{ID: id})
	if err != nil {
		ctx.Data["message_categorie"] = "negative"
		ctx.Data["message_icon"] = "user"
		ctx.Data["message_header"] = "User not found !"
		ctx.Data["message_text"] = "The user seems to not be in the database"
		ctx.Data["message_redirect"] = "/admin/users"
		ctx.HTML(200, "other/message")
	} else if user.ID == 1 {
		ctx.Data["message_categorie"] = "negative"
		ctx.Data["message_icon"] = "user"
		ctx.Data["message_header"] = "User is master !"
		ctx.Data["message_text"] = "You can't delete the master user."
		ctx.Data["message_redirect"] = "/admin/users"
		ctx.HTML(200, "other/message")
	} else if !csrf.ValidToken(ctx.Query("confirm"), "8e82e24bca448c990f69f5c364fc15ae", string(sess.Get("user").(db.User).ID), "del.user") {
		// Ask for confirmation
		ctx.Data["message_categorie"] = ""
		ctx.Data["message_icon"] = "user"
		ctx.Data["message_header"] = "Confirm user deletion!"
		ctx.Data["csrf_token"] = csrf.GenerateToken("8e82e24bca448c990f69f5c364fc15ae", string(sess.Get("user").(db.User).ID), "del.user")
		sess.Set("crsf_user_id", user.ID)
		ctx.Data["message_text"] = strings.Join([]string{"Do you really want to delete : ", user.Username}, " ")

		ctx.HTML(200, "other/confirmation")
	} else {
		// We del the user if all is good
		if sess.Get("crsf_user_id") != user.ID {
			ctx.Data["message_categorie"] = "negative"
			ctx.Data["message_icon"] = "user"
			ctx.Data["message_header"] = "Hummm ..."
			ctx.Data["message_text"] = template.HTML("It's seem there is a problem with the <a href='fr.wikipedia.org/wiki/Cross-Site_Request_Forgery' target='_blank'>CRSF</a> protection please retry.")
			ctx.Data["message_redirect"] = "/admin/users"
			ctx.HTML(200, "other/message")
			return
		}
		err := user.Delete()
		if err != nil {
			ctx.Data["message_categorie"] = ""
			ctx.Data["message_icon"] = "user"
			ctx.Data["message_header"] = "Oups"
			ctx.Data["message_text"] = "It's seem there is a problem during deletion"
			ctx.Data["message_redirect"] = "/admin/users"
			ctx.HTML(200, "other/message")
			return
		}
		ctx.Data["message_categorie"] = ""
		ctx.Data["message_icon"] = "user"
		ctx.Data["message_header"] = "User is deleted !"
		ctx.Data["message_text"] = "The user has been deleted from the database."
		ctx.Data["message_redirect"] = "/admin/users"
		ctx.HTML(200, "other/message")
		sess.Delete("crsf_user_id")
	}

}

// UserAdd generate the admin page for adding a user
func UserAdd(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if err := auth.VerificationAuth(ctx, sess, []string{"add.user"}); err != nil {
		return
	}
	fillGlobalPage(ctx, db, "admin_users")
	ctx.Data["roles"] = auth.GetRoles()
	ctx.HTML(200, "admin/add_user")
}

// UserAddPost handle the adding of a user
func UserAddPost(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if err := auth.VerificationAuth(ctx, sess, []string{"add.user"}); err != nil {
		return
	}
	err := db.CreateUser(ctx.Query("username"), ctx.Query("password"), ctx.Query("email"), ctx.Query("role"), auth.GetRoles())
	if err != nil {
		log.Println("User add failed !")
		fillGlobalPage(ctx, db, "admin_users")
		ctx.Data["roles"] = auth.GetRoles()
		ctx.Data["UserAddError"] = true
		ctx.Data["UserAddErrorText"] = err.Error()
		ctx.HTML(200, "admin/add_user")
		return
	}

	ctx.Data["message_categorie"] = "positive"
	ctx.Data["message_icon"] = "add user"
	ctx.Data["message_header"] = "User added !"
	ctx.Data["message_text"] = "The user has been added to the database and can login right now."
	ctx.Data["message_redirect"] = "/admin/users"
	ctx.HTML(200, "other/message")
}
