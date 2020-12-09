package monitor

import (
	"agent/api"
	"agent/core/engine/docker"
	libakya2 "agent/core/engine/libakya"
	"agent/core/engine/rule"
	report "agent/core/report/webhook"
	"agent/utils/log"
	"fmt"
)

type NetMonitor struct {
	MmapFile    string
	DockerKnow  *docker.DockerKnow
	RuleEngines *rule.NetRuleEngine
	EventEngine *libakya2.AkyaEventEngine
}


func (self *NetMonitor) SetMmapFile(fileName string) {
	self.MmapFile = fileName
	return
}

func (self *NetMonitor) SetDockerKnow(DockerKnow *docker.DockerKnow) {
	self.DockerKnow = DockerKnow
	return
}

func (self *NetMonitor) SetRuleEngine(RuleEngine interface{}) {
	self.RuleEngines = RuleEngine.(*rule.NetRuleEngine)
	return
}

func (self *NetMonitor) OpenMonitor()(error) {
	var err error

	self.EventEngine = libakya2.NewAkyaEventEngine(new(libakya2.NetEventEngine))
	self.EventEngine.NewEventEngine(self.MmapFile)
	if err != nil{
		log.Fatal(-1,"open %s,err:%s",self.MmapFile,err.Error())
		return err
	}
	return nil
}

func (self *NetMonitor) EventRead()(error) {
	go self.EventEngine.EventHandle(self.analyze)
	self.EventEngine.EventRead()
	return nil
}

func (self *NetMonitor)analyze(event interface{}) (err error) {
	// marshal process info
	eventlog := event.(api.AkyaNetEvent)
	// marshal process info
	info := &api.NetMonitorInfo{
		NetEvent: eventlog,
	}
	self.ResultsHandle(info)
	return nil
}

func (self *NetMonitor) ResultsHandle(value interface{}) {
	netMonitor := value.(*api.NetMonitorInfo)
	if netMonitor != nil {
		dockerinfo, _ := self.DockerKnow.Get(fmt.Sprintf("%d", netMonitor.NetEvent.Ns))
		netMonitor.DockerInfo = dockerinfo.(api.DockerInfo)
		report.Log(netMonitor)
	}
}

