package watcher

import (
	"log"
        "sync"
	"golang.org/x/net/icmp"
        "github.com/sapk/GoWatch/modules/db"
        "github.com/sapk/GoWatch/modules/rrd"
        "strconv"
//        "fmt"
//        "time"
)

type CPingMap struct{
    sync.RWMutex
    m map[string]Ping
}
//Db represent the database
type Watcher struct {
	PingListener *icmp.PacketConn
	PingToListen CPingMap
        PingSeq uint
}

var w Watcher

//TODO log in a RRD database
//TODO better handler concurrecny on map

//Init init the Watcher
func Init(d *db.Db) *Watcher {
	//TODO get ip to mintor form db at start up
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
        //c, err := icmp.ListenPacket("udp4", "0.0.0.0")

	if err != nil {
		log.Fatalf("listen err, %s", err)
	}

	w = Watcher{PingListener: c, PingToListen: CPingMap{m: make(map[string]Ping)}}

        StartPingWatcher()
        
        count ,equipements := d.GetEquipements()
        //var channels []chan PingResponse
        channels := make([]<-chan PingResponse,count)
        
        for i, equi := range equipements {
                //TODO keep the chan to log all responses
                channels[i] = RegisterPingWatch(equi.IP, 0); //We add all equipement to continuous ping
        }
        go func(){
           for rep := range merge(channels...) {
               // at each response
               //TODO log
               log.Println(rep);
               eq, _ := d.GetEquipementbyIP(db.Equipement{IP:rep.IP})
               //eq.Data=fmt.Sprintf("%v",rep)
               eq.Update()
               rrd.AddPing(strconv.FormatUint(eq.ID,10), rep.Time)
           }
        }()
        
        StartLoopPing()
	return &w
}

//Get get the Watcher
func Get() *Watcher {
	return &w
}

func merge(cs ...<-chan PingResponse) <-chan PingResponse {
    var wg sync.WaitGroup
    out := make(chan PingResponse)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan PingResponse) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}