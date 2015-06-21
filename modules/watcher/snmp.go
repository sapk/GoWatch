package watcher

import (
	"log"
	"time"

	"github.com/alouca/gosnmp"
)

//SNMPResponse hold all the information on a SNMPResponse
type SNMPResponse struct {
	IP     string
	Result bool
	Time   time.Duration
	Desc   string
	Error  string
}

//TODO monitor snmp trap

//SNMPTest try to read a proprety by SNMP
func SNMPTest(ip string, community string, timeout time.Duration) SNMPResponse {
	start := time.Now()
	s, err := gosnmp.NewGoSNMP(ip, community, gosnmp.Version2c, int64(timeout.Seconds()))
	if err != nil {
		log.Println("Error in SNMP creation : ", err)
	} else {
		resp, err := s.Get(".1.3.6.1.2.1.1.1.0")
		if err == nil {
			for _, v := range resp.Variables {
				switch v.Type {
				case gosnmp.OctetString:
					log.Printf("Response: %s : %s : %s \n", v.Name, v.Value.(string), v.Type.String())
					return SNMPResponse{IP: ip, Result: true, Time: time.Since(start), Desc: v.Value.(string), Error: ""}

				}
			}
		}
	}
	return SNMPResponse{IP: ip, Result: false, Time: time.Since(start), Desc: "", Error: "timeout"}
}
