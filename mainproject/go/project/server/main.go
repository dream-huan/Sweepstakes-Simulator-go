package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"std/go/project/gologger"
	"std/go/project/mysql"
	"std/go/project/proto"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

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
		gologger.Logwrite(err)
	}
	if mysql.Checkp(uid, tp) {
		data, err := proto.Encode("true")
		gologger.BasicLogwrite(fmt.Sprintf("New Login:\nUid:%d\nEnterPassword:%s\nResult:%s\n", uid, tp, "true"))
		if err != nil {
			gologger.Logwrite(err)
		}
		conn.Write(data)
	} else {
		data, err := proto.Encode("false")
		gologger.BasicLogwrite(fmt.Sprintf("New Login:\nUid:%d\nEnterPassword:%s\nResult:%s\n", uid, tp, "false"))
		if err != nil {
			gologger.Logwrite(err)
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
	uid := mysql.Insert(name, p)
	gologger.BasicLogwrite(fmt.Sprintf("New Register:\nUid:%d\nName:%s\nPassword:%s\n", uid, name, p))
	data, err := proto.Encode(strconv.FormatInt(uid, 10))
	if err != nil {
		gologger.Logwrite(err)
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
		gologger.Logwrite(err)
	}
	if mysql.Checkp(uid, tp) {
		mysql.Pnp(uid, np)
		gologger.BasicLogwrite(fmt.Sprintf("New Changep:\nUid:%d\nOldp:%s\nNewp:%s\nResult:%s\n", uid, tp, np, "true"))
		data, err := proto.Encode("true")
		if err != nil {
			gologger.Logwrite(err)
		}
		conn.Write(data)
	} else {
		data, err := proto.Encode("false")
		gologger.BasicLogwrite(fmt.Sprintf("New Changep:\nUid:%d\nOldp:%s\nNewp:%s\nResult:%s\n", uid, tp, np, "false"))
		if err != nil {
			gologger.Logwrite(err)
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
			gologger.Logwrite(err)
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
	listen, err := net.Listen("tcp", "10.0.0.4:30000")
	if err != nil {
		gologger.Logwrite(err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			gologger.Logwrite(err)
			continue
		}
		fmt.Printf("Copyright ©2021 dreamxw.com All Right Reserved Powered by Azure")
		fmt.Printf("The Service is starting...")
		go process(conn)
	}
}
