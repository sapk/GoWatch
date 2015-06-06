package rrd

import (
        "log"
        "time"
        "os"
        "github.com/ziutek/rrd"
        "path/filepath"
        "strings"
)


const (
    dbfolder  = "data"
    dbfile    = "icmp.rrd"
    Step      = 60 //60 secs
    Heartbeat = 5 * Step
)

//RRD represent the database
type RRD struct {
        //db *rrd.Creator
        u  map[string]*rrd.Updater
}

var r RRD

func Create(ID string)  {
        folder := strings.Join([]string{dbfolder,ID}, string(filepath.Separator))
        if src, err := os.Stat(folder); err!=nil || !src.IsDir() {
           log.Println("Creating folder : ", folder)
           err := os.Mkdir(folder, 0755)
           if err != nil {
               log.Fatal("Creating element folder failed : ",err)
           }
        }
        file := strings.Join([]string{folder,dbfile}, string(filepath.Separator))
        if _, err := os.Stat(file); os.IsNotExist(err) {
                c := rrd.NewCreator(file, time.Now(), Step)
                c.DS("ping", "GAUGE", Heartbeat, 0, Heartbeat*1000*1000) //Heartbeat in Microsecond
                c.RRA("AVERAGE", 0.5, 5, 20)
                err := c.Create(false);
                if err != nil {
                    log.Fatal("Creating element rrd db failed : ",err)
                }
        }
        
        r.u[ID] = rrd.NewUpdater(file)
        /*
        c.RRA("AVERAGE", 0.5, 1, 100)
        c.RRA("AVERAGE", 0.5, 5, 100)
        c.DS("cnt", "COUNTER", heartbeat, 0, 100)
        c.DS("g", "GAUGE", heartbeat, 0, 60)
        */
}

func Init() *RRD {
        //Create()
        // check if the source dir exist
        if src, err := os.Stat(dbfolder); err!=nil || !src.IsDir() {
           err := os.Mkdir(dbfolder, 0755)
           if err != nil {
               log.Fatal(err)
           }
        }
     
        r.u =  make(map[string]*rrd.Updater)
        return &r
}

func AddPing(ID string, t time.Duration){
        if _, ok := r.u[ID]; !ok {
            Create(ID)
        }
        log.Println("Adding ping to database ",ID," / ",t);
        r.u[ID].Update(time.Now(), float64(t)/float64(time.Microsecond))
}