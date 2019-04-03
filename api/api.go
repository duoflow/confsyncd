package api

import (
	"strconv"
)

// AliasInterface - struct for creating a new ip interface
type AliasInterface struct {
	SIZE           int    `json:"len"`
	DEVICE         string `json:"device"`
	NAME           string `json:"name"`
	TYPE           string `json:"type"`
	IPADDR         string `json:"ipaddr"`
	NETMASK        string `json:"netmask"`
	ONBOOT         string `json:"onboot"`
	BOOTPROTO      string `json:"bootproto"`
	BONDING_MASTER string `json:"bonding_master"`
	BONDING_OPTS   string `json:"bonding_opts"`
	FILENAME       string `json:"filename"`
}

// ConfigureSubIf - set parameters for interface
func (a *AliasInterface) ConfigureSubIf(subifIndex int, mainInterfaceName string, ipAddr string, netMask string) int {
	if len(mainInterfaceName) < 2 || len(ipAddr) < 7 {
		return 1
	}
	a.DEVICE = mainInterfaceName + ":" + strconv.Itoa(subifIndex)
	a.NAME = mainInterfaceName + ":" + strconv.Itoa(subifIndex)
	a.BONDING_MASTER = "yes"
	a.BONDING_OPTS = "mode=6 miimon=50"
	a.BOOTPROTO = "none"
	a.ONBOOT = "yes"
	a.TYPE = "bond"
	a.IPADDR = ipAddr
	a.NETMASK = netMask
	a.SIZE = 9
	a.FILENAME = "ifcfg-" //+ a.DEVICE
	return 0
}

// CreateConfigurationFile - collect interface parameters in multiline string
func (a *AliasInterface) CreateConfigurationFile() string {
	var intconfig = ``
	intconfig += "DEVICE=" + a.DEVICE + "\n"
	intconfig += "NAME=" + a.NAME + "\n"
	intconfig += "BONDING_MASTER=" + a.BONDING_MASTER + "\n"
	intconfig += "BONDING_OPTS=" + a.BONDING_OPTS + "\n"
	intconfig += "BOOTPROTO=" + a.BOOTPROTO + "\n"
	intconfig += "ONBOOT=" + a.ONBOOT + "\n"
	intconfig += "TYPE=" + a.TYPE + "\n"
	intconfig += "IPADDR=" + a.IPADDR + "\n"
	intconfig += "NETMASK=" + a.NETMASK + "\n"
	return intconfig
}
