package auth

import "github.com/mikespook/gorbac"

func initRbac() *gorbac.Rbac {
	//Roles //TODO
	rbac := gorbac.New()
	rbac.Set("user", []string{
		"open.equipement", /* Can see any equipement */
		"open.dashboard",  /* Can home dashboard */
	}, nil)
	rbac.Set("admin", []string{
		"add.equipement",    /* Add equipement to monitor */
		"del.equipement",    /* Remove equipement to monitor */
		"add.user",          /* Add user */
		"del.user",          /* Remove user */
		"admin.dashboard",   /* Access admin dashboard */
		"admin.equipements", /* Access admin equipements part */
		"admin.users",       /* Access admin users part */
	}, []string{"user"})
	rbac.Set("master", []string{}, []string{"admin"})
	return rbac
}
