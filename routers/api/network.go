package api

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/api/network"
	"github.com/sapk/GoWatch/modules/auth"
)

// Ping ping the ip or hostname
func Ping(ctx *macaron.Context, auth *auth.Auth, sess session.Store) {
	if err := auth.VerificationAuth(ctx, sess, []string{"api.network.ping"}); err != nil {
		return
	}
	hostorip := ctx.Req.URL.RawQuery
	ctx.JSON(200, network.Ping(hostorip))
}

// SNMPTest the snmp service of the ip or hostname
func SNMPTest(ctx *macaron.Context, auth *auth.Auth, sess session.Store) {
	if err := auth.VerificationAuth(ctx, sess, []string{"api.network.snmp"}); err != nil {
		return
	}
	hostorip := ctx.Query("host")
	community := ctx.Query("community")
	ctx.JSON(200, network.SNMPTest(hostorip, community))
}

// ReverseDNS reverse DNS the ip
func ReverseDNS(ctx *macaron.Context, auth *auth.Auth, sess session.Store) {
	if err := auth.VerificationAuth(ctx, sess, []string{"api.network.reversedns"}); err != nil {
		return
	}
	ctx.JSON(200, network.ReverseDNS(ctx.Req.URL.RawQuery))
}
