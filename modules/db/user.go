package db

import (
	"log"
	"time"
)

// User describe a user in database
type User struct {
	ID       uint64 `orm:"auto;pk"`
	Username string `orm:"unique"`
	Password string
	Roles    string
	Email    string    `orm:"unique"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

/*
//NbUsers return the number of user in database
func NbUsers() int64 {
	cnt, _ := (*db.Orm).QueryTable("user").Count() // SELECT COUNT(*) FROM USER
	return cnt
}
*/
//GetUsers return the list of user in database
func GetUsers() (int64, []*User) {
	var users []*User
	num, err := (*db.Orm).QueryTable("user").Limit(-1).All(&users)
	log.Printf("Returned Rows Num: %s, %s", num, err)
	return num, users
}

// GetUser return user by param
func GetUser(user *User) (*User, error) {
	return user, (*db.Orm).Read(user)
}

// AddUser return user by param
func AddUser(user *User) (int64, error) {
	return (*db.Orm).Insert(user)
}

// DelUser remove the user pass in param
func DelUser(user *User) error {
	_, err := (*db.Orm).Delete(user)
	return err
}
