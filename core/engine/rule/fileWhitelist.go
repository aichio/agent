package rule

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type FileWhiteRuleEngine struct {
	RuleEngine map[string]interface{}
	RuleEngineMutex  *sync.RWMutex
}

func CreatFileWlRuleEngine() (*FileWhiteRuleEngine) {
	return &FileWhiteRuleEngine{
		RuleEngine:make(map[string]interface{}),
		RuleEngineMutex: new(sync.RWMutex),
	}
}

var FileWlRuleEngine *FileWhiteRuleEngine

func (this *FileWhiteRuleEngine)Loadjson(WhiteConfigFile string) error {
	WhiteConfiginfo,err := ioutil.ReadFile(WhiteConfigFile)
	if err!=nil{
		return err
	}
	err = json.Unmarshal(WhiteConfiginfo, &this.RuleEngine)
	return err
}


func (this *FileWhiteRuleEngine)Set(key string,value interface{}) {
	this.RuleEngineMutex.Lock()
	defer this.RuleEngineMutex.Unlock()
	this.RuleEngine[key] = value
}

func (this *FileWhiteRuleEngine)Get(key string) (interface{},bool){
	this.RuleEngineMutex.RLock()
	defer this.RuleEngineMutex.RUnlock()
	v,ok:=this.RuleEngine[key]
	return v,ok
}