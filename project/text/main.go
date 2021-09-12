package main

import (
	"fmt"
)

// var file1, _ = os.OpenFile("../log/log error"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
// var file2, _ = os.OpenFile("../log/log basic"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

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

func main() {
	fmt.Println(extraction("From the [72] times onwards, the probability of obtaining 5 stars gradually increases, each increment of [5.52%]", "[", "]", 2))
}
