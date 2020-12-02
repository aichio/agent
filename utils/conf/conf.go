package conf

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)



var cfg *ini.File

func Init(file string) *ini.File {
	c, err := ini.Load(file)
	if err != nil {
		log.Print("local config.ini err.  errinfo=>",err)
		os.Exit(-1)
	}
	c.BlockMode = false
	return c
}

func Get(node string, key string) string {
	val := cfg.Section(node).Key(key).String()
	return val
}

func GetInt(node string, key string) int {
	val, _ := cfg.Section(node).Key(key).Int()
	return val
}

func Set(node string, key string,value string) {
	cfg.Section(node).Key(key).SetValue(value)
}



