package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

type AkyaConfig struct {
	Cfg *ini.File
}

func Init(file string) (*AkyaConfig) {
	c, err := ini.Load(file)
	if err != nil {
		log.Print("local config.ini err.  errinfo=>",err)
		os.Exit(-1)
	}
	c.BlockMode = false

	return &AkyaConfig{
		Cfg: c,
	}
}

func (self *AkyaConfig)Get(node string, key string) string {
	val := self.Cfg.Section(node).Key(key).String()
	return val
}

func (self *AkyaConfig)GetInt(node string, key string) int {
	val, _ := self.Cfg.Section(node).Key(key).Int()
	return val
}

/*
func Set(node string, key string,value string) {
	cfg.Section(node).Key(key).SetValue(value)
}
*/
