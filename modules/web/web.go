package web

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	_ "github.com/mattn/go-sqlite3" //used by orm
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
	//"golang.org/x/crypto/bcrypt"
	"log"

	"github.com/sapk/GoWatch/modules/rrd"
	"github.com/sapk/GoWatch/modules/watcher"

	"github.com/sapk/GoWatch/routers"
	"github.com/sapk/GoWatch/routers/admin"
	"github.com/sapk/GoWatch/routers/api"
	"github.com/sapk/GoWatch/routers/equipement"
	"github.com/sapk/GoWatch/routers/user"

	"github.com/macaron-contrib/csrf"
	"github.com/macaron-contrib/toolbox"
)

//Start init the web interface
func Start(db *db.Db, watcher *watcher.Watcher, rrd *rrd.RRD) {

	m := macaron.New()
	m.Map(db)
	m.Map(watcher)
	m.Use(macaron.Logger())
	m.Use(macaron.Gziper())
	m.Use(macaron.Static("public"))
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())
	log.Println("Macaron initialised !")

	m.Use(auth.Authentificator(auth.Options{
		Provider: db,
	}))
	log.Println("Auth initialised !")

	m.Use(macaron.Recovery())

	//TODO remove after dev  /debug
	m.Use(toolbox.Toolboxer(m))

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
	m.Group("/equipement", func() {
		m.Get("s", auth.IsLogged, equipement.Dashboard)
		m.Get("/:id([0-9]+)", auth.IsLogged, equipement.View)
	})
	m.Group("/admin", func() {
		m.Get("/", auth.IsLogged, admin.Dashboard)
		m.Get("/users", auth.IsLogged, admin.Users)
		m.Get("/user/add", auth.IsLogged, admin.UserAdd)
		m.Post("/user/add", auth.IsLogged, admin.UserAddPost)
		m.Get("/user/:id([0-9]+)/del", auth.IsLogged, admin.UserDel)
		m.Get("/equipements", auth.IsLogged, admin.Equipements)
		m.Get("/equipement/add", auth.IsLogged, admin.EquipementAdd)
		m.Post("/equipement/add", auth.IsLogged, admin.EquipementAddPost)
		m.Get("/equipement/add_multiples", auth.IsLogged, admin.EquipementAddMultiples)
		m.Get("/equipement/:id([0-9]+)/del", auth.IsLogged, admin.EquipementDel)
	})

	m.Group("/api", func() {
		m.Get("/network/ping", auth.IsLogged, api.Ping)
		m.Get("/network/snmptest", auth.IsLogged, api.SNMPTest)
		m.Get("/graph/equipement/:id([0-9]+)/:duration.png", auth.IsLogged, api.GraphPing)
	})
	m.Run()
}
