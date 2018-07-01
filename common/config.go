package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"
)

//Config 配置
type Config struct {
	LocalPort int      `json:"local_port"`
	Timeout   int      `json:"timeout"`
	Servers   []Server `json:"servers"`
}

//Server 服务器结构
type Server struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Method   string `json:"method"`
}

//ReadTimeout 连接超时时间
var ReadTimeout time.Duration

//ParseConfig 解析配置
func ParseConfig(path string) (config *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("读取配置文件错误")
		return
	}
	config = &Config{}
	if err = json.Unmarshal(data, &config); err != nil {
		return
	}
	ReadTimeout = time.Duration(config.Timeout) * time.Second
	return
}

//UpdateConfig 更新配置
func UpdateConfig(older, newer *Config) {
	newVal := reflect.ValueOf(newer).Elem()
	oldVal := reflect.ValueOf(older).Elem()
	log.Println(newVal, oldVal)
	return
}

//PrintVersion 当前版本
func PrintVersion() {
	const ver = "0.2.0"
	fmt.Println("当前版本: ", ver)
}
