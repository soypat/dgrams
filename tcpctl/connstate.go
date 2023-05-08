package tcpctl

import (
	"errors"
	"sync"

	"github.com/soypat/dgrams"
)

type connState struct {
	mu sync.Mutex
	// # Send Sequence Space
	//
	//	1         2          3          4
	//	----------|----------|----------|----------
	//		   SND.UNA    SND.NXT    SND.UNA
	//								+SND.WND
	//	1. old sequence numbers which have been acknowledged
	//	2. sequence numbers of unacknowledged data
	//	3. sequence numbers allowed for new data transmission
	//	4. future sequence numbers which are not yet allowed
	snd sendSpace
	// # Receive Sequence Space
	//
	//		1          2          3
	//	----------|----------|----------
	//		   RCV.NXT    RCV.NXT
	//					 +RCV.WND
	//	1 - old sequence numbers which have been acknowledged
	//	2 - sequence numbers allowed for new reception
	//	3 - future sequence numbers which are not yet allowed
	rcv   rcvSpace
	state State
}

// sendSpace contains Send Sequence Space data.
type sendSpace struct {
	UNA uint32 // send unacknowledged
	NXT uint32 // send next
	WND uint32 // send window
	WL1 uint32 // segment sequence number used for last window update
	WL2 uint32 // segment acknowledgment number used for last window update
	iss uint32 // initial send sequence number
	UP  bool   // send urgent pointer (deprecated)
}

// rcvSpace contains Receive Sequence Space data.
type rcvSpace struct {
	NXT uint32 // receive next
	WND uint32 // receive window
	irs uint32 // initial receive sequence number
	UP  bool   // receive urgent pointer (deprecated)

}

func (cs *connState) State() State {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.state
}

func (cs *connState) frameRcv(hdr *dgrams.TCPHeader) (err error) {
	switch {
	case hdr.Ack <= cs.snd.UNA:
		err = errors.New("bad ack")
	case hdr.Ack > cs.snd.NXT:
		err = errors.New("bad ack")
	}
	if err != nil {
		return err
	}
}
