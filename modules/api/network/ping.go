package network

import (
	"log"
	"net"
	"regexp"
	"time"

	"github.com/sapk/GoWatch/modules/tools"
	"github.com/sapk/GoWatch/modules/watcher"
)

const pingTimeout = 3 * time.Second

//Ping execute a ping and return informations
func Ping(hostorip string) watcher.PingResponse {
	ip := hostorip

	if ok, _ := regexp.MatchString(tools.ValidIpAddressRegex, ip); !ok {
		//Si ce n'est un ip on essaie de le résoudre
		i, err := net.ResolveIPAddr("ip", hostorip)
		ip = i.String()
		if err != nil {
			log.Println("Erreur in resolving : ", err)
			return watcher.PingResponse{IP: "", Result: false, Time: 0, Error: "hostname-unresolved"}
		}
	}
	//Si cela ne match toujours pas une ip c'est un echec
	if ok, _ := regexp.MatchString(tools.ValidIpAddressRegex, ip); !ok {
		return watcher.PingResponse{IP: "", Result: false, Time: 0, Error: "hostname-unresolved"}
	}
	log.Println("IP to scan ", ip)
	return watcher.PingTest(ip, pingTimeout)
}
