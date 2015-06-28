package main

import (
	"log"

	"github.com/sapk/GoWatch/modules/collector"
	"github.com/sapk/GoWatch/modules/watcher"
	"github.com/sapk/GoWatch/modules/web"
)

func main() {

	log.Println("Start !")

	//d := db.InitDb()
	log.Println("Db initialised !")

	w := watcher.Init()
	log.Println("Watcher initialised !")

	collector.Init()
	log.Println("Collector initialised !")

	web.Start(w)
}
