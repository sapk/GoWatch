package db

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
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

//NbUsers return the number of user in database
func (db *Db) NbUsers() int64 {
	cnt, _ := (*db.Orm).QueryTable("user").Count() // SELECT COUNT(*) FROM USER
	return cnt
}

//GetUsers return the list of user in database
func (db *Db) GetUsers() (int64, []*User) {
	var users []*User
	num, err := (*db.Orm).QueryTable("user").Limit(-1).All(&users)
	log.Printf("Returned Rows Num: %s, %s", num, err)
	return num, users
}

//ContainMaster verifiy if the master is in db (init)
func (db *Db) ContainMaster() bool {
	user := User{ID: 1}

	err := (*db.Orm).Read(&user)

	if err == orm.ErrNoRows {
		return false
	} else if err != nil {
		log.Printf("Error : %v", err)
		return false
	} else {
		return true
	}
}

// GetUser return user by param
func (db *Db) GetUser(user User) (*User, error) {
	return &user, (*db.Orm).Read(&user)
}

// DelUser remove the user pass in param
func (db *Db) DelUser(user *User) error {
	_, err := (*db.Orm).Delete(user)
	return err
}

// Delete remove the user
func (user *User) Delete() error {
	return db.DelUser(user)
}

// CreateUser verify the data and add a user
func (db *Db) CreateUser(username, password, email, role string, autorizedroles []string) error {
	if ok, _ := regexp.MatchString("[a-z0-9_-]{3,16}", username); !ok {
		return errors.New("Username invalid !")
	}
	if ok, _ := regexp.MatchString("[a-z0-9_-]{6,18}", password); !ok {
		return errors.New("Password invalid !")
	}
	if ok, _ := regexp.MatchString("([a-z0-9_.-]+)@([a-z.-]+).([a-z]{2,6})", email); !ok {
		return errors.New("Email invalid !")
	}
	//*
	log.Printf("Role : %v", role)
	log.Printf("AutorizedRoles : %v", autorizedroles)
	if !stringInSlice(role, autorizedroles) {
		return errors.New("Role invalid !")
	}
	//*/
	//TODO check roles and email regex
	//TODo checkc if user exist (normaly done by the orm)
	user := &User{Username: username, Roles: role, Email: email}
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return err
	}

	user.Password = string(pass)
	_, err = (*db.Orm).Insert(user)
	if err != nil {
		return err
	}
	log.Printf("User %s created !", username)
	return nil
}

//GetGravatar return the url of gravatar img form the email of the User.
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
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
