package admin

import (
	"github.com/Unknwon/macaron"
)

func Dashboard(ctx *macaron.Context) {
	ctx.Data["admin_dashboard"] = true
	ctx.HTML(200, "admin/dashboard")
}
