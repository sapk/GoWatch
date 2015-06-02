package watcher

import (
	"log"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"
)

//Db represent the database
type Watcher struct {
	PingListener *icmp.PacketConn
	PingToListen map[string]Ping
}

var w Watcher

//TODO log in a RRD database
//TODO better handler concurrecny on map

//Init init the Watcher
func Init() *Watcher {
	//TODO get ip to mintor form db at start up
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("listen err, %s", err)
	}

	w = Watcher{PingListener: c, PingToListen: make(map[string]Ping)}
	//Clearer check if needed
	go func() {
		//TODO support continuous ping
		for {
			//every minutes we check for timeout and clean the map
			time.Sleep(1 * time.Minute)
			log.Println("Scanning PingToListen map:", w.PingToListen)
			for ip, ping := range w.PingToListen {
				log.Println("IP:", ip, "Ping:", ping)
				if ping.Nb == 0 || time.Since(ping.Start) > ping.Timeout {
					log.Println("Clearing IP:", ip, "Ping:", ping)
					close(ping.Ch)
					delete(w.PingToListen, ip)
				}
			}
		}
	}()
	//Ping watcher
	go func() {
		//TODO clear PingTowatchList for timeout
		for {
			rb := make([]byte, 1500)

			n, peer, err := w.PingListener.ReadFrom(rb)
			if err != nil {
				log.Fatal(err)
			}

			rm, err := icmp.ParseMessage(iana.ProtocolICMP, rb[:n])
			if err != nil {
				log.Fatal(err)
			}

			switch rm.Type {
			case ipv4.ICMPTypeEchoReply:
				log.Printf("got reflection from %v", peer)
				if ping, ok := w.PingToListen[peer.String()]; ok {
					if time.Since(ping.Start) < ping.Timeout {
						//if timeout isn't pass
						log.Printf("Sending to chan for %v ...", peer)
						w.PingToListen[peer.String()].Ch <- PingResponse{IP: peer.String(), Result: true, Time: time.Since(ping.Start)}
					}
					ping.Nb--
					w.PingToListen[peer.String()] = ping
					if w.PingToListen[peer.String()].Nb == 0 {
						log.Println("Clearing IP:", peer, "Ping:", ping)
						close(w.PingToListen[peer.String()].Ch)
						delete(w.PingToListen, peer.String())
					}
				}
			default:
				log.Printf("got %+v; want echo reply", rm)
			}
		}
	}()

	return &w
}

//Get get the Watcher
func Get() *Watcher {
	return &w
}
