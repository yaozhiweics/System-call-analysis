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

//ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'yao321';

//end
