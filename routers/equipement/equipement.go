package equipement

import (
	//"github.com/sapk/GoWatch/modules/auth"
	"strconv"

	"github.com/Unknwon/macaron"
	"github.com/sapk/GoWatch/models/equipement"
	//"github.com/macaron-contrib/session"
)

//Dashboard render the dashboard of equipements
func Dashboard(ctx *macaron.Context) {
	ctx.Data["equipements_count"], ctx.Data["Equipements"] = equipement.GetAll()
	ctx.Data["EquipementTypes"] = equipement.Types
	ctx.HTML(200, "equipement/dashboard")
}

//View render the view of a equipement
func View(ctx *macaron.Context) {
	id, _ := strconv.ParseUint(ctx.Params(":id"), 10, 64)
	eq, ok := equipement.GetByID(id)

	if !ok {
		ctx.Data["message_categorie"] = "negative"
		ctx.Data["message_icon"] = "server"
		ctx.Data["message_header"] = "Equipement not found !"
		ctx.Data["message_text"] = "The equipement seems to not be in the database"
		ctx.Data["message_redirect"] = "/admin/equipements"
		ctx.HTML(200, "other/message")
		return
	}

	ctx.Data["equipement"] = eq
	ctx.HTML(200, "equipement/view")
}
