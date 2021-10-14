package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main(){
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn.Close() // 关闭连接
	inputReader := bufio.NewReader(os.Stdin)

	go func() {
		// 读取群聊数据
		inputFromServer := bufio.NewScanner(conn)
		for inputFromServer.Scan() {
			fmt.Println(inputFromServer.Text())
		}
	}()

	for {
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		inputInfo := strings.Trim(input, " \n")
		if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
			fmt.Println("连接已关闭")
			return
		}
		_, err = conn.Write([]byte(inputInfo)) // 发送数据
		if err != nil {
			return
		}
	}
}
