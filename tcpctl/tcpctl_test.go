package tcpctl_test

import (
	"testing"

	"github.com/soypat/dgrams/tcpctl"
)

var (
	//	192.168.1.112	192.168.1.5	TCP	74	58920 → 80 [SYN] Seq=0 Win=64240 Len=0 MSS=1460 SACK_PERM=1 TSval=144865087 TSecr=0 WS=128
	packetSyn = []byte{0xde, 0xad, 0xbe, 0xef, 0xfe, 0xff, 0x28, 0xd2, 0x44, 0x9a, 0x2f, 0xf3, 0x08, 0x00, 0x45, 0x00,
		0x00, 0x3c, 0x2c, 0xda, 0x40, 0x00, 0x40, 0x06, 0x8a, 0x1c, 0xc0, 0xa8, 0x01, 0x70, 0xc0, 0xa8,
		0x01, 0x05, 0xe6, 0x28, 0x00, 0x50, 0x3e, 0xab, 0x64, 0xf7, 0x00, 0x00, 0x00, 0x00, 0xa0, 0x02,
		0xfa, 0xf0, 0xbf, 0x4c, 0x00, 0x00, 0x02, 0x04, 0x05, 0xb4, 0x04, 0x02, 0x08, 0x0a, 0x08, 0xa2,
		0x77, 0x3f, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x03, 0x07}
)

func TestSynReceive(t *testing.T) {
	s := tcpctl.Socket{}
	s.Listen()
	pStart, pEnd, err := s.RecvEthernet(packetSyn)
	_, _ = pStart, pEnd
	if err != nil {
		t.Fatal(err)
	}
	if pStart != pEnd {
		t.Error("expected same start/end. got ", pStart, pEnd)
	}
}
