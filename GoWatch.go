package main

import (
	"log"

	"github.com/sapk/GoWatch/modules/db"
	"github.com/sapk/GoWatch/modules/web"
)

func main() {

	log.Println("Start !")
	d := db.InitDb()
	log.Println("Db initialised !")

	web.Start(d)
}
