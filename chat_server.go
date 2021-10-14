package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

//
type client chan<- string // 客户端，接收消息用的通道
var (
	entering = make(chan client)
	leaving  = make(chan client)
	msgs     = make(chan string) // 用于群发消息的通道

)

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case message := <-msgs:
			for cli := range clients {
				cli <- message
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // 服务端给每个客户端发送消息的通道
	go clientWriter(conn, ch)
	ch <- "type your name : "

	nameReader := bufio.NewReader(conn)
	var nameBuf [128]byte
	n, err := nameReader.Read(nameBuf[:]) // 读取数据
	if err != nil {
		return
	}
	name := string(nameBuf[:n])

	ch <- "You are :" + name
	entering <- ch

	msgs <- name + " has arrived"
	fmt.Println(name + " has arrived")

	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			break
		}
		recvStr := string(buf[:n])
		msgs <- name + " ：" + recvStr
		fmt.Println(name, " ：", recvStr)
	} // 发送数据

	leaving <- ch
	msgs <- name + "  has left"
	fmt.Println(name + "  has left")
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000") // 监听端口，接受来自客户端的连接请求
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
