package main

import (
	"fmt"
	"github.com/daoru"
	"io"
)

func login() {
	var uid int64
	var p string
	fmt.Println("请输入uid:")
	_, err := fmt.Scanln(&uid)
	if err != nil {
		fmt.Printf("输入存在问题,%v\n", err)
		return
	}
	fmt.Println("请输入密码:")
	_, _ = fmt.Scanln(&p)
	if daoru.Checkp(uid, p) {
		fmt.Println("登陆成功")
	} else {
		fmt.Println("密码不正确或不存在该账号")
	}
}

func check(p string) (b bool) {
	var (
		cAlphabets = 0
		lAlphabets = 0
		numbers    = 0
		characters = 0
	)
	for _, i := range p {
		// fmt.Println(i)
		if i >= 'a' && i <= 'z' {
			lAlphabets |= 1
		} else if i >= 'A' && i <= 'Z' {
			cAlphabets |= 1
		} else if i >= '0' && i <= '9' {
			numbers |= 1
		} else {
			characters |= 1
		}
	}
	// fmt.Printf("%d,%d,%d,%d", cAlphabets, lAlphabets, numbers, characters)
	if cAlphabets+lAlphabets+numbers+characters >= 3 {
		return false
	} else {
		return true
	}
}

func register() {
	var name string
	var p string
	fmt.Println("注册：\n输入姓名：要求：姓名字符不多于10个字符")
	_, _ = fmt.Scanln(&name, &p)
	if len(name) > 10 {
		fmt.Println("姓名不符合要求")
		return
	}
	fmt.Println("\n输入密码：要求：密码字符不得少于8个字符，不得多于24个字符，另外，您的密码必须包含以下任意三项：大写字母，小写字母，数字，字符")
	_, _ = fmt.Scanln(&p)
	if len(p) < 8 || len(p) > 24 || check(p) {
		fmt.Println("密码不符合要求")
		return
	} else {
		fmt.Printf("注册成功！信息如下:\nuid:%d\n姓名:%s\n密码：%s\n", daoru.Insert(name, p), name, p)
	}
}

func changep() {
	var uid int64
	var p string
	var np string
	fmt.Println("输入要修改密码的uid:")
	_, _ = fmt.Scanln(&uid)
	fmt.Println("输入旧密码:")
	_, _ = fmt.Scanln(&p)
	if daoru.Checkp(uid, p) {
		fmt.Println("输入新密码:要求：密码字符不得少于8个字符，不得多于24个字符，另外，您的密码必须包含以下任意三项：大写字母，小写字母，数字，字符")
		_, _ = fmt.Scanln(&np)
		if !check(np) {
			daoru.Pnp(uid, np)
			fmt.Println("修改成功")
		} else {
			fmt.Println("新密码不规范")
		}
	} else {
		fmt.Println("原密码不正确或不存在该账号")
	}
}

func main() {
	var input int
	for {
		fmt.Printf("欢迎！请输入要操作的事项：\n1.登录\n2.注册\n3.修改密码\n其他.退出\n")
		fmt.Printf("输入你需要进行的操作：")
		_, err := fmt.Scanln(&input)
		if err == io.EOF {
			break
		}
		if input == 1 {
			login()
		} else if input == 2 {
			register()
		} else if input == 3 {
			changep()
		} else {
			break
		}
	}
}
