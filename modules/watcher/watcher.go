package watcher

import (
	"log"

	"golang.org/x/net/icmp"
)

//Db represent the database
type Watcher struct {
	PingListener *icmp.PacketConn
	PingToListen map[string]Ping
        PingSeq uint
}

var w Watcher

//TODO log in a RRD database
//TODO better handler concurrecny on map

//Init init the Watcher
func Init() *Watcher {
	//TODO get ip to mintor form db at start up
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
        //c, err := icmp.ListenPacket("udp4", "0.0.0.0")

	if err != nil {
		log.Fatalf("listen err, %s", err)
	}

	w = Watcher{PingListener: c, PingToListen: make(map[string]Ping)}

        StartPingWatcher()
        StartLoopPing()
	return &w
}

//Get get the Watcher
func Get() *Watcher {
	return &w
}
