package monitor

/*

 */
import "C"
import (
	"agent/api"
	"agent/core/engine/docker"
	libakya2 "agent/core/engine/libakya"
	"agent/core/engine/libakya/libakya"
	"agent/core/engine/rule"
	report "agent/core/report/webhook"
	"agent/utils/log"
	"bytes"
	"fmt"
)

type FileMonitor struct {
	MmapFile    string
	DockerKnow  *docker.DockerKnow
	RuleEngines *rule.FileWhiteRuleEngine
	EventEngine *libakya2.AkyaEventEngine
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

	self.EventEngine = libakya2.NewAkyaEventEngine(new(libakya2.FileEventEngine))
	self.EventEngine.NewEventEngine(self.MmapFile)
	if err != nil{
		log.Fatal(-1,"open %s,err:%s",self.MmapFile,err.Error())
		return err
	}
	return nil
}

func (self *FileMonitor) EventRead()(error) {
	fmt.Println("ProcessMonitor->EventRead")
	go self.EventEngine.EventHandle(self.analyze)
	self.EventEngine.EventRead()
	return nil
}

func (self *FileMonitor)analyze(event interface{}) (err error) {
	eventlog := event.(libakya.AkyaFileEvent)
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
