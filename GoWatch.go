package main

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
	//"golang.org/x/crypto/bcrypt"
	"github.com/sapk/GoWatch/routers"
	"github.com/sapk/GoWatch/routers/admin"
	"github.com/sapk/GoWatch/routers/user"
	"log"
)

func main() {

	log.Println("Start !")
	d := db.InitDb()
	log.Println("Db initialised !")

	(*d).CreateAdminUser()

	m := macaron.New()
	m.Map(d)
	m.Use(macaron.Logger())
	m.Use(macaron.Gziper())
	m.Use(macaron.Static("public"))
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	log.Println("Macaron initialised !")

	//m.Map(auth.New(d))
	//TODO define auth as a middleware
	//authentificator := auth.New(d)
	//m.Map(authentificator)
	m.Use(auth.Authentificator(auth.Options{
		Provider: d,
	}))
	log.Println("Auth initialised !")

	m.Use(macaron.Recovery())

	m.Get("/", auth.IsLogged, routers.Index)

	m.Group("/user", func() {
		m.Get("/", auth.IsLogged, user.CurrentPage)
		m.Get("/login", user.Login)
		m.Post("/login", auth.SignIn)
	})

	m.Group("/admin", func() {
		m.Get("/", auth.IsLogged, admin.Dashboard)
		m.Get("/users", auth.IsLogged, admin.Users)
		m.Get("/observeds", auth.IsLogged, admin.Observeds)
	})
	m.Run()
}
