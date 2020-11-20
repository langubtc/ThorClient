package config

import (
	"encoding/json"
	"log"
	"os"
)

var config *ConfigJson

type ConfigJson struct {
	User       string
	Password   string
	Port       int
	IpRange    string
	ScanLogRow int
	Version    string
}

// 初始化全局配置
func LoadConfig() *ConfigJson {

	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}

	decoder := json.NewDecoder(file)
	config = &ConfigJson{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
	return config
}
