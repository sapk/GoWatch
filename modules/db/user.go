package db

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Id       int
	Username string `orm:"unique"`
	Password string
	Roles    string
	Email    string    `orm:"unique"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

func (this *Db) CreateUser(username, password, roles, email string) {
	//TODO not silently fail
	user := &User{Username: username, Roles: roles, Email: email}
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 7)
	user.Password = string(pass)
	(*this.Orm).Insert(user)
	log.Printf("User %s created !", username)
}
func (this *Db) CreateAdminUser() {
	this.CreateUser("admin", "0000", "master", "")
}

/*
// multiple fields unique key
func (u *User) TableUnique() [][]string {
	return [][]string{
		[]string{"Username", "Email"},
	}
}
*/
