package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
)

//Config 配置
type Config struct {
	Server      string `json:"server"`
	Port        int    `json:"port"`
	LocalServer string `json:"local_server"`
	LocalPort   int    `json:"local_port"`
	Password    string `json:"password"`
	Method      string `json:"method"`
	Timeout     int    `json:"timeout"`
}

//ReadTimeout 连接超时时间
var ReadTimeout time.Duration

//ParseConfig 解析配置
func ParseConfig(path string) (config *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("config.json not found!")
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read config.json error")
		return
	}
	config = &Config{}
	if err = json.Unmarshal(data, &config); err != nil {
		return
	}
	if config.Password == "" {
		err = errors.New("Config password cannot be empty")
		return
	}
	if config.Method == "" {
		err = errors.New("Config method cannot be empty")
		return
	}
	ReadTimeout = time.Duration(config.Timeout) * time.Second
	return
}
