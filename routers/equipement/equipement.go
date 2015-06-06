package equipement

import (
	//"github.com/sapk/GoWatch/modules/auth"
	"strconv"

	"github.com/Unknwon/macaron"
	"github.com/sapk/GoWatch/modules/db"
	//"github.com/macaron-contrib/session"
)

//Dashboard render the dashboard of equipements
func Dashboard(ctx *macaron.Context, dbb *db.Db) {
	ctx.Data["equipements_count"], ctx.Data["Equipements"] = dbb.GetEquipements()
	ctx.Data["EquipementTypes"] = dbb.GetEquipementTypes()
	ctx.HTML(200, "equipement/dashboard")
}

//View render the view of a equipement
func View(ctx *macaron.Context, dbb *db.Db) {
	id, _ := strconv.ParseUint(ctx.Params(":id"), 10, 64)
	eq, err := dbb.GetEquipement(db.Equipement{ID: id})

	if err != nil {
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
