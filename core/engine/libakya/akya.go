package libakya
//
// 引用的C头文件需要在注释中声明，紧接着注释需要有import "C"，且这一行和注释之间不能有空格
//

import (
	"errors"
	"fmt"
	"os"
)

type AkyaEvevtEngineIf interface {
	SetInterfaceFile(file string)
	SetInterfaceFileFd(*os.File)
	akyaEventRead()(error)
	akyaEventHandle(f func(event interface{}) error)
	NewEventCh()
}

type AkyaEventEngine struct {
	akyaEvevtEngine AkyaEvevtEngineIf
}

type AkyaRing struct{
	offset  uint32
	size    uint32
	in      uint32
	out     uint32
}

func NewAkyaEventEngine(s AkyaEvevtEngineIf) *AkyaEventEngine {
	return &AkyaEventEngine{akyaEvevtEngine: s}
}

func (self *AkyaEventEngine)NewEventEngine(file string) (error) {
	if  file == "" {
		return errors.New("Interface filename is nil")
	}

	f, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil{
		fmt.Println("OpenInterfaceFile failed")
		return err
	}
	fmt.Println("SetInterfaceFile->",file)
	self.akyaEvevtEngine.SetInterfaceFile(file)
	self.akyaEvevtEngine.SetInterfaceFileFd(f)
	self.akyaEvevtEngine.NewEventCh()
	return nil
}

func (self *AkyaEventEngine)EventHandle(f func(event interface{}) error){
	self.akyaEvevtEngine.akyaEventHandle(f)
}

func (self *AkyaEventEngine)EventRead(){
	err := self.akyaEvevtEngine.akyaEventRead()
	if err!=nil{
		fmt.Println("akyaEvevtEngine.akyaEventRead failed. err:",err.Error())
	}
}