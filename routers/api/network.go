package api

import (
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
	ctx.JSON(200, network.Ping("192.168.1.0"))
}
