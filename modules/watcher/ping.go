package watcher

import (
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

//TODO check PID for long running or multiple ping
// PingResponse represent ip response and stats
type PingResponse struct {
	IP     string
	Result bool
	Time   time.Duration
	Error  string
}

type Ping struct {
	Nb      int
	Ch      chan PingResponse
	Start   time.Time
	Timeout time.Duration
}

func PingTest(ip string, timeout time.Duration) <-chan PingResponse {
	ping := RegisterPingWatch(ip, timeout)
	SendPing(ip)
	return ping
}
func RegisterPingWatch(ip string, timeout time.Duration) <-chan PingResponse {
	//TODO use a global event chan
	out := make(chan PingResponse, 1)
	//TODO check up is doesn't exist how we handle multiplicity ? a array of ch ?
	//Implement timeout here
	w.PingToListen[ip] = Ping{Nb: 1, Start: time.Now(), Ch: out, Timeout: timeout}

	go func() {
		time.Sleep(timeout)
		if ping, ok := w.PingToListen[ip]; ok {
			log.Println("Clearing IP:", ip, "Ping:", ping)
			ping.Ch <- PingResponse{IP: ip, Result: false, Time: time.Since(ping.Start)}
			close(ping.Ch)
			delete(w.PingToListen, ip)
		}
	}()

	return out
}

//Get get the Watcher
func SendPing(ip string) {
	//TODO implement v6
	log.Println("Sending ping to ", ip)
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("COUCOU"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
                log.Println(err)
	}
	Ip := net.ParseIP(ip)
        if Ip == nil {
                log.Println("IP invalide",err)
        }
	if _, err := w.PingListener.WriteTo(wb, &net.IPAddr{IP: Ip}); err != nil {
		log.Fatalf("WriteTo err, %s", err)
	}
}
