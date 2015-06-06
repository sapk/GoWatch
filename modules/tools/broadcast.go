package tools

type broadcast struct {
	c chan broadcast
	v interface{}
}

//Broadcaster a element that emit broadcast
type Broadcaster struct {
	// private fields:
	Listenc chan chan (chan broadcast)
	Sendc   chan<- interface{}
}

//BroadcastReceiver a element that receive broadcast
type BroadcastReceiver struct {
	// private fields:
	C chan broadcast
}

//NewBroadcaster create a new broadcaster object.
func NewBroadcaster() Broadcaster {
	listenc := make(chan (chan (chan broadcast)))
	sendc := make(chan interface{})
	go func() {
		currc := make(chan broadcast, 1)
		for {
			select {
			case v := <-sendc:
				if v == nil {
					currc <- broadcast{}
					return
				}
				c := make(chan broadcast, 1)
				b := broadcast{c: c, v: v}
				currc <- b
				currc = c
			case r := <-listenc:
				r <- currc
			}
		}
	}()
	return Broadcaster{
		Listenc: listenc,
		Sendc:   sendc,
	}
}

//Listen start listening to the broadcasts.
func (b Broadcaster) Listen() BroadcastReceiver {
	c := make(chan chan broadcast, 0)
	b.Listenc <- c
	return BroadcastReceiver{<-c}
}

// broadcast a value to all listeners.
func (b Broadcaster) Write(v interface{}) { b.Sendc <- v }

// read a value that has been broadcast,
// waiting until one is available if necessary.
func (r *BroadcastReceiver) Read() interface{} {
	b := <-r.C
	v := b.v
	r.C <- b
	r.C = b.c
	return v
}
