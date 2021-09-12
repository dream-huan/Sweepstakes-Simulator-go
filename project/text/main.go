package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

// var file1, _ = os.OpenFile("../log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
// var file2, _ = os.OpenFile("../log/log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

func main() {
	uid := 3
	// file3, _ := os.OpenFile(strconv.Itoa(uid)+".html", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	content := []byte(`<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<title>haha</title>
</head>
<body>`)
	err := ioutil.WriteFile(strconv.Itoa(uid)+".html", content, 0644)
	file, err := os.Open(strconv.Itoa(uid) + ".txt")
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var f *os.File
	f, _ = os.OpenFile(strconv.Itoa(uid)+".html", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer f.Close()
	for {
		line, err := reader.ReadString('\n') //注意是字符
		_, _ = io.WriteString(f, line+"<br>")
		if err == io.EOF {
			break
		}
	}
	io.WriteString(f, "\n</body>\n</html>")
}
