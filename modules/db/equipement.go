package db

import (
	"log"
	"time"
)

// Equipement describe a Equipement in database
type Equipement struct {
	ID       uint64 `orm:"auto;pk"`
	IP       string `orm:"unique"`
	Hostname string
	Type     uint
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

//GetEquipements return the list of Equipement in database
func (db *Db) GetEquipements() (int64, []*Equipement) {
	var equipements []*Equipement
	num, err := (*db.Orm).QueryTable("equipement").Limit(-1).All(&equipements)
	log.Printf("Returned Rows Num: %s, %s", num, err)
	return num, equipements
}

//GetEquipementTypes return the list of types possible for Equipement
func (db *Db) GetEquipementTypes() []string {
	return []string{
		"Router",
		"Switch",
		"Server",
	}
}
