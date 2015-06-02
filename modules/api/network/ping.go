package network

import (
	"github.com/sapk/GoWatch/modules/watcher"
	"log"
	"net"
	"regexp"
	"time"
)

const ping_timeout = 3 * time.Second
const ValidIpAddressRegex = "(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])"

//Ping execute a ping and return informations
func Ping(hostorip string) watcher.PingResponse {
	ip := hostorip

	if ok, _ := regexp.MatchString(ValidIpAddressRegex, ip); !ok {
		//Si ce n'est un ip on essaie de le résoudre
		i, err := net.ResolveIPAddr("ip", hostorip)
		ip = i.String()
		if err != nil {
			log.Println("Erreur in resolving : ", err)
			return watcher.PingResponse{IP: "", Result: false, Time: 0, Error: "hostname-unresolved"}
		}
	}
	//Si cela ne match toujours pas une ip c'est un echec
	if ok, _ := regexp.MatchString(ValidIpAddressRegex, ip); !ok {
		return watcher.PingResponse{IP: "", Result: false, Time: 0, Error: "hostname-unresolved"}
	}
	log.Println("IP to scan ", ip)
	ping := watcher.PingTest(ip, ping_timeout)

	return <-ping

}
