package controls

import (
	"fmt"
	"syscall/datas"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var Db *sqlx.DB

func InitDB() {
	database, err := sqlx.Open("mysql", "root:yao321@tcp(127.0.0.1:3306)/syscall")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database
	//defer Db.Close() // 注意这行代码要写在上面err判断的下面
}

func Insert(arg *datas.SyscallData) int {
	sql := "insert into system_call(agent_ip,pid,event,time) values (?,?,?,?)"
	//2.预编译
	strTmt, err := Db.Prepare(sql)
	if err != nil {
		log.Errorln("Prepare fail:", err)
		return -1
	}
	//执行sql
	r, err := strTmt.Exec(arg.AgentIp, arg.Pid, arg.Event, arg.Time) //id自动增长
	if err != nil {
		log.Errorln("Exec fail:", err)
		return -1
	}
	//拿到插入数据的id
	id, err := r.LastInsertId()
	if err != nil {
		log.Errorln("Exec fail:", err)
		return -1
	}
	return int(id)
}
