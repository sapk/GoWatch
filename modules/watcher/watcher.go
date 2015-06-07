package watcher

import (
	"log"
	"strconv"
	"sync"

	"github.com/sapk/GoWatch/modules/db"
	"github.com/sapk/GoWatch/modules/rrd"
	"github.com/sapk/GoWatch/modules/tools"
	"golang.org/x/net/icmp"
	//        "fmt"
	//        "time"
)

//CPingMap with mutex for concurrency
type CPingMap struct {
	sync.RWMutex
	m map[string]Ping
}

//Watcher contain watcher informations
type Watcher struct {
	DB           *db.Db
	PingListener *icmp.PacketConn
	PingToListen CPingMap
	PingSeq      uint
	PingChannels chan PingResponse
}

var w Watcher

//TODO better handler concurrecny on map

//Init init the Watcher
func Init(d *db.Db) *Watcher {
	//TODO get ip to mintor form db at start up
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	//c, err := icmp.ListenPacket("udp4", "0.0.0.0")

	if err != nil {
		log.Fatalf("listen err, %s", err)
	}

	w = Watcher{PingListener: c, DB: d, PingToListen: CPingMap{m: make(map[string]Ping)}}

	StartPingWatcher()

	UpdatePingChannels()
	//TODO stop and restart

	StartLoopPing()
	return &w
}

//UpdatePingChannels clean a set the watching of channels
func UpdatePingChannels() {

	count, equipements := w.DB.GetEquipements()
	log.Println("There is ", count, " elements in db")

	//var channels []chan PingResponse
	channels := make(map[string]*tools.BroadcastReceiver) //We reset completely the map
	for _, eq := range equipements {
		channels[eq.IP], _ = RegisterPingWatch(eq.IP, 0) //We add all equipement to continuous ping
	}

	//Clearing PingToListen map from removed elements
	w.PingToListen.RLock()
	for ip, listen := range w.PingToListen.m {
		//channels[strconv.FormatUint(eq.ID,10)] = RegisterPingWatch(eq.IP, 0); //We add all equipement to continuous ping
		if listen.Timeout == 0 {
			//We only clear long running ping
			if _, ok := channels[ip]; !ok {
				//We  clear if it's not in long running ping
				delete(w.PingToListen.m, ip)
			}
		}
	}
	w.PingToListen.RUnlock()

	if w.PingChannels != nil {
		close(w.PingChannels) //TODO use WaitGroup to close the go routine parsing the PingChannels
	}
	var outIsClosed *bool
	w.PingChannels, outIsClosed = merge(channels)
	go func() {
		for {
			rep, ok := <-w.PingChannels
			if !ok {
				log.Println("Done the chan must has been reset")
				*outIsClosed = true
				return
			}
			log.Println(rep)
			eq, err := w.DB.GetEquipementbyIP(db.Equipement{IP: rep.IP}) //TODO check if it exist before logging
			//eq.Data=fmt.Sprintf("%v",rep)
			if err != nil {
				log.Println("Not found in database : ", err)
			} else {
				eq.Update()
				rrd.AddPing(strconv.FormatUint(eq.ID, 10), rep.Time)
			}
		}
	}()
}

//Get get the Watcher
func Get() *Watcher {
	return &w
}

func merge(cs map[string]*tools.BroadcastReceiver) (chan PingResponse, *bool) {
	var wg sync.WaitGroup
	outIsClosed := false
	out := make(chan PingResponse)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c *tools.BroadcastReceiver) {
		for {
			if c == nil {
				log.Println("This chan must has been close")
				wg.Done()
				return
			}
			n := c.Read().(PingResponse)
			if outIsClosed {
				log.Println("The output chan as been closed")
				wg.Done()      // clear the wait group for closing all still open chan
				SendPing(n.IP) //We resend for any other chan taht will listen after
				return
			}
			/*
			   if  _, ok := <- out; !ok {
			       log.Println("The output chan as been closed")
			       for _, c := range cs {
			               if _, ok := <- c; ok {
			                       wg.Done()// clear the wait group for closing all still open chan
			               }
			       }
			       return;
			   }
			*/
			out <- n

		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		if !outIsClosed {
			close(out)
		}
	}()
	return out, &outIsClosed
}
