package libakya
/*
#include "akya.h"
*/
import "C"
import (
	"agent/core/engine/libakya/libakya"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type ProcessEventEngine struct {
	interfaceFile	string
	interfaceFileFd *os.File
	eventCh 			chan libakya.AkyaProcessEvent
}

func (self *ProcessEventEngine)SetInterfaceFile(file string) {
	self.interfaceFile = file
}

func (self *ProcessEventEngine)SetInterfaceFileFd(InterfaceFileFd *os.File) {
	self.interfaceFileFd = InterfaceFileFd
}

func (self *ProcessEventEngine)NewEventCh() {
	self.eventCh = make(chan libakya.AkyaProcessEvent,1024)
}

func (self *ProcessEventEngine)akyaEventRead()(error){

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
			log := *((* libakya.AkyaProcessEvent)(unsafe.Pointer(&bt[0])))
			self.eventCh <- log
		}
	}

	syscall.Munmap(data);


	return err
}

func (self *ProcessEventEngine)akyaEventHandle(f func(event interface{}) error){
	for{
		select{
		case x,ok := <-self.eventCh :
			if ok {
				err := f(x)
				if err!=nil {
					fmt.Println(err)
				}
			}else {
				fmt.Println("process eventCh close")
			}
		}
	}
}