package admin

import (
	"github.com/Unknwon/macaron"
)

func Observeds(ctx *macaron.Context) {
	ctx.Data["admin_observeds"] = true
	ctx.HTML(200, "admin/observeds")
}
