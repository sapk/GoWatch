package api

import (
	"log"
	"time"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/api/network"
	"github.com/sapk/GoWatch/modules/auth"
)

// Ping ping the ip or hostname
func Ping(ctx *macaron.Context, auth *auth.Auth, sess session.Store) {
	if err := verificationAuth(ctx, auth, sess); err != nil {
		return
	}
	//TODO
	ip := ctx.Query("ip")
	timer := time.AfterFunc(3*time.Second, func() {
		log.Printf("Time out!")
		ctx.JSON(200, network.PingResponse{IP: ip, Result: false, Time: 3 * time.Second})
		ctx.Resp.Flush()
		ctx.Next()
		log.Printf("Written ? %v", ctx.Written())
	})
	defer timer.Stop()
	rep := network.Ping(ip)
	ctx.JSON(200, rep)
	log.Printf("Clearing timer !")
}
