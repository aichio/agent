package rule

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)


type ProcessWhiteRuleEngine struct {
	RuleEngine map[string]interface{}
	RuleEngineMutex  *sync.RWMutex
}

func CreatProcessWlRuleEngine() (*ProcessWhiteRuleEngine) {
	return &ProcessWhiteRuleEngine{
		RuleEngine:make(map[string]interface{}),
		RuleEngineMutex: new(sync.RWMutex),
	}
}

var ProcessWlRuleEngine *ProcessWhiteRuleEngine

func (this *ProcessWhiteRuleEngine)Loadjson(processWhiteConfigFile string) error {
	processWhiteConfiginfo,err := ioutil.ReadFile(processWhiteConfigFile)
	if err!=nil{
		return err
	}
	err = json.Unmarshal(processWhiteConfiginfo, &this.RuleEngine)
	return err
}


func (this *ProcessWhiteRuleEngine)Set(key string,value interface{}) {
	this.RuleEngineMutex.Lock()
	defer this.RuleEngineMutex.Unlock()
	this.RuleEngine[key] = value
}

func (this *ProcessWhiteRuleEngine)Get(key string) (interface{},bool){
	this.RuleEngineMutex.RLock()
	defer this.RuleEngineMutex.RUnlock()
	v,ok:=this.RuleEngine[key]
	return v,ok
}