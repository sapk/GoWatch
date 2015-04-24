package main

import (
	"github.com/Unknwon/macaron"
	//	"github.com/xyproto/permissions2"
)

func main() {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Gziper())
	m.Use(macaron.Static("public"))
	m.Use(macaron.Renderer())

	m.Use(macaron.Recovery())

	m.Get("/", func(ctx *macaron.Context) {
		ctx.Data["Name"] = "jeremy"
		ctx.HTML(200, "index") // 200 is the response code.
	})

	m.Run()
}
