package watcher

import (
	"errors"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/sapk/GoWatch/modules/rrd"
	"github.com/sapk/GoWatch/modules/tools"

	"bytes"
	"encoding/binary"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"
)

const maxUniqePingTimeout = 15 * time.Second

//TODO mutex on listenpinglist

// PingResponse represent ip response and stats
type PingResponse struct {
	//TODO check PID for long running or multiple ping
	IP     string
	Result bool
	Time   time.Duration
	Error  string
}

//PingSend format of a Ping send
type PingSend struct {
	at time.Time
}

//Ping contain information about a runnnig ping
type Ping struct {
	Ch      tools.Broadcaster
	Start   time.Time //TODO remove this and base calcul on Send[] average
	Send    map[int]PingSend
	Timeout time.Duration // a infinite send will have 0 here
}

//PingTest execute a ping and return un litenerfor response
func PingTest(ip string, timeout time.Duration) PingResponse {
	ping, err := RegisterPingWatch(ip, timeout)
	SendPing(ip)

	if err != nil {
		return PingResponse{IP: ip, Result: false, Time: 0, Error: err.Error()}
	}

	return ping.Read().(PingResponse)
}

//RegisterPingWatch add to the watch list
func RegisterPingWatch(ip string, timeout time.Duration) (*tools.BroadcastReceiver, error) {
	//TODO use a global event chan
	out := tools.NewBroadcaster()
	log.Println("Adding ", ip, "to watch list")
	//If ip is invali we do nothing except send a massaeg contian error
	if ok, _ := regexp.MatchString(tools.ValidIPAddressRegex, ip); !ok {
		//out <- PingResponse{IP: ip, Result: false, Time: 0, Error: "Invalid IP"}
		return nil, errors.New("Invalid IP")
	}

	//If we don't have it we make it
	w.PingToListen.RLock()
	if ping, ok := w.PingToListen.m[ip]; !ok {
		log.Println("Creating ", ip, " element to watch list for ", timeout)
		w.PingToListen.m[ip] = Ping{Start: time.Now(), Ch: out, Timeout: timeout, Send: make(map[int]PingSend)}
	} else if timeout == 0 {
		//If we register for unlimited listen and the ip is already in listen but not neccesrry in unltimate
		log.Println("There is ", ip, " element setting him for unlimited ", timeout)
		ping.Timeout = timeout
		w.PingToListen.m[ip] = ping
	} else {
		//We reset start counter for timeout
		log.Println("There is ", ip, " element resseting him for ping at Start ", timeout)
		ping.Start = time.Now()
		w.PingToListen.m[ip] = ping

	}
	w.PingToListen.RUnlock()
	//Si le timeout est supérieur à 0 le minimum on active le timeout
	if timeout > 0 {
		go func() {
			time.Sleep(timeout)
			ClearPingIfNeeded(ip)
		}()
	}
	listener := w.PingToListen.m[ip].Ch.Listen()
	return &(listener), nil
}

//ClearPingIfNeeded analyze if we need to do clean of the ip related objects
func ClearPingIfNeeded(ip string) {
	w.PingToListen.RLock()
	//Verify that it still exist
	if ping, ok := w.PingToListen.m[ip]; ok {
		//We verify that it doesn't become unlimited during the timeout
		if ping.Timeout > 0 && time.Since(ping.Start) > ping.Timeout {
			log.Println("Clearing IP:", ip, "Ping:", ping)
			ping.Ch.Write(PingResponse{IP: ip, Result: false, Time: time.Since(ping.Start)})
			delete(w.PingToListen.m, ip)
		} else {
			//If it 's unlimited or not timeouted we clear all Send maxuniqtimeouted
			for seq, send := range ping.Send {
				if time.Since(send.at) > maxUniqePingTimeout {
					delete(ping.Send, seq)
				}
			}

			w.PingToListen.m[ip] = ping
		}
	}
	w.PingToListen.RUnlock()
}

//SendPing send a ping packet
func SendPing(ip string) int {
	//TODO implement v6
	//If ip is invali we do nothing
	if ok, _ := regexp.MatchString(tools.ValidIPAddressRegex, ip); !ok {
		log.Println("Invalid IP")
		return -1
	}

	w.PingToListen.RLock()
	ping, ok := w.PingToListen.m[ip]
	//If we don't wait for a response we don't send anything
	if !ok {
		log.Println("Don't send a Ping we don't listen to is response")
		return -1
	}
	w.PingSeq++
	seq := int(w.PingSeq) & 0xffff
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: seq,
			Data: []byte("COUCOU"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Sending ping to ", ip)
	if _, err := w.PingListener.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(ip)}); err != nil {
		log.Printf("WriteTo err, %s", err)
	} else {
		ping.Send[seq] = PingSend{at: time.Now()}
		w.PingToListen.m[ip] = ping
	}

	w.PingToListen.RUnlock()
	return seq
}

//StartPingWatcher start the goroutine for watch receiving of ping response
func StartPingWatcher() {

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
				//http://www.hsc.fr/ressources/articles/protocoles/icmp/index.html
				log.Printf("got reflection from %v", peer)
				ip := peer.String()
				ClearPingIfNeeded(ip) //do all the cleaning stuff for timeouted

				w.PingToListen.RLock()
				if ping, ok := w.PingToListen.m[ip]; ok {
					b, _ := rm.Body.Marshal(4)
					buf := bytes.NewReader(b) // b is []byte
					var uid, useq uint16
					binary.Read(buf, binary.BigEndian, &uid)
					binary.Read(buf, binary.BigEndian, &useq)

					seq := int(useq)
					if send, ok := ping.Send[seq]; ok {
						log.Printf("Sending to chan for %v ...", ip)
						ping.Ch.Write(PingResponse{IP: ip, Result: true, Time: time.Since(send.at)})
						delete(ping.Send, seq)
						w.PingToListen.m[ip] = ping
					}
				}
				w.PingToListen.RUnlock()
			default:
				//log.Printf("got %+v; want echo reply", rm)
			}
		}
	}()

}

//StartLoopPing start the go routing sending recursively ping for element in database
func StartLoopPing() {
	//Loop ping
	go func() {
		//TODO support continuous ping
		for {
			//every maxUniqePingTimeout we check for timeout and clean the map
			//time.Sleep(maxUniqePingTimeout)
			//log.Println("Scanning PingToListen map:", w.PingToListen)
			for ip := range w.PingToListen.m {
				ClearPingIfNeeded(ip) //We do cleanup before every thing

				w.PingToListen.RLock()
				if ping, ok := w.PingToListen.m[ip]; ok && ping.Timeout == 0 && len(ping.Send) == 0 {
					//the pin has not been cleared and it's a contnious and we are not expecting a ping
					//So we could send another
					//TODO make sure that we pass each step for each el so here we take amargin of /2  but could better if coudl ping at excatly Step bettween each
					SendPing(ip)
					timetowait := (int64(rrd.Step*time.Second) / int64(len(w.PingToListen.m)))
					//log.Println("Wainting :",time.Duration(timetowait))
					time.Sleep(time.Duration(timetowait)) //scale for the number waiting
				} else {
					log.Println("Skipping because it's not a continous ping or a ping is already pending", ip)
				}
				w.PingToListen.RUnlock()
			}
		}
	}()
}
