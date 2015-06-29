package api

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/api/graph"
	"github.com/sapk/GoWatch/modules/auth"
)

// GraphPing graph ping the data from ip
func GraphPing(ctx *macaron.Context, auth *auth.Auth, sess session.Store) {
	if err := auth.VerificationAuth(ctx, sess, []string{"api.graph.ping"}); err != nil {
		return
	}
	ctx.Header().Set("Expires", "0")
	ctx.Header().Set("Cache-Control", "must-revalidate")
	ctx.Header().Set("Content-Type", "image/png")
	graph.EquipementPing(ctx.Params(":id"), ctx.Params(":duration"), ctx.RW())
}

// GraphPingData send data for ping of ip in database
func GraphPingData(ctx *macaron.Context, auth *auth.Auth, sess session.Store) {
	if err := auth.VerificationAuth(ctx, sess, []string{"api.graph.ping"}); err != nil {
		return
	}
	ctx.Header().Set("Expires", "0")
	ctx.Header().Set("Cache-Control", "must-revalidate")
	ctx.JSON(200, graph.EquipementPingData(ctx.Params(":id"), ctx.Params(":duration")))
}
