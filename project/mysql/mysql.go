package mysql

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var file1, _ = os.OpenFile("../log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

var u struct {
	uid        int64
	u_name     string
	u_password string
	pool       int
}

func Checkp(uid int64, s string) (b bool) {
	sqlStr := "select * from users where uid=?"
	_ = db.QueryRow(sqlStr, uid).Scan(&u.uid, &u.u_name, &u.u_password, &u.pool)
	if u.u_password == s {
		return true
	} else {
		return false
	}
}

func Checkpool(uid int64) int {
	sqlStr := "select pool from users where uid=?"
	_ = db.QueryRow(sqlStr, uid).Scan(&u.pool)
	return u.pool
}

func init() {
	var err error
	dsn := "root:SUIbianla123@@tcp(127.0.0.1:3306)/users?charset=utf8mb4&parseTime=True"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
	err = db.Ping()
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
}

func Insert(name, p string) (uid int64) {
	sqlStr := "insert into users(u_name,u_password) values(?,?)"
	ret, err := db.Exec(sqlStr, name, p)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
	uid, err = ret.LastInsertId()
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
	return uid
}

func Pnp(uid int64, np string) {
	sqlStr := "update users set u_password=? where uid=?"
	_, err := db.Exec(sqlStr, np, uid)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
}

func Toggle(uid int64, pool int) {
	sqlStr := "update users set pool=? where uid=?"
	_, err := db.Exec(sqlStr, pool, uid)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
}
