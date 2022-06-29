package sysMonitor

import (
	"bufio"
	"io"
	"os/exec"
	"sync"
	"syscall/datas"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartSysCallMonitor(localAddress string, dc chan datas.ResponseData) {
	cmd := "/home/yao/syscall/probe/dist/tracee-ebpf --output json"

	c := exec.Command("/bin/bash", "-c", cmd) // mac or linux
	stdout, err := c.StdoutPipe()
	if err != nil {
		log.Error(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				log.Error(err)
				return
			}
			//fmt.Print(readString)
			//组装response ，写入chan
			response := &datas.ResponseData{
				AgentIp:    localAddress,
				ModuleType: datas.MODULE_SYSCALL,
				Datas:      readString,
			}
			dc <- *response
			time.Sleep(100 * time.Millisecond)

		}
	}()
	err = c.Start()
	wg.Wait()

	// sysdata := &datas.SyscallData{
	// 	AgentIp: localAddress,
	// 	Pid:     "1232",
	// 	Event:   "read",
	// 	Time:    "114480914969",
	// }
	// str, _ := json.Marshal(sysdata)
	// response := &datas.ResponseData{
	// 	AgentIp:    localAddress,
	// 	ModuleType: datas.MODULE_SYSCALL,
	// 	Datas:      string(str),
	// }
	// go func() {
	// 	for {
	// 		dc <- *response
	// 		time.Sleep(2 * time.Second)

	// 	}
	// }()

}
