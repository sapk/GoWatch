package web

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
	//"golang.org/x/crypto/bcrypt"
	"log"

	"github.com/sapk/GoWatch/routers"
	"github.com/sapk/GoWatch/routers/admin"
	"github.com/sapk/GoWatch/routers/user"
)

func Start(db *db.Db) {

	m := macaron.New()
	m.Map(db)
	m.Use(macaron.Logger())
	m.Use(macaron.Gziper())
	m.Use(macaron.Static("public"))
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	log.Println("Macaron initialised !")

	m.Use(auth.Authentificator(auth.Options{
		Provider: db,
	}))
	log.Println("Auth initialised !")

	m.Use(macaron.Recovery())

	m.Get("/", routers.Index)
	m.Get("/install", routers.Install)
	m.Post("/install", routers.InstallPost)
	//TODO  determine if we protect the landing page
	//	m.Get("/", auth.IsLogged, routers.Index)

	m.Group("/user", func() {
		m.Get("/", auth.IsLogged, user.CurrentPage)
		m.Get("/login", user.Login)
		m.Post("/login", auth.SignIn)
		m.Get("/logout", auth.LogOut)
	})

	m.Group("/admin", func() {
		m.Get("/", auth.IsLogged, admin.Dashboard)
		m.Get("/users", auth.IsLogged, admin.Users)
		m.Get("/user/add", auth.IsLogged, admin.UserAdd)
		m.Post("/user/add", auth.IsLogged, admin.UserAddPost)
		m.Get("/equipements", auth.IsLogged, admin.Equipements)
	})

	m.Run()
}
