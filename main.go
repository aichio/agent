// +build linux

package main

import (
	"agent/base/config"
	"agent/base/lib"
	"agent/core/engine/docker"
	"agent/core/engine/monitor"
	"agent/core/engine/rule"
	"flag"
	"fmt"
)

var(
	akyaConfigFile string
	processWhiteConfigFile string
	fileWhiteConfigFile		string
)

func init(){
	flag.StringVar(&akyaConfigFile,"conf", "config.ini",
		`配置信息.`)
	flag.StringVar(&processWhiteConfigFile,"pwlc", "processwhite.json",
		`进程白名单配置.`)
	flag.StringVar(&fileWhiteConfigFile,"fwlc", "filewhite.json",
		`文件白名单配置.`)
}

var akyaConfig *config.AkyaConfig

func main() {
	flag.Parse()
	//加载容器信息
	defer lib.TryE()

	akyaConfig = config.Init(akyaConfigFile)

	dockerKnow := docker.DockerKnowNew()
	go dockerKnow.RunDockerKnow()

	if akyaConfig.Get("processMonitor","enable") == "true" {
		fmt.Println("processMonitor enable")
		go processEventToEnable(dockerKnow)
	}

	if akyaConfig.Get("fileMonitor","enable") == "true" {
		fmt.Println("fileMonitor enable")
		go fileEventToEnable(dockerKnow)
	}
	if akyaConfig.Get("netMonitor","enable") == "true" {
		fmt.Println("netMonitor enable")
		go netMonitorToEnable(dockerKnow)
	}

	select {}
}

func processEventToEnable(dockerKnow *docker.DockerKnow){
	defer lib.TryE()
	RuleEngines := rule.CreatProcessWlRuleEngine()
	RuleEngines.Loadjson(processWhiteConfigFile)
	ProcessMonitorEngine:=  monitor.NewMonitorEngine(new(monitor.ProcessMonitor))
	ProcessMonitorEngine.SetMmapFile(akyaConfig.Get("processMonitor","interfaceFile"))
	ProcessMonitorEngine.SetDockerKnow(dockerKnow)
	ProcessMonitorEngine.SetRuleEngine(RuleEngines)
	ProcessMonitorEngine.MonitorOpen()
	ProcessMonitorEngine.MonitorEventRead()
}

func fileEventToEnable(dockerKnow *docker.DockerKnow){
	defer lib.TryE()
	RuleEngines := rule.CreatFileWlRuleEngine()
	RuleEngines.Loadjson(fileWhiteConfigFile)
	fileMonitorEngine:=  monitor.NewMonitorEngine(new(monitor.FileMonitor))
	fileMonitorEngine.SetMmapFile(akyaConfig.Get("fileMonitor","interfaceFile"))
	fileMonitorEngine.SetDockerKnow(dockerKnow)
	fileMonitorEngine.SetRuleEngine(RuleEngines)
	fileMonitorEngine.MonitorOpen()
	fileMonitorEngine.MonitorEventRead()
}

func netMonitorToEnable(dockerKnow *docker.DockerKnow){
	defer lib.TryE()
	MonitorEngine:=  monitor.NewMonitorEngine(new(monitor.NetMonitor))
	MonitorEngine.SetMmapFile(akyaConfig.Get("netMonitor","interfaceFile"))
	MonitorEngine.SetDockerKnow(dockerKnow)
	MonitorEngine.MonitorOpen()
	MonitorEngine.MonitorEventRead()
}
