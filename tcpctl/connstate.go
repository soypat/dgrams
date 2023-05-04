package tcpctl

import (
	"sync"

	"github.com/soypat/dgrams"
)

type connState struct {
	mu sync.Mutex
	// # Send Sequence Space
	//
	//  1         2          3          4
	//  ----------|----------|----------|----------
	//  	   SND.UNA    SND.NXT    SND.UNA
	//  							+SND.WND
	//  1. old sequence numbers which have been acknowledged
	//  2. sequence numbers of unacknowledged data
	//  3. sequence numbers allowed for new data transmission
	//  4. future sequence numbers which are not yet allowed
	sndUNA int // send unacknowledged
	sndNXT int // send next
	sndWND int // send window
	sndUP  int // send urgent pointer (deprecated)
	sndWL1 int // segment sequence number used for last window update
	sndWL2 int // segment acknowledgment number used for last window update
	iss    int // initial send sequence number
	// # Receive Sequence Space
	//
	//  	1          2          3
	//  ----------|----------|----------
	//  	   RCV.NXT    RCV.NXT
	//  				 +RCV.WND
	//  1 - old sequence numbers which have been acknowledged
	//  2 - sequence numbers allowed for new reception
	//  3 - future sequence numbers which are not yet allowed
	rcvNXT int // receive next
	rcvWND int // receive window
	rcvUP  int // receive urgent pointer
	irs    int // initial receive sequence number
	state  State
}

func (cs *connState) State() State {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.state
}

func (cs *connState) incorporateFrame(hdr *dgrams.TCPHeader) {

}
