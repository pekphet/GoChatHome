package db

import (
	"database/sql"
	"log"
)

var gDB *sql.DB

func init(){
	db, err := sql.Open("mysql", "dev:@Dev1991@tcp(localhost:53306)/go_android_im")
	if (err != nil) {
		log.Fatal(err)
	}
	gDB = db
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