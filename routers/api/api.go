package api

import (
	"errors"

	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
)

func verificationAuth(ctx *macaron.Context, auth *auth.Auth, sess session.Store) error {
	if !auth.IsGranted("admin.users", sess) {
		ctx.Data["message_categorie"] = "negative"
		ctx.Data["message_icon"] = "warning sign"
		ctx.Data["message_header"] = "Access forbidden"
		ctx.Data["message_text"] = "It's seem you don't have the right to be there"
		ctx.HTML(403, "other/message")
		return errors.New("Not allowed")
	}
	return nil
}
