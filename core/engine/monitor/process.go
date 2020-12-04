package monitor

import "C"
import (
	"agent/api"
	"agent/base/lib"
	"agent/core/engine/docker"
	libakya2 "agent/core/engine/libakya"
	"agent/core/engine/libakya/libakya"
	"agent/core/engine/rule"
	report "agent/core/report/webhook"
	"agent/utils/log"
	"bytes"
	"fmt"
)

type ProcessMonitor struct {
	MmapFile    string
	DockerKnow  *docker.DockerKnow
	RuleEngines *rule.ProcessWhiteRuleEngine
	EventEngine *libakya2.AkyaEventEngine
}


func (self *ProcessMonitor) SetMmapFile(fileName string) {
	self.MmapFile = fileName
	return
}

func (self *ProcessMonitor) SetDockerKnow(DockerKnow *docker.DockerKnow) {
	self.DockerKnow = DockerKnow
	return
}

func (self *ProcessMonitor) SetRuleEngine(RuleEngine interface{}) {
	self.RuleEngines = RuleEngine.(*rule.ProcessWhiteRuleEngine)
	return
}

func (self *ProcessMonitor) OpenMonitor()(error) {
	var err error

	self.EventEngine = libakya2.NewAkyaEventEngine(new(libakya2.ProcessEventEngine))
	self.EventEngine.NewEventEngine(self.MmapFile)
	if err != nil{
		log.Fatal(-1,"open %s,err:%s",self.MmapFile,err.Error())
		return err
	}
	return nil
}

func (self *ProcessMonitor) EventRead()(error) {
	fmt.Println("ProcessMonitor->EventRead")
	go self.EventEngine.EventHandle(self.analyze)
	self.EventEngine.EventRead()
	return nil
}

func (self *ProcessMonitor)analyze(event interface{}) (err error) {
	eventlog := event.(libakya.AkyaProcessEvent)
	// marshal process info
	file :=string(bytes.Trim(eventlog.R1[:], "\x00"))
	filehash,err := lib.GetFileMD5(file)
	info := &api.MonitorInfo{
		Ptype: eventlog.T,
		Pid:   eventlog.Pid,
		Ppid:  eventlog.Ppid,
		Uid:   eventlog.Uid,
		Ns:    eventlog.Ns,
		File:  string(bytes.Trim(eventlog.R1[:], "\x00")),
		Args:  string(bytes.Trim(eventlog.R2[:], "\x00")),
		Path:  string(bytes.Trim(eventlog.Tpath[:], "\x00")),
		FileHash: filehash,
	}
	if _,ok := self.RuleEngines.RuleEngine[info.Path];ok {
		return
	}
	self.ResultsHandle(info)
	return nil
}

func (self *ProcessMonitor) ResultsHandle(value interface{}) {
	cprocess := value.(*api.MonitorInfo)
	if cprocess != nil {
		dockerinfo, _ := self.DockerKnow.Get(fmt.Sprintf("%d", cprocess.Ns))
		cprocess.DockerInfo = dockerinfo.(api.DockerInfo)
		report.Log(cprocess)
	}
}
