package rrd

import (
        "log"
        "time"
        "os"
        "github.com/ziutek/rrd"
)


const (
    dbfile    = "data/icmp.rrd"
    Step      = 60 //60 secs
    Heartbeat = 3 * Step
)

//RRD represent the database
type RRD struct {
        //db *rrd.Creator
        u  *rrd.Updater
}

var r RRD

func Create()  {
        if _, err := os.Stat(dbfile); os.IsNotExist(err) {
                c := rrd.NewCreator(dbfile, time.Now(), Step)
                c.DS("ping", "GAUGE", Heartbeat, 0, Heartbeat)
                c.RRA("AVERAGE", 0.5, 5, 20)
                err := c.Create(false);
                if err != nil {
                   log.Fatal(err)
                }
        }
        
        r.u = rrd.NewUpdater(dbfile)
        /*
        c.RRA("AVERAGE", 0.5, 1, 100)
        c.RRA("AVERAGE", 0.5, 5, 100)
        c.DS("cnt", "COUNTER", heartbeat, 0, 100)
        c.DS("g", "GAUGE", heartbeat, 0, 60)
        */
}

func Init() *RRD {
        Create()
        return &r
}

func Add(data float64){
        log.Println("Adding to database",data);
        r.u.Update(time.Now(), data)
}