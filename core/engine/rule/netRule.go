package rule

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type NetRuleEngine struct {
	RuleEngine map[string]interface{}
	RuleEngineMutex  *sync.RWMutex
}

func CreatNetRuleEngine() (*NetRuleEngine) {
	return &NetRuleEngine{
		RuleEngine:make(map[string]interface{}),
		RuleEngineMutex: new(sync.RWMutex),
	}
}

var NetRuleEnginePub *NetRuleEngine

func (this *NetRuleEngine)Loadjson(ConfigFile string) error {
	processWhiteConfiginfo,err := ioutil.ReadFile(ConfigFile)
	if err!=nil{
		return err
	}
	err = json.Unmarshal(processWhiteConfiginfo, &this.RuleEngine)
	return err
}


func (this *NetRuleEngine)Set(key string,value interface{}) {
	this.RuleEngineMutex.Lock()
	defer this.RuleEngineMutex.Unlock()
	this.RuleEngine[key] = value
}

func (this *NetRuleEngine)Get(key string) (interface{},bool){
	this.RuleEngineMutex.RLock()
	defer this.RuleEngineMutex.RUnlock()
	v,ok:=this.RuleEngine[key]
	return v,ok
}
