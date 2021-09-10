package g

import (
	"log"
	"sync"
	"github.com/toolkits/file"
	"encoding/json"
)

type EtcdConfig struct {
	Host				[]string 	`json:"host"`
}

type EtcdSSL struct {
	CaFile			string			`json:"cafile"`
	CertFile		string			`json:"certfile"`
	CertKeyFile		string			`json:"certkeyfile"`
}

type EtcdSetting struct {
	LoadFile		string			`json:"loadfile"`
	EtcdPath		string			`json:'etcdpath'`
	EtcdMethod		string			`json:"method"`
}


type GlobalConfig struct {
	Debug				bool			`json:"debug"`
	EtcdConfig			*EtcdConfig		`json:"etcdconfig"`
	EtcdSSL				*EtcdSSL		`json:"etcdssl"`
	EtcdSetting			*EtcdSetting	`json:"etcdsetting"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
	Version  = "1.0.0"
)


func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}


func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
