package main

import (
	"bufio"
	"dream/gologger"
	"dream/mysql"
	"dream/proto"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// var file1, _ = os.OpenFile("log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

// var file2, _ = os.OpenFile("log/log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
var fivestar1 [1024]string //储存有哪些五星
var fourstar1 [1024]string //同理上面
var threestar1 [1024]string
var fivestarsum1 = 0 //五星的数量
var fourstarsum1 = 0
var threestarsum1 = 0
var profive float64 //五星的概率
var profour float64
var prothree float64
var initprofive float64
var guafour int        //保底四星次数
var guafive int        //五星概率提升次数
var guaprofive float64 //五星概率提升
var workspace string
var ipconfig string

func extraction(s string, a string, b string, times int) (str string) {
	str = ""
	k1 := 0
	for _, r := range s {
		if string(r) == a {
			k1 += 1
			continue
		}
		if string(r) == b && k1 == times {
			break
		}
		if k1 == times {
			str += string(r)
		}
	}
	return str
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
	if mysql.Checkp(uid, tp) {
		data, _ := proto.Encode("true")
		gologger.BasicLogwrite(fmt.Sprintf("New Login: Uid:%d EnterPassword:%s Result:%s ", uid, tp, "true"))
		conn.Write(data)
	} else {
		data, _ := proto.Encode("false")
		gologger.BasicLogwrite(fmt.Sprintf("New Login: Uid:%d EnterPassword:%s Result:%s ", uid, tp, "false"))
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
	file, _ := os.Open("data/cardpool.txt")
	defer file.Close()
	reader := bufio.NewReader(file)
	gologger.BasicLogwrite(fmt.Sprintf("CheckToggle: Uid:%d ", uid))
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
	gologger.BasicLogwrite(fmt.Sprintf("Toggle: Uid:%d Pool:%d ", uid, pool))
	mysql.Toggle(uid, pool)
}

func recharge(msg *string, conn net.Conn) {
	var uid int64
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
	if mysql.Recharge(tp, uid) == true {
		data, _ := proto.Encode("充值状态:true")
		gologger.BasicLogwrite(fmt.Sprintf("Recharge: Uid:%d Keycode:%s Result:%s", uid, tp, "true"))
		conn.Write(data)
	} else {
		data, _ := proto.Encode("充值状态:false")
		gologger.BasicLogwrite(fmt.Sprintf("Recharge: Uid:%d Keycode:%s Result:%s", uid, tp, "false"))
		conn.Write(data)
	}
}

func take(msg *string, conn net.Conn) {
	var uid int64
	var times int
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
	uid, _ = strconv.ParseInt(tp, 10, 64)
	times, _ = strconv.Atoi(tu)
	fivetimes, fourtimes := mysql.Checkstatistics(uid)
	message := ""
	stone := mysql.Getstone(uid)
	gologger.BasicLogwrite(fmt.Sprintf("Take: Uid:%d times:%d Stone:%d ", uid, times, stone))
	if stone < times*160 {
		message := "原石不足，现有原石:" + strconv.Itoa(stone)
		data, _ := proto.Encode(message)
		conn.Write(data)
		return
	}
	stone -= 160 * times
	for {
		if times == 0 {
			break
		}
		rand.Seed(time.Now().UnixNano())
		probability := rand.Intn(999) + 1
		now := time.Now()
		if fourtimes == 10 {
			fourtimes = 1
			probability := rand.Intn(fourstarsum1 - 1)
			// fmt.Printf("%v(四星)", fourstar1[probability])
			mysql.Insertprop(uid, fourstar1[probability]+"(四星)", now.Format("2006/01/02 15:04:05"))
			message += (fourstar1[probability] + "(四星)")
			times -= 1
			fivetimes += 1
			if times > 1 {
				message += ","
			}
			continue
		}
		if probability >= 1 && probability <= int(profive*1000) {
			probability := rand.Intn(fivestarsum1 - 1)
			// fmt.Printf("%v(五星)", fivestar1[probability])
			mysql.Insertprop(uid, fivestar1[probability]+"(五星)", now.Format("2006/01/02 15:04:05"))
			message += (fivestar1[probability] + "(五星)")
			profive = initprofive
			fivetimes = 1
		} else if probability >= int(profive*1000)+1 && probability <= int(profour*1000) {
			probability := rand.Intn(fourstarsum1 - 1)
			// fmt.Printf("%v(四星)", fourstar1[probability])
			mysql.Insertprop(uid, fourstar1[probability]+"(四星)", now.Format("2006/01/02 15:04:05"))
			message += (fourstar1[probability] + "(四星)")
		} else {
			probability := rand.Intn(threestarsum1 - 1)
			// fmt.Printf("%v(三星)", threestar1[probability])
			mysql.Insertprop(uid, threestar1[probability]+"(三星)", now.Format("2006/01/02 15:04:05"))
			message += (threestar1[probability] + "(三星)")
		}
		if times > 1 {
			message += ","
		}
		times -= 1
		fourtimes += 1
		fivetimes += 1
		if fivetimes >= guafive {
			profive += guaprofive
		}
	}
	mysql.Changestone(stone, uid)
	mysql.Timeschange(uid, fivetimes, fourtimes)
	data, _ := proto.Encode(message)
	conn.Write(data)
}

func checkresult(msg *string, conn net.Conn) {
	var uid int64
	var tp string
	k := 0
	for _, r := range *msg {
		if r == ' ' {
			k += 1
		}
		if k == 1 && r != ' ' {
			tp += string(r)
		}
	}
	uid, _ = strconv.ParseInt(tp, 10, 64)
	os.RemoveAll(workspace + strconv.FormatInt(uid, 10))
	err := os.Mkdir(strconv.FormatInt(uid, 10), os.ModePerm)
	os.Chmod(strconv.FormatInt(uid, 10), 0777)
	mysql.Checkresult(uid)
	gologger.BasicLogwrite(fmt.Sprintf("CheckResult: Uid:%d  ", uid))
	// file3, _ := os.OpenFile(strconv.Itoa(uid)+".html", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	content := []byte(`<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<title>记录查询</title>
</head>
<body>
<table border="1" style="border-collapse: collapse;">
<caption>常驻祈愿</caption>
<tr>
<td style="text-align:center">UID</td>
<td style="text-align:center">道具</td>
<td style="text-align:center">祈愿时间</td>
</tr>`)
	err = ioutil.WriteFile(strconv.FormatInt(uid, 10)+"/"+strconv.FormatInt(uid, 10)+".html", content, 0644)
	file, err := os.Open(strconv.FormatInt(uid, 10) + "/" + strconv.FormatInt(uid, 10) + ".txt")
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var f *os.File
	f, _ = os.OpenFile(strconv.FormatInt(uid, 10)+"/"+strconv.FormatInt(uid, 10)+".html", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer f.Close()
	for {
		line, err := reader.ReadString('\n') //注意是字符
		_, _ = io.WriteString(f, "<tr>\n")
		k := 0
		temp := ""
		for _, r := range line {
			if r == '	' {
				k += 1
				if strings.Contains(temp, "三星") {
					_, _ = io.WriteString(f, "<td style=\"color:#46A3FF\">"+temp+"</td>\n")
				} else if strings.Contains(temp, "四星") {
					_, _ = io.WriteString(f, "<td style=\"color:#9F35FF\">"+temp+"</td>\n")
				} else if strings.Contains(temp, "五星") {
					_, _ = io.WriteString(f, "<td style=\"color:#FF0000\">"+temp+"</td>\n")
				} else {
					_, _ = io.WriteString(f, "<td>"+temp+"</td>")
				}
				temp = ""
				continue
			}
			temp += string(r)
		}
		_, _ = io.WriteString(f, "<td>"+temp+"</td>\n")
		_, _ = io.WriteString(f, "</tr>\n")
		if err == io.EOF {
			break
		}
	}
	io.WriteString(f, "\n</table>\n</body>\n</html>")
}

func checkstatistics(msg *string, conn net.Conn) {
	var uid int64
	var tp string
	k := 0
	for _, r := range *msg {
		if r == ' ' {
			k += 1
		}
		if k == 1 && r != ' ' {
			tp += string(r)
		}
	}
	uid, _ = strconv.ParseInt(tp, 10, 64)
	t1, t2 := mysql.Checkstatistics(uid)
	// fmt.Printf("五星已有%v次没出，四星已有%v次没出", t1, t2)
	message := "五星已有" + strconv.Itoa(t1) + "次没出，四星已有" + strconv.Itoa(t2) + "次没出"
	gologger.BasicLogwrite(fmt.Sprintf("CheckStatistics: Uid:%d  ", uid))
	data, _ := proto.Encode(message)
	conn.Write(data)
}

func getstone(msg *string, conn net.Conn) {
	var uid int64
	var tp string
	k := 0
	for _, r := range *msg {
		if r == ' ' {
			k += 1
		}
		if k == 1 && r != ' ' {
			tp += string(r)
		}
	}
	uid, _ = strconv.ParseInt(tp, 10, 64)
	message := mysql.Getstone(uid)
	data, _ := proto.Encode(strconv.Itoa(message))
	conn.Write(data)
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
	if mysql.Checkp(uid, tp) {
		mysql.Pnp(uid, np)
		gologger.BasicLogwrite(fmt.Sprintf("New Changep: Uid:%d Oldp:%s Newp:%s Result:%s ", uid, tp, np, "true"))
		data, _ := proto.Encode("true")
		conn.Write(data)
	} else {
		data, _ := proto.Encode("false")
		gologger.BasicLogwrite(fmt.Sprintf("New Changep: Uid:%d Oldp:%s Newp:%s Result:%s ", uid, tp, np, "false"))
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
		case "take":
			take(&msg, conn)
		case "checkresult":
			checkresult(&msg, conn)
		case "checkstatistics":
			checkstatistics(&msg, conn)
		case "getstone":
			getstone(&msg, conn)
		case "recharge":
			recharge(&msg, conn)
		}
		// fmt.Println("收到client发来的数据：", msg)
	}
}

func init() {
	//配置文件读取
	file, err := os.Open("data/3.pool")
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	k := 6
	for {
		line, err := reader.ReadString('\n') //注意是字符
		if extraction(line, "[", "]", 1) == "5star" {
			k = 5
		} else if extraction(line, "[", "]", 1) == "4star" {
			k = 4
		} else if extraction(line, "[", "]", 1) == "3star" {
			k = 3
		} else if extraction(line, "[", "]", 1) != "" {
			if k == 5 {
				fivestar1[fivestarsum1] = extraction(line, "[", "]", 1)
				fivestarsum1 += 1
			} else if k == 4 {
				fourstar1[fourstarsum1] = extraction(line, "[", "]", 1)
				fourstarsum1 += 1
			} else if k == 3 {
				threestar1[threestarsum1] = extraction(line, "[", "]", 1)
				threestarsum1 += 1
			}
		}
		if err == io.EOF {
			break
		}
	}
	for i := 0; i < fivestarsum1; i += 1 {
		fmt.Println(fivestar1[i])
	}
	file, err = os.Open("data/setting.in")
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	reader = bufio.NewReader(file)
	k = 1
	//概率分配
	for {
		line, err := reader.ReadString('\n') //注意是字符
		if k == 4 {
			profive, err = strconv.ParseFloat(extraction(line, "[", "]", 1), 64)
			initprofive = profive
			if err != nil {
				fmt.Printf("error:%v", err)
			}
		}
		if k == 5 {
			profour, err = strconv.ParseFloat(extraction(line, "[", "]", 1), 64)
			if err != nil {
				fmt.Printf("error:%v", err)
			}
		}
		if k == 6 {
			prothree, err = strconv.ParseFloat(extraction(line, "[", "]", 1), 64)
			if err != nil {
				fmt.Printf("error:%v", err)
			}
		}
		if k == 9 {
			guafive, err = strconv.Atoi(extraction(line, "[", "]", 1))
			if err != nil {
				fmt.Printf("error:%v", err)
			}
			guaprofive, err = strconv.ParseFloat(extraction(line, "[", "]", 2), 64)
			if err != nil {
				fmt.Printf("error:%v", err)
			}
		}
		if k == 10 {
			guafour, err = strconv.Atoi(extraction(line, "[", "]", 1))
			if err != nil {
				fmt.Printf("error:%v", err)
			}
		}
		if k == 13 {
			workspace = extraction(line, "[", "]", 1)
		}
		if k == 16 {
			ipconfig = extraction(line, "[", "]", 1)
		}
		if err == io.EOF {
			break
		}
		k += 1
	}
	fmt.Printf("%v\n%v\n%v\n%v\n%v\n%v\n", profive, profour, prothree, guafive, guaprofive, guafour)
}

func main() {
	fmt.Printf("Copyright ©2021 dreamxw.com All Right Reserved Powered by Azure")
	listen, _ := net.Listen("tcp", ipconfig)
	defer listen.Close()
	for {
		conn, _ := listen.Accept()
		go process(conn)
	}
}
