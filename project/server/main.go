package main

import (
	"bufio"
	"dream/gologger"
	"dream/mysql"
	"dream/proto"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var file1, _ = os.OpenFile("log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
var file2, _ = os.OpenFile("log/log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

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
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
	if mysql.Checkp(uid, tp) {
		data, err := proto.Encode("true")
		gologger.BasicLogwrite(fmt.Sprintf("New Login: Uid:%d EnterPassword:%s Result:%s ", uid, tp, "true"))
		if err != nil {
			log.SetOutput(file1)
			log.SetPrefix("[Error]")
			log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
			log.Printf("%v", err)
		}
		conn.Write(data)
	} else {
		data, err := proto.Encode("false")
		gologger.BasicLogwrite(fmt.Sprintf("New Login: Uid:%d EnterPassword:%s Result:%s ", uid, tp, "false"))
		if err != nil {
			log.SetOutput(file1)
			log.SetPrefix("[Error]")
			log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
			log.Printf("%v", err)
		}
		conn.Write(data)
	}
}

func checktoggle(msg *string, conn net.Conn) {
	var uid int64
	var pool int
	var tu string
	k := 0
	for _, r := range *msg {
		if r == ' ' {
			k += 1
		}
		if k == 1 && r != ' ' {
			tu += string(r)
		}
	}
	k = 1
	uid, _ = strconv.ParseInt(tu, 10, 64)
	pool = mysql.Checkpool(uid)
	file, err := os.Open("data/cardpool.txt")
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var content string
	for {
		line, err := reader.ReadString('\n') //注意是字符
		if k == pool {
			content = line
			break
		}
		k += 1
		if err == io.EOF {
			break
		}
	}
	data, _ := proto.Encode(content)
	conn.Write(data)
}

func toggle(msg *string, conn net.Conn) {
	var uid int64
	var pool int
	var tu string
	var tp string
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
	pool, _ = strconv.Atoi(tp)
	mysql.Toggle(uid, pool)
}

func take(takes=1 int){
	
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
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
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
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
	if mysql.Checkp(uid, tp) {
		mysql.Pnp(uid, np)
		gologger.BasicLogwrite(fmt.Sprintf("New Changep: Uid:%d Oldp:%s Newp:%s Result:%s ", uid, tp, np, "true"))
		data, err := proto.Encode("true")
		if err != nil {
			log.SetOutput(file1)
			log.SetPrefix("[Error]")
			log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
			log.Printf("%v", err)
		}
		conn.Write(data)
	} else {
		data, err := proto.Encode("false")
		gologger.BasicLogwrite(fmt.Sprintf("New Changep: Uid:%d Oldp:%s Newp:%s Result:%s ", uid, tp, np, "false"))
		if err != nil {
			log.SetOutput(file1)
			log.SetPrefix("[Error]")
			log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
			log.Printf("%v", err)
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
			log.SetOutput(file1)
			log.SetPrefix("[Error]")
			log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
			log.Printf("%v", err)
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
		case "checktoggle":
			checktoggle(&msg, conn)
		case "toggle":
			toggle(&msg, conn)
		default:
			break
		}
		// fmt.Println("收到client发来的数据：", msg)
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("%v", err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.SetOutput(file1)
			log.SetPrefix("[Error]")
			log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
			log.Printf("%v", err)
			continue
		}
		fmt.Printf("Copyright ©2021 dreamxw.com All Right Reserved Powered by Azure")
		fmt.Printf("The Service is starting...")
		go process(conn)
	}
}
