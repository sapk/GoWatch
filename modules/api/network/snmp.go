package network

import (
	"log"
	"net"
	"regexp"
	"time"

	"github.com/sapk/GoWatch/modules/tools"
	"github.com/sapk/GoWatch/modules/watcher"
)

const snmpTimeout = 5 * time.Second

//SNMPTest execute a snmp request and return informations for testing
func SNMPTest(hostorip, community string) watcher.SNMPResponse {
	ip := hostorip

	if ok, _ := regexp.MatchString(tools.ValidIPAddressRegex, ip); !ok {
		//Si ce n'est un ip on essaie de le résoudre
		i, err := net.ResolveIPAddr("ip", hostorip)
		ip = i.String()
		if err != nil {
			log.Println("Erreur in resolving : ", err)
			return watcher.SNMPResponse{IP: ip, Result: false, Time: 0, Desc: "", Error: "hostname-unresolved"}
		}
	}
	log.Println("IP to scan ", ip)
	return watcher.SNMPTest(ip, community, snmpTimeout)

}
