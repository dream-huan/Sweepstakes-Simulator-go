package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var file1, _ = os.OpenFile("../log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
var file2, _ = os.OpenFile("../log/log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

func main() {
	var a string
	n, _ := fmt.Scanln(&a)
	fmt.Println(n)
	if n == 0 {
		log.SetOutput(file1)
		log.SetPrefix("[Error]")
		log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
		log.Printf("123")
	}
}
