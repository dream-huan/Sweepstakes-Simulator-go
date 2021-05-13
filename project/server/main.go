package main

import (
	"bufio"
	"dream/gologger"
	"dream/mysql"
	"dream/proto"
	"fmt"
	"io"
	"net"
	"runtime"
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
		_, file, line, ok := runtime.Caller(1)
		gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
	}
	if mysql.Checkp(uid, tp) {
		data, err := proto.Encode("true")
		gologger.BasicLogwrite(fmt.Sprintf("New Login: Uid:%d EnterPassword:%s Result:%s ", uid, tp, "true"))
		if err != nil {
			_, file, line, ok := runtime.Caller(1)
			gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
		}
		conn.Write(data)
	} else {
		data, err := proto.Encode("false")
		gologger.BasicLogwrite(fmt.Sprintf("New Login: Uid:%d EnterPassword:%s Result:%s ", uid, tp, "false"))
		if err != nil {
			_, file, line, ok := runtime.Caller(1)
			gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
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
	gologger.BasicLogwrite(fmt.Sprintf("New Register: Uid:%d Name:%s Password:%s ", uid, name, p))
	data, err := proto.Encode(strconv.FormatInt(uid, 10))
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
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
		_, file, line, ok := runtime.Caller(1)
		gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
	}
	if mysql.Checkp(uid, tp) {
		mysql.Pnp(uid, np)
		gologger.BasicLogwrite(fmt.Sprintf("New Changep: Uid:%d Oldp:%s Newp:%s Result:%s ", uid, tp, np, "true"))
		data, err := proto.Encode("true")
		if err != nil {
			_, file, line, ok := runtime.Caller(1)
			gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
		}
		conn.Write(data)
	} else {
		data, err := proto.Encode("false")
		gologger.BasicLogwrite(fmt.Sprintf("New Changep: Uid:%d Oldp:%s Newp:%s Result:%s ", uid, tp, np, "false"))
		if err != nil {
			_, file, line, ok := runtime.Caller(1)
			gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
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
			_, file, line, ok := runtime.Caller(1)
			gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
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
		_, file, line, ok := runtime.Caller(1)
		gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			_, file, line, ok := runtime.Caller(1)
			gologger.Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
			continue
		}
		fmt.Printf("Copyright ©2021 dreamxw.com All Right Reserved Powered by Azure")
		fmt.Printf("The Service is starting...")
		go process(conn)
	}
}
