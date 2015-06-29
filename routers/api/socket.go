package api

import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/session"
	"github.com/sapk/GoWatch/models/equipement"
	"github.com/sapk/GoWatch/modules/auth"
	"github.com/sapk/GoWatch/modules/watcher"
	"log"
	"strconv"
)

//Message pass by socket
type Message struct {
	Rep watcher.PingResponse
}

// GraphPingData send data for ping of ip in database
func PingSocket(ctx *macaron.Context, auth *auth.Auth, sess session.Store, receiver <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, errorChannel <-chan error) {
	//TODO register for ping watching
	if err := auth.VerificationAuth(ctx, sess, []string{"api.socket.ping"}); err != nil {
		return
	}
	id, _ := strconv.ParseUint(ctx.Params(":id"), 10, 64)
	eq, _ := equipement.GetByID(id)
	ping := watcher.PingRequest{false, make(chan watcher.PingResponse)}
	err := watcher.RegisterPingWatch(eq.IP(), 0, &ping)
	if err != nil {
		log.Println("Impossible to listen to ping ", err)
	}
	for {
		select {
		case p := <-ping.Ch:
			//TODO
			sender <- &Message{Rep: p}
			/*
				case msg := <-receiver:
					//TODO
					//sender <- msg
			*/
		case <-done:
			log.Println("Socket close ! ")
			ping.IsClose = true
			// the client disconnected, so you should return / break if the done channel gets sent a message
			return
		case err := <-errorChannel:
			log.Println("Error in socket : ", err)
			// Uh oh, we received an error. This will happen before a close if the client did not disconnect regularly.
		}
	}
}
