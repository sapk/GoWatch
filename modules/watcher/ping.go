package watcher

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/sapk/GoWatch/models/equipement"
	"github.com/sapk/GoWatch/modules/rrd"
	"github.com/sapk/GoWatch/modules/tools"

	"golang.org/x/net/icmp"
	"golang.org/x/net/internal/iana"
	"golang.org/x/net/ipv4"
)

// PingWatcher represent the object used for watch ping
type PingWatcher struct {
	PingListener *icmp.PacketConn
	PingToListen PingMap
	PingSeq      uint
	PingChannels PingRequest
}

//PingMap with mutex for concurrency
type PingMap struct {
	sync.RWMutex
	m map[string]*Ping
}

//Ping contain information about a runnnig ping
type Ping struct {
	Ch      chanListPingRequest
	Send    map[int]PingSend
	Timeout time.Duration // a infinite send will have 0 here
}

//PingSend format of a Ping send
type PingSend struct {
	at time.Time
}

// PingResponse represent ip response and stats
type PingResponse struct {
	IP     string
	Result bool
	Time   time.Duration
	Error  string
}

// PingRequest represent ip request
type PingRequest struct {
	isClose bool
	ch      chan PingResponse
}

var pw PingWatcher

const maxUniqePingTimeout = 15 * time.Second

type chanListPingRequest map[int]*PingRequest

func (cs *chanListPingRequest) add(c *PingRequest) {
	if !cs.has(c) {
		(*cs)[len(*cs)] = c
		log.Println("Adding one chan to the chan list of listener")
	}
}
func (cs *chanListPingRequest) has(c *PingRequest) bool {
	for _, ch := range *cs {
		if c == ch {
			log.Println("Found the chan in the chan list")
			return true
		}
	}
	return false
}
func (cs *chanListPingRequest) send(rep PingResponse) {
	log.Println(len(*cs))
	for id, req := range *cs {
		log.Println("Trying to send to one chan", rep)
		if req != nil && !req.isClose {
			req.ch <- rep
		} else {
			log.Println("One chan seem to be dead removing it !")
			delete(*cs, id) //TODO make sure it take in accoutn in pw
		}
	}
}

//PingWatcher init the a PingWatcher
func initPingWatcher() *PingWatcher {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	//c, err := icmp.ListenPacket("udp4", "0.0.0.0")

	if err != nil {
		log.Fatalf("listen err, %s", err)
	}

	pw = PingWatcher{PingListener: c, PingToListen: PingMap{m: make(map[string]*Ping)}, PingSeq: 0, PingChannels: PingRequest{false, make(chan PingResponse)}}

	startPingWatcher()
	startWatchLongRunningPing()
	startLoopPing()
	return &pw
}

//startWatchLongRunningPing start the goroutine for parse pignResponse from long running
func startWatchLongRunningPing() {
	//we take at startthe allready in db equipement
	count, equipements := equipement.GetAll()
	log.Println("There is ", count, " elements in db")

	for _, eq := range *equipements {
		AddToPingLongRunningList(eq.IP())
	}

	//TODO goroutine parsing pw.PingChannels
	go func() {
		for {
			rep, ok := <-pw.PingChannels.ch
			if !ok {
				log.Fatalln("The chan must has been reset") //Should not happen with the new implementation
			}
			log.Println(rep)
			eq, ok := equipement.GetByIP(rep.IP) //TODO check if it exist before logging
			//eq.Data=fmt.Sprintf("%v",rep)
			if !ok {
				log.Println("Not found in database : ")
				//We should remove it from the longrunning ping list
				removeFromPingList(rep.IP)
			} else {
				if rep.Result == true {
					eq.UpdateActivity() //TODO cache in order to not lock the db
					rrd.AddPing(strconv.FormatUint(eq.ID(), 10), rep.Time)
				} else {
					//Timeout
				}
			}
		}
	}()
}

//startPingWatcher start the goroutine for watch receiving of ping response
func startPingWatcher() {

	//Ping watcher
	go func() {
		for {
			rb := make([]byte, 1500)

			n, peer, err := pw.PingListener.ReadFrom(rb)
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
				clearPingIfNeeded(ip) //do all the cleaning stuff for timeouted

				pw.PingToListen.RLock()
				if ping, ok := pw.PingToListen.m[ip]; ok {
					b, _ := rm.Body.Marshal(4)
					buf := bytes.NewReader(b) // b is []byte
					var uid, useq uint16
					binary.Read(buf, binary.BigEndian, &uid)
					binary.Read(buf, binary.BigEndian, &useq)

					seq := int(useq)
					if send, ok := ping.Send[seq]; ok {
						log.Printf("Sending to chan for %v ...", ip)
						ping.Ch.send(PingResponse{IP: ip, Result: true, Time: time.Since(send.at)})
						log.Println("Clearing Seq for response receive :", ip, "Ping:", ping, "Seq", seq)
						delete(ping.Send, seq)
					}
				}
				pw.PingToListen.RUnlock()
			default:
				//log.Printf("got %+v; want echo reply", rm)
			}
		}
	}()

}

//AddToPingLongRunningList AddToPingLongRunningList
func AddToPingLongRunningList(ip string) {
	err := registerPingWatch(ip, 0, &pw.PingChannels) //We add all equipement to continuous ping
	if err != nil {
		log.Println("error during registering long running ping", err)
	}
}

//removeFromPingList RemoveFromPingList
func removeFromPingList(ip string) {
	pw.PingToListen.RLock()
	if ping, ok := pw.PingToListen.m[ip]; ok {
		for id := range ping.Ch {
			//We don't close all chan because there are managed by upper call
			delete(ping.Ch, id)
		}
		delete(pw.PingToListen.m, ip)
	}
	pw.PingToListen.RUnlock()
}

//clearPingIfNeeded analyze if we need to do clean of the ip related objects
func clearPingIfNeeded(ip string) {
	pw.PingToListen.RLock()
	//Verify that it still exist
	if ping, ok := pw.PingToListen.m[ip]; ok {
		//We remove all timeouted packet send or maxuniqtimeouted for unlimitedping
		for id, send := range ping.Send {
			if (ping.Timeout > 0 && time.Since(send.at) > ping.Timeout) || ping.Timeout == 0 && time.Since(send.at) > maxUniqePingTimeout {
				ping.Ch.send(PingResponse{IP: ip, Result: false, Time: time.Since(send.at)})
				log.Println("Clearing Seq for timeout :", ip, "Ping:", ping, "Seq", id)
				delete(ping.Send, id)
			}
		}
		//If the ping isn't unlimited and we don't wait for ping anymore
		if ping.Timeout > 0 && len(ping.Send) == 0 {
			log.Println("Clearing IP:", ip, "Ping:", ping)
			for id := range ping.Ch {
				//We don't close all chan because there are managed by upper call
				delete(ping.Ch, id)
			}
			delete(pw.PingToListen.m, ip)
		} else {
			//If it 's unlimited or not timeouted we save all changes (doesn't needed anymore)
			//pw.PingToListen.m[ip] = ping
		}
	}
	pw.PingToListen.RUnlock()
}

//startLoopPing start the go routing sending recursively ping for element in database
func startLoopPing() {
	//Loop ping
	go func() {
		for {
			for ip, ping := range pw.PingToListen.m {
				if ping.Timeout == 0 {
					//it's a contnious and we clean only continus since other will be clean by timeout
					clearPingIfNeeded(ip) //We do cleanup before every thing
					//if ping, ok := pw.PingToListen.m[ip]; ok && len(ping.Send) == 0 {
					if len(ping.Send) == 0 {
						//the ping has not been cleared and  and we are not expecting a ping
						//So we could send another
						sendPing(ip)
					} else {
						log.Println("Skipping because a ping is already pending", ip)
					}
				} else {
					//No cleaning in order to not block the map
					log.Println("Skipping because it's not a continous ping or a ping is already pending or deleted", ip)
					//					time.Sleep(500 * time.Millisecond) //TODO scale for the number waiting timelapse to not block the entire process
				}
				timetowait := (int64(rrd.Step*time.Second) / int64(len(pw.PingToListen.m)+1))
				//log.Println("Wainting :",time.Duration(timetowait))
				time.Sleep(time.Duration(timetowait) + 5*time.Millisecond) //scale for the number waiting
			}
		}
	}()
}

//sendPing send a ping packet
func sendPing(ip string) int {
	//TODO implement v6
	//If ip is invali we do nothing
	if ok, _ := regexp.MatchString(tools.ValidIPAddressRegex, ip); !ok {
		log.Println("Invalid IP")
		return -1
	}

	pw.PingToListen.RLock()
	ping, ok := pw.PingToListen.m[ip]
	//If we don't wait for a response we don't send anything
	if !ok {
		log.Println("Don't send a Ping we don't listen to is response")
		return -1
	}
	pw.PingSeq++
	seq := int(pw.PingSeq) & 0xffff
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
	if _, err := pw.PingListener.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(ip)}); err != nil {
		log.Printf("WriteTo err, %s", err)
	} else {
		ping.Send[seq] = PingSend{at: time.Now()}
		//pw.PingToListen.m[ip] = ping
	}

	pw.PingToListen.RUnlock()
	return seq
}

//registerPingWatch add to the watch list
func registerPingWatch(ip string, timeout time.Duration, ch *PingRequest) error {
	log.Println("Adding ", ip, "to watch list")
	//If ip is invali we do nothing except send a massaeg contian error
	if ok, _ := regexp.MatchString(tools.ValidIPAddressRegex, ip); !ok {
		//out <- PingResponse{IP: ip, Result: false, Time: 0, Error: "Invalid IP"}
		return errors.New("Invalid IP")
	}

	//If we don't have it we make it
	pw.PingToListen.RLock()
	if ping, ok := pw.PingToListen.m[ip]; !ok {
		log.Println("Creating ", ip, " element to watch list for ", timeout)
		ping := Ping{Ch: make(chanListPingRequest), Timeout: timeout, Send: make(map[int]PingSend)}
		ping.Ch.add(ch)
		pw.PingToListen.m[ip] = &ping
	} else {
		//If we have the element and it's finish
		log.Println("There is ", ip, " element in the array")
		ping.Ch.add(ch) //This will add if not already in map
		if timeout == 0 {
			ping.Timeout = timeout
			log.Println(" ... setting him for unlimited ", timeout)
		}
		//ping.Timeout = timeout
	}

	pw.PingToListen.RUnlock()
	//Si le timeout est supérieur à 0 le minimum on active le timeout
	if timeout > 0 {
		go func() {
			time.Sleep(timeout)
			clearPingIfNeeded(ip)
		}()
	}
	return nil
}

//PingTest execute a ping and return un litenerfor response
func PingTest(ip string, timeout time.Duration) PingResponse {
	ping := PingRequest{false, make(chan PingResponse)}
	//defer close(ping) //TODO think about it
	err := registerPingWatch(ip, timeout, &ping)
	sendPing(ip)

	if err != nil {
		return PingResponse{IP: ip, Result: false, Time: 0, Error: err.Error()}
	}

	ret := <-ping.ch
	ping.isClose = true
	return ret
}
