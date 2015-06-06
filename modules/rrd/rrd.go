package rrd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ziutek/rrd"
)

const (
	dbfolder = "data"
	dbfile   = "icmp.rrd"
	//Step The step of logging (60 secs)
	Step = 30
	//Heartbeat minimal time to consider alive
	Heartbeat = 4 * Step
)

//RRD represent the database
type RRD struct {
	//db *rrd.Creator
	u map[string]*rrd.Updater
}

var r RRD

//Create create the rrd file of the object
func Create(ID string) {
	folder := strings.Join([]string{dbfolder, ID}, string(filepath.Separator))
	if src, err := os.Stat(folder); err != nil || !src.IsDir() {
		log.Println("Creating folder : ", folder)
		err := os.Mkdir(folder, 0755)
		if err != nil {
			log.Fatal("Creating element folder failed : ", err)
		}
	}
	file := strings.Join([]string{folder, dbfile}, string(filepath.Separator))
	if _, err := os.Stat(file); os.IsNotExist(err) {
		c := rrd.NewCreator(file, time.Now(), Step)
		c.DS("ping", "GAUGE", Heartbeat, 0, Heartbeat*1000) //Heartbeat in Millisecond
		c.RRA("AVERAGE", 0.5, 2, 60)                        //Each minute * 60 -> 1h
		c.RRA("AVERAGE", 0.5, 2*60/10, 24*10)               //Each hour/10 * 24 -> 1days
		c.RRA("AVERAGE", 0.5, 2*60*24/100, 7*100)           //Each days/100 * 7 -> 1week
		err := c.Create(false)
		if err != nil {
			log.Fatal("Creating element rrd db failed : ", err)
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

//Init initiate all
func Init() *RRD {
	//Create()
	// check if the source dir exist
	if src, err := os.Stat(dbfolder); err != nil || !src.IsDir() {
		err := os.Mkdir(dbfolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	r.u = make(map[string]*rrd.Updater)
	return &r
}

//AddPing add a ping to rrd
func AddPing(ID string, t time.Duration) {
	if _, ok := r.u[ID]; !ok {
		Create(ID)
	}
	log.Println("Adding ping to database ", ID, " / ", t, " -> ", int64(t)/int64(time.Millisecond))
	r.u[ID].Update(time.Now(), int64(t)/int64(time.Millisecond))
}

//GraphPing generate bytes of the img
func GraphPing(ID, title string, start, end time.Time) []byte {
	folder := strings.Join([]string{dbfolder, ID}, string(filepath.Separator))
	file := strings.Join([]string{folder, dbfile}, string(filepath.Separator))
	log.Println(file)
	g := rrd.NewGrapher()
	g.SetTitle(title)
	g.SetVLabel("time in ms")
	g.SetSize(800, 300)
	g.Def("p", file, "ping", "AVERAGE")
	g.VDef("avg", "p,AVERAGE")
	g.VDef("max", "p,MAXIMUM")

	g.Line(1, "p", "ff0000", "ping")
	g.GPrint("avg", "avg=%lf")
	g.GPrint("max", "max=%lf")
	infos, bs, err := g.Graph(start, end)
	if err != nil {
		log.Println("Error in generating rrd ping graph : ", err)
	}
	log.Println(infos)
	return bs
}
