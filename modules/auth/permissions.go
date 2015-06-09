package auth

import "github.com/mikespook/gorbac"

func initRbac() *gorbac.Rbac {
	//Roles //TODO
	rbac := gorbac.New()
	rbac.Set("user", []string{
		"open.equipement", /* Can see any equipement */
		"open.dashboard",  /* Can home dashboard */
		"api.graph.ping",  /* Can graph ping time */
	}, nil)
	rbac.Set("admin", []string{
		"add.equipement",         /* Add equipement to monitor */
		"del.equipement",         /* Remove equipement to monitor */
		"add.user",               /* Add user */
		"del.user",               /* Remove user */
		"admin.dashboard",        /* Access admin dashboard */
		"admin.equipements",      /* Access admin equipements part */
		"admin.users",            /* Access admin users part */
		"api.network.ping",       /* Can do ping */
		"api.network.snmp",       /* Can do snmp test */
		"api.network.reversedns", /* Can do reversedns */
	}, []string{"user"})
	rbac.Set("master", []string{}, []string{"admin"})
	return rbac
}

//GetRoles return the list of roles defined in this class
func (this *Auth) GetRoles() []string {
	roles := this.rbac.Dump()
	ret := make([]string, 0, len(roles))
	for role := range roles {
		ret = append(ret, role)
	}
	return ret
}
