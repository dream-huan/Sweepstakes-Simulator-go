package gologger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

func Logwrite(str string) {
	file, err := os.OpenFile("../log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
	}
	log.SetOutput(file)
	log.SetPrefix("[Error]")
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	log.Printf("%#v", str)
}

func BasicLogwrite(msg string) {
	file, err := os.OpenFile("../log/log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		Logwrite(fmt.Sprintf("%v", err) + " " + file + " " + strconv.Itoa(line) + strconv.FormatBool(ok))
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate)
	log.Printf("%#v", msg)
}
