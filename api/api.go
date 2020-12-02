package api

import "fmt"

type EventType uint32

var (
	ProcessFork EventType= 1002
	ProcessExec EventType= 1001
	ProcessExit EventType= 1003

	FileCreate  EventType= 201
	FileRemove  EventType= 203
	FileOpen    EventType= 2001
	FileRead	EventType= 2002
	FileWrite   EventType= 2003
	FileUnlink	EventType= 2004
	FileClose	EventType= 2005
	FileChmod	EventType= 2006
	FileLink	EventType= 2007
	FileRename  EventType= 206
)

func (etype EventType)String() string {
	switch etype {
	case 1001:
		return "命令执行"
	case 1002:
		return "进程创建"
	case 1003:
		return "进程退出"
	case 2001:
		return "文件打开"
	case 2002:
		return "文件读取"
	case 2003:
		return "文件写入"
	case 2004:
		return "FileUnlink"
	case 2005:
		return "文件关闭"
	case 2006:
		return "文件权限变更"
	case 2007:
		return "FileLink"
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
	
	DockerInfo DockerInfo `json:"docker_info"`
}

type DockerInfo struct {
	ContainerID   string	`json:"container_id"`
	ContainerName string	`json:"container_name"`
}
