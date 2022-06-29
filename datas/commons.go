package datas

//start
type ModuleType int32
type ActionType int32

// 配置文件
type Conf struct {
	Remote      []string `yaml:"remote"`
	Token       string   `yaml:"token"`
	IdIpPrefix  string   `yaml:"id_ip_prefix"`
	AllIpPrefix []string `yaml:"all_ip_prefix"`
	Period      uint64   `yaml:"period"`
	TryNum      int      `yaml:"try_num"`
	TryTime     int      `yaml:"try_time"`
	TimeHeart   int      `yaml:"time_heart"`
	Timeout     int      `yaml:"time_out"`
	AgentIp     string   `yaml:"agent_ip"`
	SafeMode    bool     `yaml:"safe_mode"`
}

// 模块名称
const (
	MODULE_INIT    ModuleType = 0 //初始化模块
	MODULE_SYSCALL ModuleType = 1 //系统调用监控
	MODULE_CONTROL ModuleType = 2 // 消息模块
)

type ResponseData struct {
	// 相应的responseId，可以为空
	RequestId    string
	ModuleType   ModuleType // 模块类型
	AgentIp      string
	ErrCode      int
	ErrorMessage interface{}
	Datas        string
}

type SyscallData struct {
	AgentIp string `json:"agent_ip"`
	Pid     string `json:"pid"`
	Event   string `json:"event"`
	Time    string `json:"time"`
}
type Event struct {
	Timestamp           int        `json:"timestamp"`
	ThreadStartTime     int        `json:"threadStartTime"`
	ProcessorID         int        `json:"processorId"`
	ProcessID           int        `json:"processId"`
	CgroupID            uint       `json:"cgroupId"`
	ThreadID            int        `json:"threadId"`
	ParentProcessID     int        `json:"parentProcessId"`
	HostProcessID       int        `json:"hostProcessId"`
	HostThreadID        int        `json:"hostThreadId"`
	HostParentProcessID int        `json:"hostParentProcessId"`
	UserID              int        `json:"userId"`
	MountNS             int        `json:"mountNamespace"`
	PIDNS               int        `json:"pidNamespace"`
	ProcessName         string     `json:"processName"`
	HostName            string     `json:"hostName"`
	ContainerID         string     `json:"containerId"`
	ContainerImage      string     `json:"containerImage"`
	ContainerName       string     `json:"containerName"`
	PodName             string     `json:"podName"`
	PodNamespace        string     `json:"podNamespace"`
	PodUID              string     `json:"podUID"`
	EventID             int        `json:"eventId,string"`
	EventName           string     `json:"eventName"`
	ArgsNum             int        `json:"argsNum"`
	ReturnValue         int        `json:"returnValue"`
	StackAddresses      []uint64   `json:"stackAddresses"`
	Args                []Argument `json:"args"` //Arguments are ordered according their appearance in the original event
}

// Argument holds the information for one argument
type Argument struct {
	ArgMeta
	Value interface{} `json:"value"`
}

// ArgMeta describes an argument
type ArgMeta struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

//ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'yao321';

//end
