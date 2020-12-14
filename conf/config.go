package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	GlobalConfig Config
)

func GetConfig() *Config {
	return &GlobalConfig
}

// 配置参数的格式
type Config struct {
	Server *ServerConf `yaml:"server"`

	EsConfig *ElasticSearchConf `yaml:"elasticsearch"`
}

// server信息
type ServerConf struct {
	// 监听端口号
	Port int `yaml:"port"`
}

// es服务器配置
type ElasticSearchConf struct {
	Addresses []string `yaml:"address"`
	User      string   `yaml:"user"`
	Password  string   `yaml:"password"`
}

// InitYmlFile 解析yaml配置文件
func InitYmlFile(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(f, &GlobalConfig)
	if err != nil {
		return err
	}
	return nil
}
