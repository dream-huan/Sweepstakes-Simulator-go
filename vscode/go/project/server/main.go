package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"std/go/project/proto"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var u struct {
	uid        int64
	u_name     string
	u_password string
}

func checkp(uid int64, s string) (b bool) {
	sqlStr := "select * from users where uid=?"
	_ = db.QueryRow(sqlStr, uid).Scan(&u.uid, &u.u_name, &u.u_password)
	if u.u_password == s {
		return true
	} else {
		return false
	}
}

func logwrite(err error) {
	file, err := os.OpenFile("./log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logwrite(err)
	}
	log.SetOutput(file)
	log.SetPrefix("[Error]")
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println(err)
}

func basiclogwrite(msg string) {
	file, err := os.OpenFile("./log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logwrite(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println(msg)
}

func initdb() {
	var err error
	dsn := "root:SUIbianla123@@tcp(127.0.0.1:3306)/users?charset=utf8mb4&parseTime=True"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		logwrite(err)
	}
	err = db.Ping()
	if err != nil {
		logwrite(err)
	}
}

func insert(name, p string) (uid int64) {
	sqlStr := "insert into users(u_name,u_password) values(?,?)"
	ret, err := db.Exec(sqlStr, name, p)
	if err != nil {
		logwrite(err)
	}
	uid, err = ret.LastInsertId()
	if err != nil {
		logwrite(err)
	}
	return uid
}

func pnp(uid int64, np string) {
	sqlStr := "update users set u_password=? where uid=?"
	_, err := db.Exec(sqlStr, np, uid)
	if err != nil {
		logwrite(err)
	}
}

func login(msg *string, conn net.Conn) {
	var uid int64
	tu := ""
	tp := ""
	k := 0
	for _, r := range *msg {
		if r == ' ' {
			k += 1
		}
		if k == 1 && r != ' ' {
			tu += string(r)
		}
		if k == 2 && r != ' ' {
			tp += string(r)
		}
	}
	uid, err := strconv.ParseInt(tu, 10, 64)
	if err != nil {
		logwrite(err)
	}
	if checkp(uid, tp) {
		data, err := proto.Encode("true")
		basiclogwrite(fmt.Sprintf("New Login:\nUid:%d\nEnterPassword:%s\nResult:%s\n", uid, tp, "true"))
		if err != nil {
			logwrite(err)
		}
		conn.Write(data)
	} else {
		data, err := proto.Encode("false")
		basiclogwrite(fmt.Sprintf("New Login:\nUid:%d\nEnterPassword:%s\nResult:%s\n", uid, tp, "false"))
		if err != nil {
			logwrite(err)
		}
		conn.Write(data)
	}
}

func register(msg *string, conn net.Conn) {
	name := ""
	p := ""
	k := 0
	for _, r := range *msg {
		if r == ' ' {
			k += 1
		}
		if k == 1 && r != ' ' {
			name += string(r)
		}
		if k == 2 && r != ' ' {
			p += string(r)
		}
	}
	uid := insert(name, p)
	basiclogwrite(fmt.Sprintf("New Register:\nUid:%d\nName:%s\nPassword:%s\n", uid, name, p))
	data, err := proto.Encode(strconv.FormatInt(uid, 10))
	if err != nil {
		logwrite(err)
	}
	conn.Write(data)
}

func changep(msg *string, conn net.Conn) {
	var uid int64
	tu := ""
	tp := ""
	np := ""
	k := 0
	for _, r := range *msg {
		if r == ' ' {
			k += 1
		}
		if k == 1 && r != ' ' {
			tu += string(r)
		}
		if k == 2 && r != ' ' {
			tp += string(r)
		}
		if k == 3 && r != ' ' {
			np += string(r)
		}
	}
	uid, err := strconv.ParseInt(tu, 10, 64)
	if err != nil {
		logwrite(err)
	}
	if checkp(uid, tp) {
		pnp(uid, np)
		basiclogwrite(fmt.Sprintf("New Changep:\nUid:%d\nOldp:%s\nNewp:%s\nResult:%s\n", uid, tp, np, "true"))
		data, err := proto.Encode("true")
		if err != nil {
			logwrite(err)
		}
		conn.Write(data)
	} else {
		data, err := proto.Encode("false")
		basiclogwrite(fmt.Sprintf("New Changep:\nUid:%d\nOldp:%s\nNewp:%s\nResult:%s\n", uid, tp, np, "false"))
		if err != nil {
			logwrite(err)
		}
		conn.Write(data)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			logwrite(err)
		}
		temp := ""
		for i := 0; i < len(msg); i++ {
			if msg[i] == ' ' {
				break
			} else {
				temp += string(msg[i])
			}
		}
		switch temp {
		case "login":
			login(&msg, conn)
		case "register":
			register(&msg, conn)
		case "changep":
			changep(&msg, conn)
		default:
			break
		}
		// fmt.Println("收到client发来的数据：", msg)
	}
}

func main() {
	initdb()
	listen, err := net.Listen("tcp", "10.0.0.4:30000")
	if err != nil {
		logwrite(err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			logwrite(err)
			continue
		}
		fmt.Printf("Copyright ©2021 dreamxw.com All Right Reserved Powered by Azure")
		fmt.Printf("The Service is starting...")
		go process(conn)
	}
}
