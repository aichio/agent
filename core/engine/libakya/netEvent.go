package libakya
/*
#include "akya.h"
*/
import "C"
import (
	"agent/api"
	"agent/core/engine/libakya/libakya"
	"agent/utils/log"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type NetEventEngine struct {
	interfaceFile	string
	interfaceFileFd *os.File
	eventCh 		chan api.AkyaNetEvent
}

func (self *NetEventEngine)SetInterfaceFile(file string) {
	self.interfaceFile = file
}

func (self *NetEventEngine)SetInterfaceFileFd(InterfaceFileFd *os.File) {
	self.interfaceFileFd = InterfaceFileFd
}

func (self *NetEventEngine)NewEventCh() {
	self.eventCh = make(chan api.AkyaNetEvent,1024)
}

func (self *NetEventEngine)akyaEventRead()(error){
	size,err:= libakya.GetAkyaMmapOpt(int(self.interfaceFileFd.Fd()) , libakya.AKFS_IOCTL_MMAP_GET_LEN)
	if err != nil{
		fmt.Println("GetAkyaMmapOpt failed")
		return err
	}

	data, err := syscall.Mmap(int(self.interfaceFileFd.Fd()) ,0 ,int(size) ,syscall.PROT_WRITE|syscall.PROT_READ ,syscall.MAP_SHARED)
	if err != nil{
		fmt.Println("OpenInterfaceFile failed")
		return err
	}

	ring := (* AkyaRing)(unsafe.Pointer(&data[0]))

	size,err = libakya.GetAkyaMmapOpt(int(self.interfaceFileFd.Fd()) , libakya.AKFS_IOCTL_MMAP_GET_NODE_LEN)
	if err != nil{
		fmt.Println("GetAkyaMmapOpt failed")
		return err
	}

	bt := make([]byte, size)
	for{
		rv := C.AkyaRingWait(C.int(self.interfaceFileFd.Fd()))
		if(rv <= 0){
			continue
		}

		for{
			c_char := (*C.char)(unsafe.Pointer(&bt[0]))
			rv = C.AkyaRingGet(unsafe.Pointer(ring) ,c_char ,C.int(size))
			if rv != C.int(size){
				break;
			}
			log := *((* api.AkyaNetEvent)(unsafe.Pointer(&bt[0])))
			self.eventCh <- log
		}
	}

	syscall.Munmap(data);


	return err
}

func (self *NetEventEngine)akyaEventHandle(f func(event interface{}) error){
	for{
		select{
		case x,ok := <-self.eventCh :
			if ok {
				err := f(x)
				if err!=nil {
					log.Error(-1,"err:%v",err.Error())
				}
			}else{
				log.Debug("akyaEventHandle: net eventCh close")
			}
		}
	}
}