package report

import (
	"agent/api"
	"agent/base/lib"
	"agent/core/engine/libakya/libakya"
	"agent/utils/log"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var serverIp *string
var serverPort *string
var enable *bool

var logFile string
func init() {
	enable = flag.Bool("webhook", false, "webhook enable.")
	serverIp = flag.String("server", "", "server Ip.")
	serverPort = flag.String("port", "", "server Port.")

	flag.StringVar(&logFile,"log", "/var/log/akya.log", "指定日志文件.")
	loginit(logFile)
}

func loginit(fileName string) {
	mode := "file"
	config :=fmt.Sprintf( `{"level":0,"filename":"%s"}`,fileName)
	log.NewLogger(0, mode, config)
}

func Log(info interface{}){

	switch info.(type) {
	case *api.MonitorInfo :
		Info := info.(*api.MonitorInfo)

		if info.(*api.MonitorInfo).Ppid == 0{
			return
		}
		event := info.(*api.MonitorInfo)
		log.Info("Ptype=%s Ns=%d Pid=%d Ppid=%d Uid=%d File=%s Path=%s Args=%s fileHash=%s DockerInfo=%v", event.Ptype.String(),
			event.Ns ,
			event.Pid ,
			event.Ppid ,
			event.Uid,
			event.File,
			event.Path,
			event.Args,
			event.FileHash,
			event.DockerInfo)
		if *enable {
			WebHook(*serverIp,*serverPort,Info)
		}
		return
	case libakya.AkyaNetEvent:
		event := info.(libakya.AkyaNetEvent)
		log.Info("Ptype=%s Ns=%d Pid=%d Ppid=%d Uid=%d File=%s Saddr=%v Sport=%d Daddr=%v Dport=%d Protocol=%v Hash=%v",event.T.String(),
			event.Ns ,
			event.Pid ,
			event.Ppid ,
			event.Uid,
			event.Tpath,
			event.R1.Saddr.BigEndianPut().String(),
			lib.BigEndianPut(event.R1.Sport),
			event.R1.Daddr.BigEndianPut().String(),
			lib.BigEndianPut(event.R1.Dport),
			event.R1.Protocol.String(),
			event.R1.Hash)
	default:
		log.Debug("无法识别的类型")
	}
}

func WebHook(destip string, destport string, ReportContent interface{}) {
	// 判断 WebHook 通知
	bytesData, _ := json.Marshal(ReportContent)
	reader := bytes.NewReader(bytesData)
	
	request, _ := http.NewRequest("POST", "http://"+destip+":"+destport+"/Monitor", reader)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		log.Debug("%v  上报记录失败. err:(%v)\n", destip, err)
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug("WebHook read resp body fail.err: (%v)\n", err)
	}
	//log.Printf("回应：%v\n", string(body))
}