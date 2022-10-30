package config

import (
	"io/ioutil"
	"sync/atomic"
	"time"

	"github.com/hust-tianbo/go_lib/log"
	"gopkg.in/yaml.v2"
)

var gConfig = &atomic.Value{}

type Config struct {
	DBUser   string `yaml:"db_user"`
	DBSecret string `yaml:"db_secret"`
	DBIP     string `yaml:"db_ip"`
}

func GetConfig() *Config {
	if configVar, ok := gConfig.Load().(*Config); ok {
		return configVar
	}
	return nil
}

// InitConfig 从配置平台拉取配置并定时更新
func InitConfig() {
	LoadAndWatch(ReadConfig)
}

// 解析配置内容
func ReadConfig(content string) error {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.DefaultLogger.Debugf("[ReadConfig]read yaml failed:%+v", err)
		return err
	}

	//fmt.Printf("[ReadConfig]read success:%+v", string(yamlFile))
	config := &Config{}

	unmarshalErr := yaml.Unmarshal(yamlFile, config)
	if unmarshalErr != nil {
		log.DefaultLogger.Debugf("[ReadConfig]unmarshal failed:%+v", unmarshalErr)
		return unmarshalErr
	}

	//fmt.Printf("[ReadConfig]unmarshal success:%+v", config)

	gConfig.Store(config)
	return nil
}

// 主动刷新+周期刷新
func LoadAndWatch(handle func(string) error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.DefaultLogger.Debugf("[LoadAndWatch]read yaml failed:%+v", err)
		panic(err)
	}

	handle(string(yamlFile)) // 刷新配置

	go func() {
		ticket := time.NewTicker(10 * time.Second)
		for _ = range ticket.C {
			yamlFile, err := ioutil.ReadFile("config.yaml")
			if err != nil {
				log.DefaultLogger.Debugf("[LoadAndWatch]read yaml failed:%+v", err)
				continue
			}

			handle(string(yamlFile)) // 刷新配置
		}
	}()

	return
}
