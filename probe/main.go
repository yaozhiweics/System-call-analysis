package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"syscall/datas"
	"syscall/probe/appclient"

	"syscall/probe/sysMonitor"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var conf *datas.Conf
var usedAddress string
var allAddress []string

func main() {
	conf = getConf()
	usedAddress = getIfAddress(conf.IdIpPrefix)
	allAddress = getAllIfAddress(conf.AllIpPrefix)
	dataChannel := make(chan datas.ResponseData)

	sysMonitor.StartSysCallMonitor(usedAddress, dataChannel)
	go appclient.StartClient(conf.Remote, usedAddress, allAddress, dataChannel, conf)
	time.Sleep(10000 * time.Second)

}
func getConf() *datas.Conf {
	var c = new(datas.Conf)
	// /root/go/src/kraken_agent/config.yml
	///etc/kraken_agent/config.yml
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return c
}
func getAllIfAddress(prefixes []string) []string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	addresses := make([]string, 0)

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			for _, prefix := range prefixes {
				if ipnet.IP != nil && strings.HasPrefix(ipnet.IP.String(), prefix) {
					//return ipnet.IP.String()
					// addresses.Add(ipnet.IP)
					addresses = append(addresses, ipnet.IP.String())
				}
			}

		}
	}
	if len(addresses) == 0 {
		log.Errorf("no suitable address for report, please check your config")
		os.Exit(1)
	}

	return addresses
}

func getIfAddress(prefix string) string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP != nil && strings.HasPrefix(ipnet.IP.String(), prefix) {
				return ipnet.IP.String()
			}

		}
	}
	log.Errorf("no suitable address for report, please check your config")
	os.Exit(1)
	return ""
}
