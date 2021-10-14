package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// 客户端
func main1() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn.Close() // 关闭连接
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString(' ') // 读取用户输入
		inputInfo := strings.Trim(input, " \n")
		if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
			fmt.Println("连接已关闭")
			return
		}
		_, err = conn.Write([]byte(inputInfo)) // 发送数据
		if err != nil {
			return
		}
		//buf := [512]byte{}
		//n, err := conn.Read(buf[:])
		//if err != nil {
		//	fmt.Println("recv failed, err:", err)
		//	return
		//}
		//msgFromServer := string(buf[:n])
		//fmt.Println("收到server端发来的数据：", msgFromServer)
	}
}
