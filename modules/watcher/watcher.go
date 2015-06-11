package watcher

import "github.com/sapk/GoWatch/modules/db"

//        "fmt"
//        "time"

//Watcher contain watcher informations
type Watcher struct {
	DB   *db.Db
	Ping *PingWatcher
}

var w Watcher

//TODO better handler concurrecny on map

//Init init the Watcher
func Init(d *db.Db) *Watcher {
	w := Watcher{DB: d, Ping: initPingWatcher(d)}
	return &w
}
