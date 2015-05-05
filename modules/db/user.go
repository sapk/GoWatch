package db

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"time"
)

type User struct {
	Id       uint64 `orm:"auto;pk"`
	Username string `orm:"unique"`
	Password string
	Roles    string
	Email    string    `orm:"unique"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

func (this *Db) NbUsers() int64 {
	cnt, _ := (*this.Orm).QueryTable("user").Count() // SELECT COUNT(*) FROM USER
	return cnt
}
func (this *Db) ContainMaster() bool {
	user := User{Id: 1}

	err := (*this.Orm).Read(&user)

	if err == orm.ErrNoRows {
		return false
	} else if err != nil {
		log.Printf("Error : %v", err)
		return false
	} else {
		return true
	}
}
func (this *Db) CreateUser(username, password, email, roles string) error {
	if ok, _ := regexp.Match("/^[a-z0-9_-]{3,16}$/", []byte(username)); !ok {
		return errors.New("Username invalid !")
	}
	if ok, _ := regexp.Match("/^[a-z0-9_-]{6,18}$/", []byte(password)); !ok {
		return errors.New("Password invalid !")
	}
	if ok, _ := regexp.Match("/^([a-z0-9_.-]+)@([a-z.-]+).([a-z]{2,6})$/", []byte(email)); !ok {
		return errors.New("Email invalid !")
	}
	//TODO check roles adn email regex
	user := &User{Username: username, Roles: roles, Email: email}
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return err
	}

	user.Password = string(pass)
	_, err = (*this.Orm).Insert(user)
	if err != nil {
		return err
	}
	log.Printf("User %s created !", username)
	return nil
}

/*
// multiple fields unique key
func (u *User) TableUnique() [][]string {
	return [][]string{
		[]string{"Username", "Email"},
	}
}
*/
