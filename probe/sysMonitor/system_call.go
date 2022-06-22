package sysMonitor

import (
	"encoding/json"
	"syscall/datas"
	"time"
)

func StartSysCallMonitor(localAddress string, dc chan datas.ResponseData) {

	sysdata := &datas.SyscallData{
		AgentIp: localAddress,
		Pid:     "1232",
		Event:   "read",
		Time:    "114480914969",
	}
	str, _ := json.Marshal(sysdata)
	response := &datas.ResponseData{
		AgentIp:    localAddress,
		ModuleType: datas.MODULE_SYSCALL,
		Datas:      string(str),
	}
	go func() {
		for {
			dc <- *response
			time.Sleep(2 * time.Second)

		}
	}()

}
