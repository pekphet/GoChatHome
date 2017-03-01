package db

import (
	"database/sql"
	"log"
	"encoding/json"
)

type Equips map[int] *Equip

var gDB 	*sql.DB
var equips	Equips

func init(){
	db, err := sql.Open("mysql", "dev:@Dev1991@tcp(localhost:53306)/go_android_im")
	if (err != nil) {
		log.Fatal(err)
	}
	gDB = db
	equips = make(Equips)
}

func CloseDB() {
	if gDB != nil {
		gDB.Close()
	}
}

func Reg(name string, acc string, pwd string) bool {
	return insertAcc(name, acc, pwd)
}

func Login(acc string, pwd string) (uid int, name string, token string, ok bool){
	return login(acc, pwd)
}

func GetEquip(eid int) string {
	equip, ok := equips[eid]
	if !ok {
		equip = getEquipInfo(eid)
		equips[eid] = equip
	}
	b, _ := json.Marshal(equip)
	return string(b)
}