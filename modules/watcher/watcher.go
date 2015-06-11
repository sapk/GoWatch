package watcher

//        "fmt"
//        "time"

//Watcher contain watcher informations
type Watcher struct {
	Ping *PingWatcher
}

var w Watcher

//TODO better handler concurrecny on map

//Init init the Watcher
func Init() *Watcher {
	w := Watcher{Ping: initPingWatcher()}
	return &w
}
