package collector

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	//	"github.com/ziutek/rrd"
	"github.com/boltdb/bolt"
)

const (
	dbfolder = "data"
	dbfile   = "collect.db"
	Step     = 15 * time.Second //interval of time between ping //TODO use config file for these time
)

//Collector represent the database
type Collector struct {
	//db *rrd.Creator
	DBs map[string]*bolt.DB
}

var c Collector

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
	//if _, err := os.Stat(file); os.IsNotExist(err) { //We do not need to check file existence because bolt do it
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		log.Fatal("Creating element collector db failed : ", err)
	}
	c.DBs[ID] = db
	//}
}

//Init initiate all
func Init() *Collector {
	// check if the source dir exist
	if src, err := os.Stat(dbfolder); err != nil || !src.IsDir() {
		err := os.Mkdir(dbfolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	c.DBs = make(map[string]*bolt.DB)
	return &c
}

//AddPing add a ping to database
func AddPing(ID string, t time.Duration) {
	if _, ok := c.DBs[ID]; !ok {
		Create(ID)
	}
	log.Println("Adding ping to database ", ID, " / ", t, " -> ", int64(t)/int64(time.Millisecond))

	c.DBs[ID].Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("PingTime"))
		if err != nil {
			return fmt.Errorf("Error during creating bucket PingTime: %s", err)
		} //TODO create bucket at creation of file ?
		//b := tx.Bucket([]byte("MyBucket"))
		timestamp, err := time.Now().MarshalBinary()
		if err != nil {
			return fmt.Errorf("Error during timestamp generation: %s", err)
		}
		return b.Put(timestamp, []byte(t.String()))
	})
	//c.DBs[ID].Update(time.Now(), int64(t)/int64(time.Millisecond))
}

//GraphPing generate bytes of the img
func GraphPing(ID, title string, start, end time.Time) []byte {
	//TODO
	/*
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
	*/
	return []byte{}
}
