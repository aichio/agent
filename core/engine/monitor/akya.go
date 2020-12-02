package monitor
//
// 引用的C头文件需要在注释中声明，紧接着注释需要有import "C"，且这一行和注释之间不能有空格
//

/*
#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/types.h>
#include <fcntl.h>
#include <poll.h>

typedef struct akya_security_log_s{
    unsigned int t;
    pid_t pid;
    pid_t ppid;
    uid_t uid;
    unsigned int ns;
    char tpath[256];
    char reserve1[256];
    char reserve2[256];
}__attribute__((aligned(8))) akya_security_log_t;

typedef struct akya_ring_s{
    unsigned int offset;
    unsigned int size;
    unsigned int in;
    unsigned int out;
}akya_ring_t;

#define min(x, y) ({                        \
        typeof(x) _min1 = (x);          \
        typeof(y) _min2 = (y);          \
        (void) (&_min1 == &_min2);      \
        _min1 < _min2 ? _min1 : _min2; })


unsigned int akya_ring_get(akya_ring_t *ring ,unsigned char *data ,unsigned int len)
{
    unsigned int l;
    unsigned char *buffer = (unsigned char *)ring + ring->offset;

    len = min(len, ring->in - ring->out);

    l = min(len ,ring->size - (ring->out & (ring->size - 1)));
    memcpy(data ,buffer + (ring->out & (ring->size - 1)) ,l);

    memcpy(data + l ,buffer,len - l);

    ring->out += len;

    return len;
}

int AkyaRingWait(int fd)
{
    struct pollfd fds;
    int ret;

    fds.fd = fd;
    fds.events = POLLIN;

    return poll(&fds ,1 ,5000);
}

int AkyaRingGet(void *d ,char *data ,int len)
{
	akya_ring_t *ring = (akya_ring_t *)d;
    int ret = 0;

    ret = akya_ring_get(ring ,data ,len);
    return ret;
}
*/
import "C"

import (
	"agent/core/libakya"
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type AkyaEventEngine struct {
	InterfaceFile	string
	InterfaceFileFd *os.File
	AkyaCh 			chan libakya.AkyaSecurityLogt

}

type AkyaRing struct{
	offset  uint32
	size    uint32
	in      uint32
	out     uint32
}



func NewAkyaEventEngine(file string) ( *AkyaEventEngine,error) {
	if file == "" {
		return nil,errors.New("Interface filename is nil")
	}

	f, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil{
		fmt.Println("OpenInterfaceFile failed")
		return nil,err
	}

	return &AkyaEventEngine{
		InterfaceFile:file,
		InterfaceFileFd:f,
		AkyaCh:make(chan libakya.AkyaSecurityLogt,1024),
	},nil
}



func (self *AkyaEventEngine)AkyaRun()(error){
	size,err:= libakya.GetAkyaMmapOpt(int(self.InterfaceFileFd.Fd()) , libakya.AKFS_IOCTL_MMAP_GET_LEN)
	if err != nil{
		fmt.Println("GetAkyaMmapOpt failed")
		return err
	}

	data, err := syscall.Mmap(int(self.InterfaceFileFd.Fd()) ,0 ,int(size) ,syscall.PROT_WRITE|syscall.PROT_READ ,syscall.MAP_SHARED)
	if err != nil{
		fmt.Println("OpenInterfaceFile failed")
		return err
	}

	ring := (* AkyaRing)(unsafe.Pointer(&data[0]))

	size,err = libakya.GetAkyaMmapOpt(int(self.InterfaceFileFd.Fd()) , libakya.AKFS_IOCTL_MMAP_GET_NODE_LEN)
	if err != nil{
		fmt.Println("GetAkyaMmapOpt failed")
		return err
	}

	bt := make([]byte, size)
	for{
		rv := C.AkyaRingWait(C.int(self.InterfaceFileFd.Fd()))
		if(rv <= 0){
			continue
		}

		for{
			c_char := (*C.char)(unsafe.Pointer(&bt[0]))
			rv = C.AkyaRingGet(unsafe.Pointer(ring) ,c_char ,C.int(size))
			if rv != C.int(size){
				break;
			}
			log := *((* libakya.AkyaSecurityLogt)(unsafe.Pointer(&bt[0])))
			self.AkyaCh <- log
		}
	}

	syscall.Munmap(data);


	return err
}

func (self *AkyaEventEngine)Akyahandle(f func(libakya.AkyaSecurityLogt) error){
	for{
		select{
		case x := <-self.AkyaCh :
			err := f(x)
			if err!=nil {
				fmt.Println(err)
			}
		}
	}
}