package gologger

import (
	"log"
	"os"
	"time"
)

// var file1 *os.File
var file2 *os.File

func init() {
	// file1, _ = os.OpenFile("../log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	file2, _ = os.OpenFile("../log/log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

/* func Logwrite(err error) {
	log.SetOutput(file1)
	log.SetPrefix("[Error]")
	log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
	log.Printf("%#v", err)
} */

func BasicLogwrite(msg string) {
	log.SetOutput(file2)
	log.SetPrefix("[Basic]")
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	log.Printf("%#v", msg)
}
