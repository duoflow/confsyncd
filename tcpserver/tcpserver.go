package tcpserver

import (
	"github.com/duoflow/confsyncd/appconfig"
	"github.com/duoflow/confsyncd/loggers"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const buffersize = 1024

// SyncServer - application config
type SyncServer struct {
	FileToSync           string `json:"filetosync"`
	FileVersion          int64  `json:"fileversion"`
	SyncTimeout          int64  `json:"synctimeout"`
	Peer                 string `json:"peer"`
	Protocol             string `json:"protocol"`
	Port                 string `json:"port"`
	fieldFileNameSize    int
	fieldFileSizeSize    int
	fieldFileVersionSize int
}

// Init - create new struct object
func (s *SyncServer) Init(confData *appconfig.AppConfigStruct) {
	s.FileToSync = confData.FileToSync
	s.FileVersion = 1
	s.SyncTimeout = confData.SyncTimeout
	s.Peer = confData.Peer
	s.Protocol = confData.Protocol
	s.Port = confData.Port
}

// StartServer - listen to protocol and port
func (s *SyncServer) StartServer() {
	// create TCP-server
	server, err := net.Listen(s.Protocol, "localhost:"+s.Port)
	if err != nil {
		loggers.Error.Println("Error listetning: ", err)
		os.Exit(1)
	}
	// close server at close condition
	defer server.Close()
	// start permanent cecle loop
	loggers.Info.Println("Server started! Waiting for connections...")
	for {
		connection, err := server.Accept()
		if err != nil {
			loggers.Error.Println("Error: ", err)
			os.Exit(1)
		}
		loggers.Info.Println("Client connected")
		s.SendFileToClient(connection)
	}

}

// SendFileToClient - transfer conficuration file to peer
func (s *SyncServer) SendFileToClient(connection net.Conn) {
	loggers.Info.Println("A client has connected!")
	// close the connection in the end
	defer connection.Close()
	// read configuration file
	file, err := os.Open(s.FileToSync)
	if err != nil {
		loggers.Error.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		loggers.Error.Println(err)
		return
	}
	// prepare tech file information
	fileSize := s.fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := s.fillString(fileInfo.Name(), 256)
	fileVersion := s.fillString(strconv.FormatInt(s.FileVersion, 10), 10)
	// send technical information about file
	loggers.Info.Println("Sending fileSize, fileName and fileVersion!")
	loggers.Info.Println("FileSize", fileSize)
	loggers.Info.Println("FileName", fileName)
	loggers.Info.Println("FileVersion", fileVersion)
	//
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	connection.Write([]byte(fileVersion))
	// send the config file
	sendBuffer := make([]byte, buffersize)
	loggers.Info.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	loggers.Info.Println("File has been sent, closing connection!")
	return
}

// fillString - fill string with symbols
func (s *SyncServer) fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

// RequestFileFromServer - receive configuration file from peer
func (s *SyncServer) RequestFileFromServer() {
	connection, err := net.Dial(s.Protocol, s.Peer+":"+s.Port)
	if err != nil {
		panic(err)
	}
	loggers.Info.Println("Connected to server")
	defer connection.Close()
	//
	bufferFileName := make([]byte, s.fieldFileNameSize)
	bufferFileSize := make([]byte, s.fieldFileSizeSize)
	bufferFileVersion := make([]byte, s.fieldFileVersionSize)
	// read tech parameters - FileSize
	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	// read tech parameters - FileName
	connection.Read(bufferFileName)
	//fileName := strings.Trim(string(bufferFileName), ":")
	// read tech parameters - FileVersion
	connection.Read(bufferFileVersion)
	fileVersion, _ := strconv.ParseInt(strings.Trim(string(bufferFileVersion), ":"), 10, 64)
	if fileVersion < s.FileVersion {
		loggers.Error.Println("File version from server is older then current file version!")
		return
	}
	// rename current configuration file
	os.Rename(s.FileToSync, s.FileToSync+"_"+time.Now().String())
	// read file
	newFile, err := os.Create(s.FileToSync)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < buffersize {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+buffersize)-fileSize))
			break
		}
		io.CopyN(newFile, connection, buffersize)
		receivedBytes += buffersize
	}
	loggers.Info.Println("Received file completely!")
}

// CheckConfigurationChange - check for changes in configuration
func (s *SyncServer) CheckConfigurationChange() {
	time.Sleep(time.Second * 7)
	loggers.Info.Println("CheckConfigurationChange function")
}

// DetermineActiveVRRPNode - check keepalive daemon status
func (s *SyncServer) DetermineActiveVRRPNode() {
	addrs, err := net.InterfaceAddrs()
    if err != nil {
        loggers.Error.Println(err)
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                loggers.Info.Println(ipnet.IP.String())
            }
        }
    }
}
