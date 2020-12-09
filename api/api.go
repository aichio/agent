package api

import (
	"agent/base/lib"
	"agent/utils/log"
	"strconv"
)

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
	default:
		return strconv.Itoa(int(etype))
	}
	return "unknow"
}

type AkyaProcessEvent struct{
	T       EventType
	Pid     uint32
	Ppid    uint32
	Uid     uint32
	Ns      uint32
	Tpath   [256]byte
	R1      [256]byte
	R2      [256]byte
}

type AkyaFileEvent struct{
	T       EventType
	Pid     uint32
	Ppid    uint32
	Uid     uint32
	Ns      uint32
	Tpath   [256]byte
	R1      [256]byte
	R2      [256]byte
}


type AkyaNetEvent struct{
	T       EventType
	Pid     uint32
	Ppid    uint32
	Uid     uint32
	Ns      uint32
	Tpath   [256]byte
	NetInfo      AkyaNetInfo
}

/**
 * @brief 网络日志需要的五元组
 */
type  AkyaNetInfo struct{
	Saddr	lib.IP;
	Sport	uint32;
	Daddr	lib.IP;
	Dport	uint32;
	Protocol	Prototype;
	Hash	uint32;
}

type Prototype uint32

var (
	PrototypeTCP Prototype = 6
	PrototypeUDP Prototype = 17
)

func (Proto Prototype)String() string {
	switch Proto {
	case PrototypeTCP:
		return "TCP"
	case PrototypeUDP:
		return "UDP"
	}
	return ""
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

func (event *MonitorInfo)Log(){

	switch event.Ptype {
	case ProcessFork,ProcessExit:
		log.Info("Ptype=%s Ns=%d Pid=%d Ppid=%d Uid=%d File=%s Path=%s fileHash=%s ContainerName=%v ContainerID=%v", event.Ptype.String(),
			event.Ns ,
			event.Pid ,
			event.Ppid ,
			event.Uid,
			event.File,
			event.Path,
			event.FileHash,
			event.DockerInfo.ContainerName,
			event.DockerInfo.ContainerID)
		return
	case FileOpen,FileRead,FileWrite,FileUnlink,FileClose,FileChmod:
		return
	case ProcessExec:
		log.Info("Ptype=%s Ns=%d Pid=%d Ppid=%d Uid=%d File=%s Path=%s Args=%s fileHash=%s ContainerName=%v ContainerID=%v", event.Ptype.String(),
			event.Ns ,
			event.Pid ,
			event.Ppid ,
			event.Uid,
			event.File,
			event.Path,
			event.Args,
			event.FileHash,
			event.DockerInfo.ContainerName,
			event.DockerInfo.ContainerID)
		return
	}
}


type NetMonitorInfo struct {
	NetEvent   AkyaNetEvent
	DockerInfo DockerInfo `json:"docker_info"`
}

func (event *NetMonitorInfo)Log(){
	log.Info("Ptype=%s Ns=%d Pid=%d Ppid=%d Uid=%d File=%s Saddr=%v Sport=%d Daddr=%v Dport=%d Protocol=%v Hash=%v ContainerName=%v ContainerID=%v",
		event.NetEvent.T.String() ,
		event.NetEvent.Ns,
		event.NetEvent.Pid ,
		event.NetEvent.Ppid ,
		event.NetEvent.Uid,
		event.NetEvent.Tpath,
		event.NetEvent.NetInfo.Saddr.BigEndianPut().String(),
		lib.BigEndianPut(event.NetEvent.NetInfo.Sport),
		event.NetEvent.NetInfo.Daddr.BigEndianPut().String(),
		lib.BigEndianPut(event.NetEvent.NetInfo.Dport),
		event.NetEvent.NetInfo.Protocol.String(),
		event.NetEvent.NetInfo.Hash,
		event.DockerInfo.ContainerName,
		event.DockerInfo.ContainerID)
}


type DockerInfo struct {
	ContainerID   string	`json:"container_id"`
	ContainerName string	`json:"container_name"`
}


