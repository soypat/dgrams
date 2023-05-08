package tcpctl

import (
	"fmt"

	"github.com/soypat/dgrams"
)

// State enumerates states a TCP connection progresses through during its lifetime.
//
//go:generate stringer -type=State
type State uint8

const (
	// CLOSED - represents no connection state at all.
	StateClosed State = iota
	// LISTEN - represents waiting for a connection request from any remote TCP and port.
	StateListen
	// SYN-RECEIVED - represents waiting for a confirming connection request acknowledgment
	// after having both received and sent a connection request.
	StateSynRcvd
	// SYN-SENT - represents waiting for a matching connection request after having sent a connection request.
	StateSynSent
	// ESTABLISHED - represents an open connection, data received can be delivered
	// to the user.  The normal state for the data transfer phase of the connection.
	StateEstablished
	// FIN-WAIT-1 - represents waiting for a connection termination request
	// from the remote TCP, or an acknowledgment of the connection
	// termination request previously sent.
	StateFinWait1
	// FIN-WAIT-2 - represents waiting for a connection termination request
	// from the remote TCP.
	StateFinWait2
	// CLOSING - represents waiting for a connection termination request
	// acknowledgment from the remote TCP.
	StateClosing
	// TIME-WAIT - represents waiting for enough time to pass to be sure the remote
	// TCP received the acknowledgment of its connection termination request.
	StateTimeWait
	// CLOSE-WAIT - represents waiting for a connection termination request
	// from the local user.
	StateCloseWait
	// LAST-ACK - represents waiting for an acknowledgment of the
	// connection termination request previously sent to the remote TCP
	// (which includes an acknowledgment of its connection termination request).
	StateLastAck
)

type Socket struct {
	cs connState
}

func (s *Socket) rx(hdr *dgrams.TCPHeader) {
	s.cs.mu.Lock()
	defer s.cs.mu.Unlock()
	switch s.cs.state {
	case StateClosed:
		return // Not interested in business.
	case StateListen:
		if hdr.Flags() != dgrams.FlagTCP_SYN {
			return //
		}
		fmt.Println("SYN received!")
		s.cs.rcv(hdr)
	case StateSynRcvd:

	default:
		fmt.Println("[ERR] unhandled state transition:" + s.cs.state.String())
	}
}
