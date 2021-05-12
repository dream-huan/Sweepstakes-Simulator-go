package gologger

import (
	"log"
	"os"
	"time"
)

func Logwrite(err error) {
	file, err := os.OpenFile("../log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Logwrite(err)
	}
	log.SetOutput(file)
	log.SetPrefix("[Error]")
	log.SetFlags(log.Ldate)
	log.Printf("%#v", *err)
}

func BasicLogwrite(msg string) {
	file, err := os.OpenFile("../log/log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Logwrite(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate)
	log.Printf("%#v", msg)
}
