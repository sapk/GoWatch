package admin

import (
	"github.com/Unknwon/macaron"
)

func Users(ctx *macaron.Context) {
	ctx.Data["admin_users"] = true
	ctx.HTML(200, "admin/users")
}
