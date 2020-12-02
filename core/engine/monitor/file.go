package monitor

/*

 */
import "C"
import (
	"agent/api"
	"agent/core/engine/docker"
	"agent/core/engine/rule"
	"agent/core/libakya"
	report "agent/core/report/webhook"
	"bytes"
	"fmt"
)

type FileMonitor struct {
	MmapFile    string
	DockerKnow  *docker.DockerKnow
	RuleEngines *rule.FileWhiteRuleEngine
	EventEngine *AkyaEventEngine
}

func (self *FileMonitor) SetMmapFile(fileName string) {
	self.MmapFile = fileName
	return
}

func (self *FileMonitor) SetDockerKnow(DockerKnow *docker.DockerKnow) {
	self.DockerKnow = DockerKnow
	return
}

func (self *FileMonitor) SetRuleEngine(RuleEngine interface{}) {
	self.RuleEngines = RuleEngine.(*rule.FileWhiteRuleEngine)
	return
}


func (self *FileMonitor) OpenMonitor()(error) {
	var err error
	self.EventEngine,err = NewAkyaEventEngine(self.MmapFile)
	if err != nil{
		return err
	}
	return nil
}

func (self *FileMonitor) EventRead()(error) {
	go self.EventEngine.Akyahandle(self.analyze)
	self.EventEngine.AkyaRun()
	return nil
}

func (self *FileMonitor)analyze(eventlog libakya.AkyaSecurityLogt) (err error) {
	// marshal process info
	info := &api.MonitorInfo{
		Ptype: eventlog.T,
		Pid:   eventlog.Pid,
		Ppid:  eventlog.Ppid,
		Uid:   eventlog.Uid,
		Ns:    eventlog.Ns,
		File:  fmt.Sprintf("%s",string(bytes.Trim(eventlog.R1[:], "\x00"))),
		Args:  fmt.Sprintf("%s",string(bytes.Trim(eventlog.R2[:], "\x00"))),
		Path:  fmt.Sprintf("%s",string(bytes.Trim(eventlog.Tpath[:], "\x00"))),
	}
	if _,ok := self.RuleEngines.RuleEngine[info.Path];ok {
		return
	}
	self.ResultsHandle(info)
	return nil
}

func (self *FileMonitor) ResultsHandle(value interface{}) {
	cprocess := value.(*api.MonitorInfo)
	if cprocess != nil {
		dockerinfo, _ := self.DockerKnow.Get(fmt.Sprintf("%d", cprocess.Ns))
		cprocess.DockerInfo = dockerinfo.(api.DockerInfo)
		report.Log(cprocess)
	}
}
