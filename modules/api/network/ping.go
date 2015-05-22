package network

import (
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"
)

// PingResponse represent ip response and stats
type PingResponse struct {
	IP     string
	Result bool
	Time   time.Duration
}

//Ping execute a ping and return informations
func Ping(ip string) PingResponse {
	//TODO
	//TODO implement v6
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("listen err, %s", err)
	}

	defer c.Close()
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("COUCOU"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := c.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(ip)}); err != nil {
		log.Fatalf("WriteTo err, %s", err)
	}

	start := time.Now()
	//TODO break on timeout
PingRespLoop:
	for {
		rb := make([]byte, 1500)

		n, peer, err := c.ReadFrom(rb)
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
			if peer.String() == ip {
				break PingRespLoop
			}
		default:
			log.Printf("got %+v; want echo reply", rm)
		}
	}

	return PingResponse{IP: ip, Result: true, Time: time.Since(start)}
}
