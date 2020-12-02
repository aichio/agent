package docker

import (
	"agent/api"
	"sync"
)

type DockerKnow struct {
	Dict      map[string]api.DockerInfo
	DictMutex *sync.RWMutex
}

func DockerKnowNew() *DockerKnow {
	return &DockerKnow{
		Dict:      make(map[string]api.DockerInfo),
		DictMutex: new(sync.RWMutex),
	}
}

func (self *DockerKnow) Set(key string, value interface{}) {
	self.DictMutex.Lock()
	//fmt.Println(key)
	self.Dict[key] = value.(api.DockerInfo)
	self.DictMutex.Unlock()
}
func (self *DockerKnow) Get(key string) (interface{}, bool) {
	self.DictMutex.RLock()
	defer self.DictMutex.RUnlock()
	value, ok := self.Dict[key]
	return value, ok
}
