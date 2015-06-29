package collector

import (
	"fmt"
	"github.com/boltdb/bolt"
	//"github.com/sapk/GoWatch/models/equipement"
	"log"
	"os"
	"path/filepath"
	//"strconv"
	"strings"
	"time"
)

const (
	dbfolder  = "data"
	dbfile    = "collect.db"
	Step      = 5 * time.Second //interval of time between ping //TODO use config file for these time
	MaxMissed = 3               //Use to determine if no ping response
)

// Ping represent ip response and stats
type Ping struct {
	/*IP     string*/
	Result bool
	Time   time.Duration
	At     time.Time
}

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
	//
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
	log.Println("Adding ping to database ", ID, " : ", t)

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
		return b.Put(timestamp, []byte(t.String())) //TODO determine if better in nanosecond ?
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

//DataPing send array of ping response
func DataPing(ID string, start, end time.Time) []Ping {
	//TODO
	ret := make([]Ping, 0, 1000) //TODO better
	//id, _ := strconv.ParseUint(ID, 10, 64)
	//eq, _ := equipement.GetByID(id)
	if _, ok := c.DBs[ID]; !ok {
		Create(ID)
	}
	log.Println("Getting ping of database ", ID)
	//finish := make(chan error)
	c.DBs[ID].View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PingTime"))
		c := b.Cursor()
		///TODO fileter base on end and start
		for k, v := c.First(); k != nil; k, v = c.Next() {
			//log.Printf("key=%s, value=%s\n", k, v)
			dur, _ := time.ParseDuration(string(v))
			at := time.Now() //TODO better creation
			at.UnmarshalBinary(k)

			if at.After(start) && at.Before(end) {
				//ret = append(ret, Ping{IP: eq.IP(), Result: true, Time: dur, At: at}) //TODO add false base on step and MaxMissed
				ret = append(ret, Ping{Result: true, Time: dur, At: at}) //TODO add false base on step and MaxMissed
			}
		}
		log.Println("Db parsing finish sending to chan", ret)
		//finish <- nil
		return nil
	})
	//err := <-finish
	/*
		if err != nil {
			log.Println("error during DB parsing", err)
		}
	*/
	log.Println("Db parsing finish returning array", ret)

	return ret
}
