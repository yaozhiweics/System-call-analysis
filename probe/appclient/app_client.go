package appclient

import (
	"encoding/json"
	"math/rand"
	"net"
	"sync"
	"syscall/datas"
	"syscall/utils"

	"time"

	log "github.com/sirupsen/logrus"
)

var safeMode = false                                  //重连失败后是否进入测试模式
var timeHeart = 6                                     //心跳时间间隔，从配置文件读取，单位：秒
var timeout = 10                                      //读写超时时间，从配置文件读取，单位：秒
var tryNum = 2                                        //重连尝试次数, 从配置文件读取
var tryTime = 3                                       //重连尝试时间间隔, 从配置文件读取，单位：秒
var agentIp = "127.0.0.1"                             //从配置文件读取agentip，目前为手动修改
var dataChannels = make(chan datas.ResponseData, 100) //这个结构有缓存，也阻塞了。
var wg sync.WaitGroup                                 //用于保证不同协程的执行顺序
var aliveConnP = make(chan net.Conn, 1)               //若当前连接发生异常，则新建连接并放入此处。解析模块携程从此处读，上传模块从此处写入。
var aliveConnU = make(chan net.Conn, 1)               //若当前连接发生异常，则新建连接并放入此处。上传模块携程从此处读，解析模块从此处写入。
var aliveConnH = make(chan net.Conn, 1)               //若当前连接发生异常，则新建连接并放入此处。心跳模块携程从此处读，心跳模块从此处写入。

var allIfAddress = make([]string, 0)

//StartClient 通信模块母模块
func StartClient(remoteAddress []string, localAddress string, allAddress []string, dc chan datas.ResponseData, conf *datas.Conf) {
	dataChannels = dc
	agentIp = localAddress
	allIfAddress = allAddress
	timeHeart = conf.TimeHeart //心跳时间间隔，从配置文件读取，单位：秒
	timeout = conf.Timeout     //读写超时时间，从配置文件读取，单位：秒
	tryNum = conf.TryNum       //重连尝试次数, 从配置文件读取
	tryTime = conf.TryTime     //重连尝试时间间隔, 从配置文件读取，单位：秒
	safeMode = conf.SafeMode
	//多个服务端地址remoteaddress
	log.Infoln("remote address is ", remoteAddress)

	var tcpAddrs []*net.TCPAddr
	for _, ipAddr := range remoteAddress {
		tcpAddr, err := net.ResolveTCPAddr("tcp", ipAddr)
		if err != nil {
			log.Error(err)
		}
		tcpAddrs = append(tcpAddrs, tcpAddr)
	}

	//建立一个初始连接
	conn, _ := createMultiConn(tcpAddrs)
	//发送初始化消息
	initAgent(conn, conf)
	aliveConnH <- conn
	aliveConnU <- conn
	aliveConnP <- conn
	//启动一个携程，检测连接状态，操作变量isAlive
	heartAlive(conn, conf, tcpAddrs)
	//解析datachan，组装数据上传给server
	uploadData(conn)
	//从server读取命令，解析命令，将命令按照类型放入三种cmdchan中
	wg.Wait()
}

//调用createConn同时尝试像多个本地服务建立连接
//任何一个连接建立成功，则通知其它协程不再继续尝试
//返回建立成功的连接conn和真正建立连接的net.TCPAddr
func createMultiConn(remoteAddress []*net.TCPAddr) (net.Conn, *net.TCPAddr) {
	var conn net.Conn
	var realtcpAddr *net.TCPAddr
	var err error

	allErr := true

	for {
		//打乱IP池
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(remoteAddress), func(i, j int) {
			remoteAddress[i], remoteAddress[j] = remoteAddress[j], remoteAddress[i]
		})
		for _, tcpAddr := range remoteAddress {
			//尝试IP建立连接
			conn, err = createConn(tcpAddr)
			if err != nil {
				log.Error("client: dial: ", err)
				log.Info("change IP")
			} else {
				log.Info(tcpAddr, "连接建立成功，返回该连接")
				realtcpAddr = tcpAddr
				allErr = false
				break
			}
		}
		if allErr {
			// 全部都失败了，进入goTest
			goTest()
		} else {
			break
		}
	}

	return conn, realtcpAddr
}

//所有ip都重连失败调用这个函数
func goTest() {
	//读取配置文件，确定是否需要安全支持
	if safeMode {
		log.Info("建立连接失败")
	}
	//继续重连
}

//建立tcp连接
//若建立连接失败，尝试一定次数重连，重连次数由配置文件读取，超过次数后，则记录日志，报错，关闭客户端
//即，本函数在最大限度内尝试建立连接
//return: Conn, error
//连接建立失败，则err不为空
func createConn(remoteAddress *net.TCPAddr) (net.Conn, error) {
	var conn net.Conn
	var err error
	//尝试连接次数，相隔时间间隔
	//一定时间内一直重连对poc有压力，且粗暴不优雅
	for i := 0; i < tryNum; i++ {
		//增加证书验证等
		// caCert, _ := ioutil.ReadFile("/etc/kraken_agent/ca.crt")
		// caCertPool := x509.NewCertPool()
		// caCertPool.AppendCertsFromPEM(caCert)
		// //双向认证需要提供客户端证书，参数：二者是未加密公私钥对
		// // clientCert, err := tls.LoadX509KeyPair("./client.crt", "./client.key.text")
		// // if err != nil {
		// // 	log.Error(err)
		// // }
		// config := &tls.Config{
		// 	// Certificates:       []tls.Certificate{clientCert},
		// 	RootCAs:            caCertPool,
		// 	InsecureSkipVerify: false,
		// }
		// conn, err = tls.Dial("tcp", remoteAddress.String(), config)
		conn, err = net.DialTCP("tcp", nil, remoteAddress)
		if err != nil {
			log.Error("client: dial: ", err)
			time.Sleep(time.Second * time.Duration(tryTime))
		}
		// conn, err := net.DialTCP("tcp", nil, remoteAddress)
		if err == nil {
			// conn.SetKeepAlive(true) //设置保持长连接，可以降低TCP连接时的握手开销
			log.Info("连接建立成功，返回该连接")
			break
		}
	}
	return conn, err
}

//通过心跳判断连接是否可用
//如果不可用，就调用keepAlive更新并同步连接
func heartAlive(conn net.Conn, conf *datas.Conf, remoteAddressList []*net.TCPAddr) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		var try bool
		tick := time.Tick(time.Second * time.Duration(timeHeart))
		for now := range tick {
			select {
			case conn = <-aliveConnH:
				log.Info("H:连接更新")
				try = false
			default:
				//------------------------------------发送数据------------------------------------------------------------------------------------
				//每过一个时间间隔就发送心跳包，探测是否连接可用
				err := conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second)) //设置读写超时时间
				if err != nil {
					log.Error("心跳：写入数据前，设置deadline出错", err)
				}
				log.Debug("当前时间", now)
				//写入心跳包
				heartbeatdata := datas.ResponseData{}
				heartbeatdata.RequestId = "0"
				heartbeatdata.Datas = "HEARTBEAT"
				heartbeatdata.AgentIp = agentIp
				heartbeatdata.ModuleType = datas.MODULE_CONTROL
				jsonbytes, err := json.Marshal(heartbeatdata)
				if err != nil {
					log.Error("发送心跳包时，转json出错")
					continue
				}
				//jsonbytes = append(jsonbytes, byte('\n'))
				// jsonbytes = jsonbytes.join([]byte("\n"))
				data, err := utils.Encode(string(jsonbytes))
				if err != nil {
					log.Error("encode msg failed, err:", err)
				}
				n, err := conn.Write(data) //向连接写入心跳包
				if err != nil {
					log.Error("心跳：写入数据时出错： ", err)
					log.Info("尝试重连")
					//有个问题：server断开后，不断建立协程尝试重连，当server正常时，会累计许多新连接。
					//方法：设定一个局部变量，保证只有一个协程在尝试建立连接而不是随时间增加而增加协程数量。
					//写入数据出错，认为连接有问题，重新建立连接，并通知广播连接模块
					if !try { //使用try保证只会启动一个协程尝试连接。
						keepAlive(remoteAddressList, conf)
						try = true
					}
					log.Info("H:记录一下，心跳出问题，需要重连")
					// time.Sleep(time.Second * time.Duration(tryTime))
					continue
				}
				log.Debug("send ", (n), " heartbeat bytes to ", conn.RemoteAddr(), "\n")
			}
		}
	}()
}

//作为拓展区，存储数据提供给Datas
type ExtraData struct {
	//可以扩展字段
	// Conf            datas.Conf              `json:"conf"`
	AllIfAddress []string `json:"allIfAddress"`
}

//发送初始化包，如果发送失败，会被动等待。
func initAgent(conn net.Conn, conf *datas.Conf) {
	//发送初始化包
	initpackage := datas.ResponseData{}
	initpackage.AgentIp = agentIp
	initpackage.ModuleType = datas.MODULE_INIT
	initpackage.RequestId = "INIT"
	extraData := ExtraData{
		AllIfAddress: allIfAddress,
	}
	extraDatasBytes, err := json.Marshal(extraData)
	if err != nil {
		log.Error(err)
	}
	initpackage.Datas = string(extraDatasBytes)
	jsonbytes, err := json.Marshal(initpackage)
	if err != nil {
		log.Error("发送初始化包时，转json出错")
		//疑问，出错后该干啥，不过一般不会错
		//待处理
	}
	// println(string(jsonbytes))
	// jsonbytes = append(jsonbytes, byte('\n'))
	data, err := utils.Encode(string(jsonbytes))
	if err != nil {
		log.Error("encode msg failed, err:", err)
	}
	n, err := conn.Write(data) //向连接写入初始化包
	if err != nil {
		log.Error("初始化：写入数据时出错： ", err)
		log.Info("尝试重连")
	}
	log.Info("send ", (n), " init bytes to ", conn.RemoteAddr(), "\n")
}

//keepAlive负责创建新的连接，并创建活跃的连接存放于aliveConn中。
//同步三个模块之间的连接同步。
//类似广播功能，将最新的连接广播出去
//如果建立连接出点error，可能是由于此模块还没有将conn建立好，其他地方就调用了
func keepAlive(remoteAddressList []*net.TCPAddr, conf *datas.Conf) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		// getconn:
		conn, _ := createMultiConn(remoteAddressList)
		//发送初始化包
		initAgent(conn, conf)

		//给每个模块的专属管道通知一下新的conn
		aliveConnH <- conn
		aliveConnP <- conn
		aliveConnU <- conn
	}()
}

//监听cmd chan，使用对应的zcmd获取对应的数据chan，返回给于服务器的连接
func uploadData(conn net.Conn) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			//连接更新，取新连接
			case conn = <-aliveConnU:
				log.Info("U:连接更新")
			case d := <-dataChannels:
				log.Info("data is ", d)
				jsonbytes, err := json.Marshal(d)
				if err != nil {
					log.Error(err)
				}
				//把数据写入连接
				data, err := utils.Encode(string(jsonbytes))
				if err != nil {
					log.Error("encode msg failed, err:", err)
				}
				//jsonbytes = append(jsonbytes, byte('\n'))
				_, err = conn.Write(data)
				if err != nil {
					log.Error("U:连接写入数据出错 ", err)
					log.Info("尝试重连")
					// time.Sleep(time.Second * time.Duration(tryTime))
				}

			}
		}
	}()
}
