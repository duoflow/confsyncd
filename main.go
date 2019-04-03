package main

import (
	//"fmt"
	"io/ioutil"
	"os"
	//"sync"

	//"github.com/duoflow/systemtasks/api"
	"github.com/duoflow/confsyncd/appconfig"
	"github.com/duoflow/confsyncd/loggers"
	"github.com/duoflow/confsyncd/tcpserver"
	//"github.com/duoflow/confsyncd/webserver"
)

func main() {
	// make loggers initialization
	loggers.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	// read application configuration file
	var appConfig appconfig.AppConfigStruct
	appConfig.ReadConfig()
	loggers.Info.Println(appConfig)

	// start TCP-server to listen to connection
	var tcpsrv tcpserver.SyncServer
	tcpsrv.Init(&appConfig)
	go tcpsrv.StartServer()
	// start client application
	for {
		tcpsrv.CheckConfigurationChange()
		tcpsrv.DetermineActiveVRRPNode()
	}
}
