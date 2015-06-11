package db

import (
	"log"
	"time"
)

// Equipement describe a Equipement in database
type Equipement struct {
	ID           uint64 `orm:"auto;pk"`
	IP           string `orm:"unique"`
	Hostname     string
	Type         int
	Data         string
	Created      time.Time `orm:"auto_now_add;type(datetime)"`
	Updated      time.Time `orm:"auto_now;type(datetime)"`
	LastActivity time.Time `orm:"auto_now;type(datetime)"`
}

//The name of the machine should be determine from the master domain given in config

/*
//NbEquipements return the number of Equipement in database
func NbEquipements() int64 {
	cnt, _ := (*db.Orm).QueryTable("equipement").Count()
	return cnt
}
*/

// GetEquipement return Equipement by param
func GetEquipement(equi Equipement) (*Equipement, error) {
	return &equi, (*db.Orm).Read(&equi)
}

// GetEquipementbyIP return Equipement by param
func GetEquipementbyIP(equi Equipement) (*Equipement, error) {
	return &equi, (*db.Orm).Read(&equi, "IP")
}

//GetEquipements return the list of Equipement in database
func GetEquipements() (int64, []*Equipement) {
	var equipements []*Equipement
	num, err := (*db.Orm).QueryTable("equipement").Limit(-1).All(&equipements)
	log.Printf("Returned Rows Num: %s, %s", num, err)
	return num, equipements
}

// DelEquipement remove the Equipement pass in param
func DelEquipement(equi *Equipement) error {
	_, err := (*db.Orm).Delete(equi)
	return err
}

// AddEquipement add AddEquipement by param
func AddEquipement(equi *Equipement) (int64, error) {
	return (*db.Orm).Insert(equi)
}
