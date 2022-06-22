package controls

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"syscall/datas"
	"syscall/utils"

	log "github.com/sirupsen/logrus"
)

func InitServer() {
	// 1. 绑定ip和端口，设置监听
	listener, err := net.Listen("tcp", "127.0.0.1:8086")
	if err != nil {
		log.Error("Failed to Listen", err)
	}
	// 延迟关闭，释放资源
	defer listener.Close()

	// 2. 循环等待新连接
	for {
		// 从连接列表获取新连接
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Failed to Accept", err)
		}
		// 3. 与新连接通信(为了不同步阻塞，这里开启异步协程进行函数调用)
		go handle_conn(conn)
	}
}

func handle_conn(conn net.Conn) {
	defer conn.Close()
	log.Info("New connection ", conn.RemoteAddr())
	// 通信
	// buf := make([]byte, 409600)
	reader := bufio.NewReader(conn)
	for {
		msg, err := utils.Decode(reader)
		if err == io.EOF {
			return
		}
		// // 从网络中读
		// readBytesCount, err := conn.Read(buf)
		if err != nil {
			log.Error("Failed to read", err)
			break
		}
		// // 提示：buf[:n]的效果为：读取buf[总长度-n]至buf[n]处的字节
		//log.Info("get data from ", conn.RemoteAddr(), ":", msg)
		saveData(msg)
	}
}

func saveData(data string) {
	responseData := &datas.ResponseData{}
	err := json.Unmarshal([]byte(data), responseData)
	if err != nil {
		log.Error("parse json err: ", err)
	}
	//过滤消息，不处理初始化消息和心跳检测消息
	if responseData.ModuleType == datas.MODULE_SYSCALL {
		sysMsg := &datas.SyscallData{}
		err := json.Unmarshal([]byte(responseData.Datas), sysMsg)
		if err != nil {
			log.Error("parse json err: ", err)
		}
		//插入数据库

		res := Insert(sysMsg)
		if res < 0 {
			log.Error("insert err:", sysMsg)
		} else {
			log.Info("inser success:", sysMsg)
		}

	}
}
