package daoru

import (
	"database/sql"
	"fmt"

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
		fmt.Printf("DateBase Open Error:%v", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("DateBase Access Error:%v", err)
	}
}

func Insert(name, p string) (uid int64) {
	sqlStr := "insert into users(u_name,u_password) values(?,?)"
	ret, _ := db.Exec(sqlStr, name, p)
	uid, _ = ret.LastInsertId()
	return uid
}

func Pnp(uid int64, np string) {
	sqlStr := "update users set u_password=? where uid=?"
	_, _ = db.Exec(sqlStr, np, uid)
}
