package admin

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/csrf"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/db"
        "github.com/sapk/GoWatch/modules/tools"
        "net"
        "log"
        "regexp"
        "strings"
        "strconv"
        "html/template"
)

// Equipements generate the admin page for Equipement management
func Equipements(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {

	if err := verificationAuth(ctx, auth, sess); err != nil {
		return
	}
	fillGlobalPage(ctx, db, "admin_equipements")
	ctx.Data["equipements_count"], ctx.Data["Equipements"] = db.GetEquipements()
	ctx.Data["EquipementTypes"] = db.GetEquipementTypes()
	ctx.HTML(200, "admin/equipements")
	//TODO representation in tmeplate
}

// EquipementDel handle deletion of one user
func EquipementDel(ctx *macaron.Context, auth *auth.Auth, sess session.Store, dbb *db.Db, x csrf.CSRF) {
	if err := verificationAuth(ctx, auth, sess); err != nil {
            return
        }
        id, _ := strconv.ParseUint(ctx.Params(":id"), 10, 64)
        equi, err := dbb.GetEquipement(db.Equipement{ID: id})
        if err != nil {
            ctx.Data["message_categorie"] = "negative"
            ctx.Data["message_icon"] = "server"
            ctx.Data["message_header"] = "Equipement not found !"
            ctx.Data["message_text"] = "The equipement seems to not be in the database"
            ctx.Data["message_redirect"] = "/admin/equipements"
            ctx.HTML(200, "other/message")
        } else if !x.ValidToken(ctx.Query("confirm")) {
            // Ask for confirmation
            ctx.Data["message_categorie"] = ""
            ctx.Data["message_icon"] = "server"
            ctx.Data["message_header"] = "Confirm equipement deletion!"
            ctx.Data["csrf_token"] = x.GetToken()
            sess.Set("crsf_user_id", equi.ID)
            ctx.Data["message_text"] = strings.Join([]string{"Do you really want to delete : ", equi.Hostname}, " ")
    
            ctx.HTML(200, "other/confirmation")
        } else {
            // We del the user if all is good
            if sess.Get("crsf_user_id") != equi.ID {
                ctx.Data["message_categorie"] = "negative"
                ctx.Data["message_icon"] = "server"
                ctx.Data["message_header"] = "Hummm ..."
                ctx.Data["message_text"] = template.HTML("It's seem there is a problem with the <a href='fr.wikipedia.org/wiki/Cross-Site_Request_Forgery' target='_blank'>CRSF</a> protection please retry.")
                ctx.Data["message_redirect"] = "/admin/equipements"
                ctx.HTML(200, "other/message")
                return
            }
            err := equi.Delete()
            if err != nil {
                ctx.Data["message_categorie"] = ""
                ctx.Data["message_icon"] = "server"
                ctx.Data["message_header"] = "Oups"
                ctx.Data["message_text"] = "It's seem there is a problem during deletion"
                ctx.Data["message_redirect"] = "/admin/equipements"
                ctx.HTML(200, "other/message")
                return
            }
            ctx.Data["message_categorie"] = ""
            ctx.Data["message_icon"] = "server"
            ctx.Data["message_header"] = "Equipement is deleted !"
            ctx.Data["message_text"] = "The equipement has been deleted from the database."
            ctx.Data["message_redirect"] = "/admin/equipements"
            ctx.HTML(200, "other/message")
            sess.Delete("crsf_user_id")
        }
    
}

// EquipementAdd generate the admin page for adding a user
func EquipementAdd(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	if err := verificationAuth(ctx, auth, sess); err != nil {
		return
	}
	fillGlobalPage(ctx, db, "admin_equipements")
	//TODO ?	ctx.Data["organizations"] = auth.GetOrganizations()
	//TODO ?	ctx.Data["locations"] = auth.GetLocations()
	ctx.Data["types"] = db.GetEquipementTypes()
	ctx.HTML(200, "admin/add_equipement")
}

// EquipementAddPost handle the adding of a user
func EquipementAddPost(ctx *macaron.Context, auth *auth.Auth, sess session.Store, db *db.Db) {
	//TODO add support for adding wildcard dns (find host in DNS) && support ip scan of range
	//TODO convert hostname to IP
	//TODO use config for DNS resolver if it'snt host
	//TODO
	err := verificationAuth(ctx, auth, sess);
	if  err != nil {
            return
        }
        
        ip := ctx.Query("iporhostname")
        host := ctx.Query("iporhostname")
        
        //TODO resolve if iporhostname is HOST
        if ok, _ := regexp.MatchString(tools.ValidIpAddressRegex, ip); !ok {
            //Si ce n'est pas un ip on essaie de résoudrele host
            log.Println("Resolve IP : ", ip)
            var i *net.IPAddr
            i, err = net.ResolveIPAddr("ip", ip)
            ip = i.String()
            if err != nil {
                log.Println("Erreur in resolving : ", err)
            }
        }else { 
                    //TODO reverse DNS if iporhostname is IP
                    log.Println("Reverse DNS : ", host)
                    //Si c'est une ip on essaie de faire un reverse dns
                    hosts, err := net.LookupAddr(host)
                    log.Println("Hosts discover in reverse : ", hosts)
                    if err != nil || len(hosts) == 0 {
                        log.Println("Erreur in resolving : ", err)
                        host = ip //in case of failure we ip has host also
                    }else{
                        host = strings.Trim(hosts[0],".")       
                    }
        }
        
        log.Println("CreateEquipement : ", ip,host, ctx.Query("type"))
        if err == nil {
                err = db.CreateEquipement(ip,host, ctx.Query("type"))
                log.Println("Err : ", err)
        }
        
        log.Println("Err : ", err)
        if err != nil {
            log.Println("Equipement add failed !")
            fillGlobalPage(ctx, db, "admin_equipements")
            //TODO ?    ctx.Data["organizations"] = auth.GetOrganizations()
            //TODO ?    ctx.Data["locations"] = auth.GetLocations()
            ctx.Data["types"] = db.GetEquipementTypes()
            ctx.Data["EquipementAddError"] = true
            ctx.Data["EquipementAddErrorText"] = err.Error()
            ctx.HTML(200, "admin/add_equipement")
            return
        }
    
        ctx.Data["message_categorie"] = "positive"
        ctx.Data["message_icon"] = "add equipement"
        ctx.Data["message_header"] = "Equipement added !"
        ctx.Data["message_text"] = "The equipement has been added to the database."
        ctx.Data["message_redirect"] = "/admin/equipements"
        ctx.HTML(200, "other/message")
}
