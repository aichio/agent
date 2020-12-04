package libakya

import (
	"agent/api"
	"agent/base/lib"
	"syscall"
	"unsafe"
)

type AkyaIoctlStruct struct{
	key uint32
	size uint32
}

func GetAkyaMmapOpt(fd ,t int) (uint32,error) {
	r := AkyaIoctlStruct{}

	r.key = uint32(t)

	_, _, e := syscall.Syscall(syscall.SYS_IOCTL ,uintptr(fd) , IO('a' ,uintptr(0)) ,uintptr(unsafe.Pointer(&r)))
	if(e != 0){
		return 0,syscall.Errno(e)
	}
	return r.size,nil
}

const(
	AKFS_IOCTL_MMAP_GET_LEN      = 0
	AKFS_IOCTL_MMAP_GET_NODE_LEN = 1
	AKFS_IOCTL_MMAP_INIT_ON      = 2
	AKFS_IOCTL_MMAP_INIT_OFF     = 3
)

type AkyaProcessEvent struct{
	T       api.EventType
	Pid     uint32
	Ppid    uint32
	Uid     uint32
	Ns      uint32
	Tpath   [256]byte
	R1      [256]byte
	R2      [256]byte
}

type AkyaFileEvent struct{
	T       api.EventType
	Pid     uint32
	Ppid    uint32
	Uid     uint32
	Ns      uint32
	Tpath   [256]byte
	R1      [256]byte
	R2      [256]byte
}


type AkyaNetEvent struct{
	T       api.EventType
	Pid     uint32
	Ppid    uint32
	Uid     uint32
	Ns      uint32
	Tpath   [256]byte
	R1      AkyaNetInfo
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


