package user

import (
	//"./../../modules/auth"
	"github.com/Unknwon/macaron"
	//"github.com/macaron-contrib/session"
)

func Login(ctx *macaron.Context) {
	ctx.HTML(200, "user/login")
}

func CurrentPage(ctx *macaron.Context) {
	ctx.HTML(200, "user/page")
}

/*
func SignIn(ctx *macaron.Context, authentificator *auth.Auth, sess session.Store) {
	authentificator.SignIn(ctx, sess)
}
*/
