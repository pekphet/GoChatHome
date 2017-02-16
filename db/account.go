package db

import (
	"log"
	"crypto/md5"
	"io"
	"strconv"
	"time"
	"fmt"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func insertAcc(db *sql.DB, name string, acc string, pwd string) bool {
	stmt, err := db.Prepare("INSERT account SET acc=?,pwd=?")
	if err != nil {
		log.Fatal(err)
		return false
	}
	res, a := stmt.Exec(acc, pwd)

	if a != nil {
		if a.(*mysql.MySQLError).Number == 1062 {
			return false
		}
		log.Fatal(a)
		return false
	}
	affect, err2 := res.RowsAffected()
	if err2 != nil {
		log.Fatal(err2)
		return false
	}
	log.Println(affect)
	row := db.QueryRow("SELECT _id FROM account WHERE acc=?", acc)
	var _id int
	row.Scan(&_id)
	insertUser(db, _id, name)
	insertToken(db, _id, "")
	return true
}

func login(db *sql.DB, acc string, pwd string) (uid int, name string, token string, ok bool){
	rows, err := db.Query("SELECT _id FROM account WHERE acc=? AND pwd=?", acc, pwd)
	defer rows.Close()
	if (err != nil) {
		log.Fatal(err)
		return 0, "", "", false
	}
	if !rows.Next() {
		return 0, "", "", false
	}
	var _id int
	rows.Scan(&_id)
	ok = true
	uid = _id
	err2 := db.QueryRow("SELECT nick FROM usr WHERE uid=?", _id).Scan(&name)
	if (err2 != nil) {
		log.Fatal(err2)
		return 0, "", "", false
	}
	token = generateToken(uid)
	stmtToken, err3 := db.Prepare("UPDATE token SET token=? WHERE uid=?")
	if (err3 != nil) {
		log.Fatal(err3)
		return 0, "", "", false
	}
	stmtToken.Exec(token, uid)
	return
}

func insertUser(db *sql.DB, uid int, name string) {
	stmt, err := db.Prepare("INSERT usr SET uid=?,nick=?")
	if (err != nil) {
		log.Fatal(err)
		return
	}
	stmt.Exec(uid, name)
}

func insertToken(db *sql.DB, uid int, token string) {
	stmt, err := db.Prepare("INSERT token SET uid=?,token=?")
	if (err != nil) {
		log.Fatal(err)
		return
	}
	stmt.Exec(uid, token)
}

func generateToken(uid int) string{
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}