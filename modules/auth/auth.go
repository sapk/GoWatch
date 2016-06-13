// auth
package auth

import (
	//"encoding/json"
	"errors"
	"log"

	"github.com/Unknwon/macaron"
	"github.com/astaxie/beego/orm"
	"github.com/macaron-contrib/session"
	"github.com/mikespook/gorbac"
	"github.com/sapk/GoWatch/models/user"
	"github.com/sapk/GoWatch/modules/db"
	"golang.org/x/crypto/bcrypt"
)

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}
	//TODO default value and error
	return opt
}

// Options represents a struct for specifying configuration options for the session middleware.
type Options struct {
	// Name of provider for the rules.
	Provider *db.Db
}
type Auth struct {
	db   *db.Db
	rbac *gorbac.RBAC
}

func Authentificator(options ...Options) macaron.Handler {
	opt := prepareOptions(options)
	auth := New(opt.Provider)
	return func(ctx *macaron.Context, sess session.Store) {
		ctx.Map(auth)
		//log.Printf("Uri : %v", ctx.Req.RequestURI != "/install")
		//log.Printf("Containmaster : %v", (*auth.db).ContainMaster())
		//TODO used a var in config
		if !(user.ContainMaster()) && ctx.Req.RequestURI != "/install" {
			ctx.Redirect("/install")
			return
			//(*auth.db).CreateUser("master", "0000", "master", "master@localhost")
		}
		if sess.Get("uid") != nil {
			ctx.Data["user"] = sess.Get("user").(db.User)
			ctx.Data["role"] = auth.rbac.Get(sess.Get("user").(db.User).Roles)
		}
	}
}
func New(db *db.Db) *Auth {
	return &Auth{
		db:   db,
		rbac: initRbac(),
	}
}

func LogOut(ctx *macaron.Context, sess session.Store) {
	sess.Flush()
	ctx.Redirect("/")
}
func SignIn(ctx *macaron.Context, sess session.Store, auth *Auth) {

	o := auth.db.Orm
	user := db.User{Username: ctx.Query("username")}
	err := (*o).Read(&user, "username")

	if err == orm.ErrNoRows {
		log.Printf("User %s not found.", ctx.Query("username"))
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found in user table.")
	} else {
		//fmt.Println(user.ID, user.Username)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ctx.Query("password")))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Println("Bad password with username %s.", user.Username)
		} else if err != nil {
			log.Println("Unkwon error")
		} else {
			log.Printf("User %s login succesfully.", user.Username)
			sess.Set("uid", user.ID)
			sess.Set("user", user)
			//sess.Set("roles", user.Roles)
			if sess.Get("auth.redirect_after_login") == nil {
				ctx.Redirect("/")
			} else {
				ctx.Redirect(sess.Get("auth.redirect_after_login").(string))
			}
			return
		}
	}

	ctx.Data["AuthLoginError"] = true
	ctx.HTML(200, "user/login")
}

func IsLogged(ctx *macaron.Context, sess session.Store, auth *Auth) {
	if sess.Get("uid") == nil {
		sess.Set("auth.redirect_after_login", ctx.Req.RequestURI)
		ctx.Redirect("/user/login")
		return
	}
	ctx.Next()
}

func (this *Auth) IsGranted(action string, sess session.Store) bool {
	return IsGranted(action, sess, this)
}

func IsGranted(action string, sess session.Store, auth *Auth) bool {
	/*
		var roles []interface{}
		json.Unmarshal([]byte(sess.Get("user").(db.User).Roles), roles)
		//
		for _, role := range roles {
			if auth.rbac.IsGranted(role.(string), action, nil) {
				return true
			}
		}
	*/
	//First only role allowed
	return auth.rbac.IsGranted(sess.Get("user").(db.User).Roles, action, nil)
}

//VerificationAuth verify if the needed are filled
func (auth *Auth) VerificationAuth(ctx *macaron.Context, sess session.Store, needed []string) error {
	// if we need to verify that it's a admin : admin.users

	// We check if all the other need are filled
	for _, need := range needed {
		if !auth.IsGranted(need, sess) {
			ctx.Data["message_categorie"] = "negative"
			ctx.Data["message_icon"] = "warning sign"
			ctx.Data["message_header"] = "Access forbidden"
			ctx.Data["message_text"] = "It's seem you don't have the right to be there"
			ctx.HTML(403, "other/message")
			return errors.New("Not allowed")
		}
	}
	return nil
}
