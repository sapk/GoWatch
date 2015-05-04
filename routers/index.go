package routers

import (
	"github.com/Unknwon/macaron"
)

func Index(ctx *macaron.Context) {
	ctx.HTML(200, "index") // 200 is the response code.
}
