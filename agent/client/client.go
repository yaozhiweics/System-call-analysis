package main

import (
	"fmt"
	"net"
	"syscall/agent/proto"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8086")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			msg := `Hello, Hello. How are you?`
			data, err := proto.Encode(msg)
			if err != nil {
				fmt.Println("encode msg failed, err:", err)
				return
			}
			conn.Write(data)
			time.Sleep(2 * time.Second)
		}

	}()
	// for i := 0; i < 20; i++ {
	// 	msg := `Hello, Hello. How are you?`
	// 	data, err := proto.Encode(msg)
	// 	if err != nil {
	// 		fmt.Println("encode msg failed, err:", err)
	// 		return
	// 	}
	// 	conn.Write(data)
	// }
	time.Sleep(2000000 * time.Second)
}
