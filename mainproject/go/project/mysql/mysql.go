package mysql

import (
	"database/sql"
	"std/go/project/gologger"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var u struct {
	uid        int64
	u_name     string
	u_password string
}

func Checkp(uid int64, s string) (b bool) {
	sqlStr := "select * from users where uid=?"
	_ = db.QueryRow(sqlStr, uid).Scan(&u.uid, &u.u_name, &u.u_password)
	if u.u_password == s {
		return true
	} else {
		return false
	}
}

func init() {
	var err error
	dsn := "root:SUIbianla123@@tcp(127.0.0.1:3306)/users?charset=utf8mb4&parseTime=True"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		gologger.Logwrite(err)
	}
	err = db.Ping()
	if err != nil {
		gologger.Logwrite(err)
	}
}

func Insert(name, p string) (uid int64) {
	sqlStr := "insert into users(u_name,u_password) values(?,?)"
	ret, err := db.Exec(sqlStr, name, p)
	if err != nil {
		gologger.Logwrite(err)
	}
	uid, err = ret.LastInsertId()
	if err != nil {
		gologger.Logwrite(err)
	}
	return uid
}

func Pnp(uid int64, np string) {
	sqlStr := "update users set u_password=? where uid=?"
	_, err := db.Exec(sqlStr, np, uid)
	if err != nil {
		gologger.Logwrite(err)
	}
}