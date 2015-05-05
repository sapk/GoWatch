// auth
package auth

import (
	"encoding/json"
	"github.com/Unknwon/macaron"
	"github.com/astaxie/beego/orm"
	"github.com/macaron-contrib/session"
	"github.com/mikespook/gorbac"
	"github.com/sapk/GoWatch/modules/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	//"net/http"
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
	rbac *gorbac.Rbac
}

func Authentificator(options ...Options) macaron.Handler {
	opt := prepareOptions(options)
	auth := New(opt.Provider)
	return func(ctx *macaron.Context, sess session.Store) {
		ctx.Map(auth)
	}
}
func New(db *db.Db) *Auth {
	return &Auth{
		db:   db,
		rbac: initRbac(db),
	}
}

func initRbac(db *db.Db) *gorbac.Rbac {
	//Roles //TODO
	rbac := gorbac.New()
	rbac.Set("user", []string{"open.equipement"}, nil)
	rbac.Add("user", []string{"open.dashboard"}, nil)
	rbac.Set("admin", []string{"add.equipement", "del.equipement", "add.user", "del.user", "admin.dashboard"}, []string{"user"})
	rbac.Set("master", []string{}, []string{"admin"})
	return rbac
}

//TODO
func SignIn(ctx *macaron.Context, sess session.Store, auth *Auth) {

	o := auth.db.Orm
	user := db.User{Username: ctx.Query("username")}
	err := (*o).Read(&user, "username")

	if err == orm.ErrNoRows {
		log.Printf("User %s not found.", ctx.Query("username"))
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found in user table.")
	} else {
		//fmt.Println(user.Id, user.Username)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ctx.Query("password")))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Println("Bad password with username %s.", user.Username)
		} else if err != nil {
			log.Println("Unkwon error")
		} else {
			log.Printf("User %s login succesfully.", user.Username)
			sess.Set("uid", user.Id)
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

//TODO
func IsLogged(ctx *macaron.Context, sess session.Store, auth *Auth) {
	if sess.Get("uid") == nil {
		sess.Set("auth.redirect_after_login", ctx.Req.RequestURI)
		ctx.Redirect("/user/login")
		return
	}
	log.Printf("rbac : %v", auth.rbac.Get("master"))
	//TODO
	//	auth.rbac.IsGranted(sess.Get("user").(string), "page.article", nil)
	ctx.Data["user"] = sess.Get("user").(db.User)
	ctx.Next()
}

//TODO
func (this *Auth) IsGranted(action string, sess session.Store) bool {
	return IsGranted(action, sess, this)
}
func IsGranted(action string, sess session.Store, auth *Auth) bool {
	var roles []interface{}
	json.Unmarshal([]byte(sess.Get("user").(db.User).Roles), roles)
	//
	for _, role := range roles {
		if auth.rbac.IsGranted(role.(string), action, nil) {
			return true
		}
	}
	return false
}
