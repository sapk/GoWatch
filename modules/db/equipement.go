package db

import (
	"log"
	"time"
	"regexp"
	"errors"
	"strings"
	"github.com/sapk/GoWatch/modules/tools"
)


// Equipement describe a Equipement in database
type Equipement struct {
	ID       uint64 `orm:"auto;pk"`
	IP       string `orm:"unique"`
	Hostname string
	Type     int
	Data     string
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

//The name of the machine should be determine from the master domain given in config

//NbEquipements return the number of Equipement in database
func (db *Db) NbEquipements() int64 {
	cnt, _ := (*db.Orm).QueryTable("equipement").Count()
	return cnt
}
// GetEquipement return Equipement by param
func (db *Db) GetEquipement(equi Equipement) (*Equipement, error) {
    return &equi, (*db.Orm).Read(&equi)
}
 // GetEquipement return Equipement by param
 func (db *Db) GetEquipementbyIP(equi Equipement) (*Equipement, error) {
     return &equi, (*db.Orm).Read(&equi,"IP")
 }
//GetEquipements return the list of Equipement in database
func (db *Db) GetEquipements() (int64, []*Equipement) {
	var equipements []*Equipement
	num, err := (*db.Orm).QueryTable("equipement").Limit(-1).All(&equipements)
	log.Printf("Returned Rows Num: %s, %s", num, err)
	return num, equipements
}
// DelEquipement remove the Equipement pass in param
func (db *Db) DelEquipement(equi *Equipement) error {
    _, err := (*db.Orm).Delete(equi)
    return err
}
// Delete remove the Equipement
func (equi *Equipement) Delete() error {
    return db.DelEquipement(equi)
}
//GetEquipementTypes return the list of types possible for Equipement
func (db *Db) GetEquipementTypes() []string {
	return []string{
		"Router",
		"Switch",
		"Server",
                "Computer",
	}
}

func  posInSlice(slice []string,value string) int {
    for p, v := range slice {
        if (v == value) {
            return p
        }
    }
    return -1
}
// CreateUser verify the data and add a user
func (db *Db) CreateEquipement(ip, host, typ string) error {
    log.Println("CreateEquipement : ", ip,host, typ)
    if ok, _ := regexp.MatchString(tools.ValidIpAddressRegex, ip); !ok {
        return errors.New("IP invalid !")
    }
    if ok, _ := regexp.MatchString(tools.ValidHostAddressRegex, host); !ok {
        return errors.New("Host invalid !")
    }
    //*
    log.Printf("Type : %v", typ)
    log.Printf("AutorizedTypes : %v", db.GetEquipementTypes())
     
    idType := -1
    if idType = posInSlice(db.GetEquipementTypes(),typ);idType == -1 {
        return errors.New("Type invalid !")
    }
    //*/
    //TODo check if equipement exist (normaly done by the orm)
    equi := &Equipement{IP: ip, Hostname: host, Type: idType}
    
    _, err := (*db.Orm).Insert(equi)
    if err != nil {
        return err
    }
    log.Printf("Equi %s created !", equi.Hostname)
    return nil
}


//GetTypeIcon return the class of icon for the Equipement
func (equi *Equipement) GetTypeIcon() string {
        return strings.ToLower(db.GetEquipementTypes()[equi.Type])
}

//UpdatedFormated 
func (equi *Equipement) UpdatedFormated() string {
        sec := int(time.Since(equi.Updated).Seconds())
        return (time.Duration(sec)*time.Second).String()
}
//Update 
func (equi *Equipement) Update() error {
        _, err := (*db.Orm).Update(equi, "Data", "Updated")
        return err
}