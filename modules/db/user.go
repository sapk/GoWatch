package db

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"strings"
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
func (this *Db) GetUsers() (int64, []*User) {
	var users []*User
	num, err := (*this.Orm).QueryTable("user").Limit(-1).All(&users)
	for u := range users {
		log.Printf("Row : %v", u)
	}
	log.Printf("Returned Rows Num: %s, %s", num, err)
	return num, users
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
	if ok, _ := regexp.MatchString("[a-z0-9_-]{3,16}", username); !ok {
		return errors.New("Username invalid !")
	}
	if ok, _ := regexp.MatchString("[a-z0-9_-]{6,18}", password); !ok {
		return errors.New("Password invalid !")
	}
	if ok, _ := regexp.MatchString("([a-z0-9_.-]+)@([a-z.-]+).([a-z]{2,6})", email); !ok {
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
func (user *User) GetGravatar() string {
	md5 := md5.Sum([]byte(strings.ToLower(strings.Trim(user.Email, " "))))
	return "http://www.gravatar.com/avatar/" + hex.EncodeToString(md5[:16])
	//. md5( strtolower( trim( $email ) ) ) . "?d=" . urlencode( $default ) . "&s=" . $size;
}

/*
// multiple fields unique key
func (u *User) TableUnique() [][]string {
	return [][]string{
		[]string{"Username", "Email"},
	}
}
*/
