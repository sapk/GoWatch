package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/sapk/GoWatch/modules/db"
	"github.com/sapk/GoWatch/modules/tools"
)

// User describe a user in database
type User struct {
	dbPointer *db.User
}

// Users describe a list of Equipement
type Users map[uint64]User

var listUser Users
var isInit = false

const ()

//init
func init() {
	if !isInit {
		log.Println("Initialisation of User model...")
		listUser = make(Users)
		_, users := db.GetUsers()
		for _, u := range users {
			listUser[u.ID] = User{u}
		}
		log.Println("Users find :", NbUsers())
		isInit = true
	}
}

//GetAll return listUser
func GetAll() (int, Users) {
	return NbUsers(), listUser
}

//GetByID return user of id
func GetByID(id uint64) (User, bool) {
	user, ok := listUser[id]
	return user, ok
}

//NbUsers return the number of User in database
func NbUsers() int {
	return len(listUser)
}

//ContainMaster verifiy if the master is in db (init)
func ContainMaster() bool {
	_, ok := listUser[1]
	return ok
}

// CreateUser verify the data and add a user
func CreateUser(username, password, email, role string, autorizedroles []string) error {
	if ok, _ := regexp.MatchString(tools.ValidUsernameRegex, username); !ok {
		return errors.New("Username invalid !")
	}
	if ok, _ := regexp.MatchString(tools.ValidPasswordRegex, password); !ok {
		return errors.New("Password invalid !")
	}
	if ok, _ := regexp.MatchString(tools.ValidEmailRegex, email); !ok {
		return errors.New("Email invalid !")
	}
	//*
	log.Printf("Role : %v", role)
	log.Printf("AutorizedRoles : %v", autorizedroles)
	if !tools.StringInSlice(role, autorizedroles) {
		return errors.New("Role invalid !")
	}
	//*/
	//TODO check roles and email regex
	//TODo checkc if user exist (normaly done by the orm)
	user := &db.User{Username: username, Roles: role, Email: email}
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return err
	}

	user.Password = string(pass)
	_, err = db.AddUser(user)
	//_, err = (*db).Orm.Insert(user)
	if err != nil {
		return err
	}

	listUser[user.ID] = User{user}
	log.Printf("User %s created !", username)
	return nil
}

//GetGravatar return the url of gravatar img form the email of the User.
func (u User) GetGravatar() string {
	md5 := md5.Sum([]byte(strings.ToLower(strings.Trim(u.dbPointer.Email, " "))))
	return "http://www.gravatar.com/avatar/" + hex.EncodeToString(md5[:16])
	//. md5( strtolower( trim( $email ) ) ) . "?d=" . urlencode( $default ) . "&s=" . $size;
}

//ID return the ID of User in database
func (u User) ID() uint64 {
	return u.dbPointer.ID
}

//Roles return the ID of User in database
func (u User) Roles() string {
	return u.dbPointer.Roles
}

//Username return the Username of User in database
func (u User) Username() string {
	return u.dbPointer.Username
}

//Email return the Email of User in database
func (u User) Email() string {
	return u.dbPointer.Email
}

//Delete delete a User in database
func (u User) Delete() error {
	//TODO
	id := u.ID()
	err := db.DelUser(u.dbPointer)
	if err != nil {
		return err
	}
	delete(listUser, id)
	return nil
}
