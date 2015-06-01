package network

import (
	"time"
	"github.com/sapk/GoWatch/modules/watcher"
)

const maxtimeout = 3*time.Second
// PingResponse represent ip response and stats
type PingResponse struct {
	IP     string
	Result bool
	Time   time.Duration
}

//Ping execute a ping and return informations
func Ping(ip string) PingResponse {
	
	//TODO
	//TODO break on timeout
	ch := watcher.RegisterPingWatch(ip,maxtimeout)
	//defer close(ch)
	watcher.SendPing(ip)
	start := time.Now()
	
	timeout := make(chan bool, 1)
	go func() {
	    time.Sleep(maxtimeout)
	    timeout <- true
	}()
	//defer close(timeout)
	
	select {
		case <-ch:
		    // a read from ch has occurred
			return PingResponse{IP: ip, Result: true, Time: time.Since(start)}
		case <-timeout:
		    // the read from ch has timed out
			return PingResponse{IP: ip, Result: false, Time: time.Since(start)}
	}
}
