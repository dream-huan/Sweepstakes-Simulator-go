package main

import (
	"os"
	"strconv"
)

// var file1, _ = os.OpenFile("../log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
// var file2, _ = os.OpenFile("../log/log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

func main() {
	var uid int64
	uid = 3
	os.RemoveAll("/workspaces/go-chat/project/" + strconv.FormatInt(uid, 10))
}
