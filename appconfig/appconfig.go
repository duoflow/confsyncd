package appconfig

import (
	"encoding/json"
	"github.com/duoflow/confsyncd/loggers"
	"io/ioutil"
)

const (
	// ApplicationConfigFilePath - application configuration file
	applicationConfigFilePath = "./test/confsyncd.conf"
)

// appconfig.AppConfigStruct - application config
type AppConfigStruct struct {
	FileToSync   string `json:"filetosync"`
	SyncTimeout  int64  `json:"synctimeout"`
	Peer         string `json:"peer"`
	Protocol     string `json:"protocol"`
	Port         string `json:"port"`
	AuthPassword string `json:"authpassword"`
	VRRPVIP      string `json:"vrrpvip"`
	VRRPMasterFlag string 
}

// ReadConfig - read configuration from filesystem
func (a *AppConfigStruct) ReadConfig() int {
	// read file
	data, err := ioutil.ReadFile(applicationConfigFilePath)
	if err != nil {
		loggers.Error.Println(err)
		return 1
	}
	// parse config file
	err = json.Unmarshal(data, a)
	if err != nil {
		loggers.Error.Println(err)
		return 2
	}
	return 0
}
