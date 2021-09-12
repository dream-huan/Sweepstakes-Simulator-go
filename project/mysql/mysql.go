package mysql

import (
	"database/sql"
	"log"
	"os"
	"strconv"
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
	fivetimes  int
	fourtimes  int
	stone      int
	addstone   int
	effective  bool
}

func Checkp(uid int64, s string) (b bool) {
	sqlStr := "select * from users where uid=?"
	_ = db.QueryRow(sqlStr, uid).Scan(&u.uid, &u.u_name, &u.u_password, &u.pool, &u.fivetimes, &u.fourtimes, &u.stone)
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

func Getstone(uid int64) int {
	sqlStr := "select stone from users where uid=?"
	_ = db.QueryRow(sqlStr, uid).Scan(&u.stone)
	return u.stone
}

func init() {
	var err error
	dsn := "root:SUIbianla123@@tcp(127.0.0.1:3306)/users"
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

func Insertprop(uid int64, np string, time string) {
	sqlStr := "insert into data(uid,prop,time) values(?,?,?)"
	_, err := db.Exec(sqlStr, uid, np, time)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
}

func Changestone(newstone int, uid int64) {
	sqlStr := "update users set stone=? where uid=?"
	_, err := db.Exec(sqlStr, newstone, uid)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
}

func Checkresult(uid int64) {
	sqlStr := "select uid,prop,time from data where uid=? into outfile '/workspaces/go-chat/project/" + strconv.FormatInt(uid, 10) + ".txt'"
	_, err := db.Exec(sqlStr, uid)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
}

func Recharge(key string, uid int64) bool {
	sqlStr := "select stone,effective from recharge where keycode=?"
	err := db.QueryRow(sqlStr, key).Scan(&u.addstone, &u.effective)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
	if u.effective == false {
		return false
	}
	u.stone = Getstone(uid)
	u.stone += u.addstone
	sqlStr = "update users set stone=? where uid=?"
	_, err = db.Exec(sqlStr, u.stone, uid)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
	sqlStr = "update recharge set effective=0 where keycode=?"
	_, err = db.Exec(sqlStr, key)
	return u.effective
}

func Checkstatistics(uid int64) (int, int) {
	sqlStr := "select fivetimes,fourtimes from users where uid=?"
	_ = db.QueryRow(sqlStr, uid).Scan(&u.fivetimes, &u.fourtimes)
	return u.fivetimes, u.fourtimes
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

func Timeschange(uid int64, fivetimes int, fourtimes int) {
	sqlStr := "update users set fivetimes=? where uid=?"
	_, err := db.Exec(sqlStr, fivetimes, uid)
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
	sqlStr = "update users set fourtimes=? where uid=?"
	_, err = db.Exec(sqlStr, fourtimes, uid)
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
