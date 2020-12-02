// +build linux

package main

import (
	"agent/base/config"
	"agent/core/engine/docker"
	"agent/core/engine/monitor"
	"agent/core/engine/rule"
	"flag"
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


	akyaConfig = config.Init(akyaConfigFile)

	dockerKnow := docker.DockerKnowNew()
	go dockerKnow.RunDockerKnow()

	if akyaConfig.Get("processEvent","enable") == "true" {
		processEventToEnable(dockerKnow)
	}

	if akyaConfig.Get("fileEvent","enable") == "true" {
		fileEventToEnable(dockerKnow)
	}
	select {}
}

func processEventToEnable(dockerKnow *docker.DockerKnow){
	RuleEngines := rule.CreatProcessWlRuleEngine()
	RuleEngines.Loadjson(processWhiteConfigFile)
	ProcessMonitorEngine:=  monitor.NewMonitorEngine(new(monitor.ProcessMonitor))
	ProcessMonitorEngine.SetMmapFile(akyaConfig.Get("processEvent","ringfile"))
	ProcessMonitorEngine.SetDockerKnow(dockerKnow)
	ProcessMonitorEngine.SetRuleEngine(RuleEngines)
	ProcessMonitorEngine.MonitorOpen()
	ProcessMonitorEngine.MonitorEventRead()
}

func fileEventToEnable(dockerKnow *docker.DockerKnow){
	RuleEngines := rule.CreatFileWlRuleEngine()
	RuleEngines.Loadjson(fileWhiteConfigFile)
	fileMonitorEngine:=  monitor.NewMonitorEngine(new(monitor.FileMonitor))
	fileMonitorEngine.SetMmapFile(akyaConfig.Get("fileEvent","ringfile"))
	fileMonitorEngine.SetDockerKnow(dockerKnow)
	fileMonitorEngine.SetRuleEngine(RuleEngines)
	fileMonitorEngine.MonitorOpen()
	fileMonitorEngine.MonitorEventRead()
}

