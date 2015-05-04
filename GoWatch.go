package main

import (
	"./modules/auth"
	"./modules/db"
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	_ "github.com/mattn/go-sqlite3"
	//"golang.org/x/crypto/bcrypt"
	"./routers"
	"./routers/admin"
	"./routers/user"
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
	authentificator := auth.New(d)
	m.Map(authentificator)
	//m.Use(authentificator.PermissionHandler)
	log.Println("Auth initialised !")

	m.Use(macaron.Recovery())

	m.Get("/", authentificator.IsLogged, routers.Index)

	m.Group("/user", func() {
		m.Get("/", authentificator.IsLogged, user.CurrentPage)
		m.Get("/login", user.Login)
		m.Post("/login", authentificator.SignIn)
	})

	m.Group("/admin", func() {
		m.Get("/", authentificator.IsLogged, admin.Dashboard)
		m.Get("/users", authentificator.IsLogged, admin.Users)
		m.Get("/observeds", authentificator.IsLogged, admin.Observeds)
	})
	m.Run()
}
