// auth
package auth

import (
	"./../db"
	"github.com/Unknwon/macaron"
	"github.com/astaxie/beego/orm"
	"github.com/macaron-contrib/session"
	"github.com/mikespook/gorbac"
	"golang.org/x/crypto/bcrypt"
	"log"
	//"net/http"
)

type User db.User

type Auth struct {
	db   *db.Db
	rbac *gorbac.Rbac
}

func New(db *db.Db) *Auth {
	//Roles //TODO
	rbac := gorbac.New()
	rbac.Set("user", []string{"open.equipement"}, nil)
	rbac.Add("user", []string{"open.dashboard"}, nil)
	rbac.Set("admin", []string{"add.equipement", "del.equipement"}, []string{"user"})
	rbac.Set("master", []string{}, []string{"admin"})
	return &Auth{
		db:   db,
		rbac: rbac,
	}
}

//TODO
func (this *Auth) SignIn(ctx *macaron.Context, sess session.Store) {

	o := this.db.Orm
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
func (this *Auth) IsLogged(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") == nil {
		sess.Set("auth.redirect_after_login", ctx.Req.RequestURI)
		ctx.Redirect("/user/login")
		return
	}
	log.Printf("rbac : %v", this.rbac.Get("master"))
	//TODO
	rbac.IsGranted(sess.Get("user"), "page.article", nil)
	ctx.Data["user"] = sess.Get("user")
	ctx.Next()
}

//TODO
func (this *User) IsGranted(action string) bool {
	return false
}
