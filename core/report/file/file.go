package report

import (
	"agent/api"
	"log"
)

func Log1(info interface{}){
	switch info.(type) {
	case *api.MonitorInfo :
		Info := info.(*api.MonitorInfo)
		log.Printf("Ptype: %d Pid: %d Ppid: %d Uid: %d Args: %s Path: %s namespace: %d DockerInfo: %s",
			Info.Ptype,Info.Pid,Info.Ppid,Info.Uid,Info.Args,Info.Path,Info.Ns,Info.DockerInfo)
		return
	default:
		log.Println("无法识别的类型")
	}
}