package tcpctl

import (
	"errors"
	"fmt"
	"math"

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

func (s *Socket) Listen() {
	s.cs.SetState(StateListen)
}

func (s *Socket) RecvEthernet(buf []byte) (payloadStart, payloadEnd uint16, err error) {
	buflen := uint16(len(buf))
	switch {
	case len(buf) > math.MaxUint16:
		err = errors.New("buffer too long")
	case buflen < dgrams.SizeEthernetHeaderNoVLAN+dgrams.SizeIPHeader+dgrams.SizeTCPHeaderNoOptions:
		err = errors.New("buffer too short to contain TCP")

	}
	if err != nil {
		return 0, 0, err
	}
	eth := dgrams.DecodeEthernetHeader(buf)
	if eth.IsVLAN() {
		return 0, 0, errors.New("VLAN not supported")
	}
	if eth.SizeOrEtherType != uint16(dgrams.EtherTypeIPv4) {
		return 0, 0, errors.New("support only IPv4")
	}
	ip := dgrams.DecodeIPv4Header(buf[dgrams.SizeEthernetHeaderNoVLAN:])
	payloadEnd = ip.TotalLength + dgrams.SizeEthernetHeaderNoVLAN
	if payloadEnd > buflen {
		return 0, 0, fmt.Errorf("IP.TotalLength exceeds buffer size %d/%d", payloadEnd, buflen)
	}
	if ip.Protocol != 6 { // Ensure TCP protocol.
		fmt.Printf("%+v\n%s\n", ip, ip.String())
		return 0, 0, fmt.Errorf("expected TCP protocol (6) in IP.Proto field; got %d", ip.Protocol)
	}
	tcp := dgrams.DecodeTCPHeader(buf[dgrams.SizeEthernetHeaderNoVLAN+dgrams.SizeIPHeader:])
	nb := tcp.OffsetInBytes()
	if nb < 20 {
		return 0, 0, errors.New("garbage TCP.Offset")
	}
	payloadStart = nb + dgrams.SizeEthernetHeaderNoVLAN + dgrams.SizeIPHeader
	if payloadStart > buflen {
		return 0, 0, fmt.Errorf("malformed packet, got payload offset %d/%d", payloadStart, buflen)
	}
	err = s.rx(&tcp)
	if err != nil {
		return 0, 0, err
	}
	return payloadStart, payloadEnd, err
}

func (s *Socket) rx(hdr *dgrams.TCPHeader) (err error) {
	s.cs.mu.Lock()
	defer s.cs.mu.Unlock()
	switch s.cs.state {
	case StateClosed:
		err = errors.New("connection closed")
	case StateListen:
		if hdr.Flags() != dgrams.FlagTCP_SYN {
			return //
		}
		var iss uint32 = 0 // TODO: use random start sequence when done debugging.
		fmt.Println("SYN received!")
		// Initialize connection state:
		s.cs.snd = sendSpace{
			iss: iss,
			UNA: iss,
			NXT: iss + 1,
			WND: 10,
			// UP, WL1, WL2 defaults to zero values.
		}
		s.cs.rcv = rcvSpace{
			irs: hdr.Seq,
			NXT: hdr.Seq,
			WND: hdr.WindowSize,
		}
		// Send ACK!
	case StateSynRcvd:

	default:
		err = errors.New("[ERR] unhandled state transition:" + s.cs.state.String())
		fmt.Println("[ERR] unhandled state transition:" + s.cs.state.String())
	}
	return err
}
