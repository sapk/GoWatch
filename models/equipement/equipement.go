package equipement

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/sapk/GoWatch/modules/db"
	"github.com/sapk/GoWatch/modules/tools"
)

// Equipement describe a Equipement in database
type Equipement struct {
	dbPointer *db.Equipement
}

// Equipements describe a list of Equipement
type Equipements map[uint64]Equipement

var listEquipement Equipements
var isInit = false

//Types all the allowed
var Types = []string{
	"Router",
	"Switch",
	"Server",
	"Firewall",
	"WLAN",
	"Computer",
	"EndPoint",
	"Other",
}

//init
func init() {
	if !isInit {
		log.Println("Initialisation of Equipement model...")
		listEquipement = make(Equipements)
		_, equipements := db.GetEquipements()
		for _, eq := range equipements {
			listEquipement[eq.ID] = Equipement{eq}
		}
		log.Println("Equipements find :", NbEquipements())
		isInit = true
	}
}

//NbEquipements return the number of Equipement in database
func NbEquipements() int {
	return len(listEquipement)
}

//GetAll return listEquipement
func GetAll() (int, *Equipements) {
	return NbEquipements(), &listEquipement
}

//GetByID return user of id
func GetByID(id uint64) (*Equipement, bool) {
	user, ok := listEquipement[id]
	return &user, ok
}

//GetByIP return user of id
func GetByIP(ip string) (*Equipement, bool) {
	for _, el := range listEquipement {
		if el.dbPointer.IP == ip {
			return &el, true
		}
	}
	return nil, false
}

// CreateEquipement verify the data and add a user
func CreateEquipement(ip, host, typ string) error {
	log.Println("CreateEquipement : ", ip, host, typ)
	if ok, _ := regexp.MatchString(tools.ValidIPAddressRegex, ip); !ok {
		return errors.New("IP invalid")
	}
	if ok, _ := regexp.MatchString(tools.ValidHostAddressRegex, host); !ok {
		return errors.New("Host invalid")
	}
	//*
	log.Printf("Type : %v", typ)
	log.Printf("AutorizedTypes : %v", Types)

	idType := -1
	if idType = tools.PosInSlice(Types, typ); idType == -1 {
		return errors.New("Type invalid")
	}
	//*/
	//TODo check if equipement exist (normaly done by the orm)
	equi := &db.Equipement{IP: ip, Hostname: host, Type: idType}

	//	_, err := (*db.Orm).Insert(equi)
	_, err := db.AddEquipement(equi)

	if err != nil {
		return err
	}

	listEquipement[equi.ID] = Equipement{equi}
	log.Printf("Equi %s created !", equi.Hostname)
	return nil
}

//ID return the ID of the Equipement
func (eq Equipement) ID() uint64 {
	return eq.dbPointer.ID
}

//Hostname return the Hostname of the Equipement
func (eq Equipement) Hostname() string {
	return eq.dbPointer.Hostname
}

//IP return the IP of the Equipement
func (eq Equipement) IP() string {
	return eq.dbPointer.IP
}

//Created return the Created of the Equipement
func (eq Equipement) Created() time.Time {
	return eq.dbPointer.Created
}

//Updated return the Updated of the Equipement
func (eq Equipement) Updated() time.Time {
	return eq.dbPointer.Updated
}

//LastActivity return the LastActivity of the Equipement
func (eq Equipement) LastActivity() time.Time {
	return eq.dbPointer.LastActivity
}

//Data return the Data of the Equipement
func (eq Equipement) Data() string {
	return eq.dbPointer.Data
}

//GetTypeIcon return the class of icon for the Equipement
func (eq Equipement) GetTypeIcon() string {
	return strings.ToLower(Types[eq.dbPointer.Type])
}

//UpdatedFormated return the update date formated
func (eq Equipement) UpdatedFormated() string {
	sec := int(time.Since(eq.dbPointer.Updated).Seconds())
	return (time.Duration(sec) * time.Second).String()
}

//LastActivityFormated return the update date formated
func (eq Equipement) LastActivityFormated() string {
	sec := int(time.Since(eq.dbPointer.LastActivity).Seconds())
	return (time.Duration(sec) * time.Second).String()
}

//UpdateActivity update the equipement in database
func (eq Equipement) UpdateActivity() error {
	//_, err := (*db.Orm).Update(equi, "Data", "Updated")
	//_, err := (*db.Orm).Update(eq, "Updated")
	//TODO save some time
	eq.dbPointer.LastActivity = time.Now()
	return nil
}

//Delete delete a Equipement in database
func (eq Equipement) Delete() error {
	//TODO
	id := eq.ID()
	err := db.DelEquipement(eq.dbPointer)
	if err != nil {
		return err
	}
	delete(listEquipement, id)
	return nil
}
