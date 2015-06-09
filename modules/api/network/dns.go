package network

import (
	"log"
	"net"
	"regexp"

	"github.com/sapk/GoWatch/modules/tools"
)

//ReverseDNSResponse represent DNS reverse response
type ReverseDNSResponse struct {
	IP     string
	Result bool
	Host   string
	Error  string
}

//ReverseDNS execute a reverse DNS
func ReverseDNS(ip string) ReverseDNSResponse {

	if ok, _ := regexp.MatchString(tools.ValidIPAddressRegex, ip); ok {
		//TODO reverse DNS if ip is IP
		log.Println("Reverse DNS : ", ip)
		hosts, err := net.LookupAddr(ip)
		log.Println("Hosts discover in reverse : ", hosts)
		if err != nil || len(hosts) == 0 {
			log.Println("Error in resolving : ", err)
			return ReverseDNSResponse{ip, false, "unknown", "Error during resolving"}
		}

		return ReverseDNSResponse{ip, true, hosts[0], ""}

	}

	return ReverseDNSResponse{ip, false, "unknown", "Invalid IP"}
}
