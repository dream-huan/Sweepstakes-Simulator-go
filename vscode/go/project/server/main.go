package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"net"
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

func initdb() {
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

func insert(name, p string) (uid int64) {
	sqlStr := "insert into users(u_name,u_password) values(?,?)"
	ret, _ := db.Exec(sqlStr, name, p)
	uid, _ = ret.LastInsertId()
	return uid
}

func pnp(uid int64, np string) {
	sqlStr := "update users set u_password=? where uid=?"
	_, _ = db.Exec(sqlStr, np, uid)
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
	uid, _ = strconv.ParseInt(tu, 10, 64)
	fmt.Printf("%v new login:\nuid:%d\nparsep:%s\n", time.Now(), uid, tp)
	if checkp(uid, tp) {
		data, _ := proto.Encode("true")
		conn.Write(data)
	} else {
		data, _ := proto.Encode("false")
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
	fmt.Printf("%v new register:\nuid:%d\nname:%s\npassword:%s\n", time.Now(), uid, name, p)
	data, _ := proto.Encode(strconv.FormatInt(uid, 10))
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
	uid, _ = strconv.ParseInt(tu, 10, 64)
	fmt.Printf("%v new changep:\nuid:%d\noldp:%s\nnewp:%s\n", time.Now(), uid, tp, np)
	if checkp(uid, tp) {
		pnp(uid, np)
		data, _ := proto.Encode("true")
		conn.Write(data)
	} else {
		data, _ := proto.Encode("false")
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
			fmt.Println("decode msg failed, err:", err)
			return
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
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		fmt.Printf("Copyright ©2021 dreamxw.com All Right Reserved Powered by Azure")
		fmt.Printf("The Service is starting...")
		go process(conn)
	}
}
