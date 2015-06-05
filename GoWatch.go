package main

import (
	"log"

	"github.com/sapk/GoWatch/modules/db"
	"github.com/sapk/GoWatch/modules/watcher"
	"github.com/sapk/GoWatch/modules/web"
        "github.com/sapk/GoWatch/modules/rrd"
)

func main() {

	log.Println("Start !")
	
	d := db.InitDb()
	log.Println("Db initialised !")
	
	w := watcher.Init(d)
	log.Println("Watcher initialised !")
	
        r := rrd.Init()
        log.Println("RRD initialised !")

	web.Start(d,w,r)
}
