package api

import "fmt"

type EventType uint32

var (

	ProcessExec EventType= 1001
	ProcessFork EventType= 1002
	ProcessExit EventType= 1003

	FileOpen    EventType= 2001
	FileRead	EventType= 2002
	FileWrite   EventType= 2003
	FileUnlink	EventType= 2004
	FileClose	EventType= 2005
	FileChmod	EventType= 2006
	FileLink	EventType= 2007
	FileRename  EventType= 206

	NetSend EventType= 3001
	NetConn	EventType=3002
	NetAccept EventType=3003
)

func (etype EventType)String() string {
	switch etype {
	case ProcessExec:
		return "命令执行"
	case ProcessFork:
		return "进程创建"
	case ProcessExit:
		return "进程退出"
	case FileOpen:
		return "文件打开"
	case FileRead:
		return "文件读取"
	case FileWrite:
		return "文件写入"
	case FileUnlink:
		return "FileUnlink"
	case FileClose:
		return "文件关闭"
	case FileChmod:
		return "文件权限变更"
	case FileLink:
		return "FileLink"
	case NetSend:
		return "数据发送"
	case NetConn:
		return "网络外联"
	case NetAccept:
		return "收到连接"
	}
	fmt.Println(etype)
	return "unknow"
}



const ParamMaxSize  int = 256

type MonitorInfo struct {
	Ptype      EventType      `json:"ptype"`
	Pid        uint32      `json:"pid"`
	Ppid       uint32      `json:"ppid"`
	Uid        uint32      `json:"uid"`
	Ns         uint32     `json:"namespace"`
	Path       string  `json:"path"`
	File		string	`json:"file"`
	Args       string  `json:"args"`
	FileHash	string	`json:"filehash"`
	DockerInfo DockerInfo `json:"docker_info"`
}


type NetMonitorInfo struct {
	Ptype      EventType      `json:"ptype"`
	Pid        uint32      `json:"pid"`
	Ppid       uint32      `json:"ppid"`
	Uid        uint32      `json:"uid"`
	Ns         uint32     `json:"namespace"`
	Path       string  `json:"path"`
	File		string	`json:"file"`
	Args       string  `json:"args"`
	FileHash	string	`json:"filehash"`
	DockerInfo DockerInfo `json:"docker_info"`
}

type DockerInfo struct {
	ContainerID   string	`json:"container_id"`
	ContainerName string	`json:"container_name"`
}
