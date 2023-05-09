//go:build taptest

package dgrams_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/songgao/water"
	"github.com/soypat/dgrams"
	"github.com/soypat/dgrams/tcpctl"
)

func TestTap(t *testing.T) {
	iface, err := water.New(water.Config{
		DeviceType: water.TAP,
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf [1504]byte
	// for {
	// 	iface.Write(buf[:10])
	// }
	sock := tcpctl.Socket{}
	for {
		n, err := iface.Read(buf[:])
		if err != nil {
			t.Fatal(err)
		}
		start, end, err := sock.RecvEthernet(buf[:n])
		if err != nil {
			fmt.Println("[ERR] ", err)
		} else {
			fmt.Println("[RCV] ", string(buf[start:end]))
		}
		// err = tcpparse(buf[:n])
		// if err != nil {
		// 	fmt.Println("error:", err)
		// }
	}
	_ = iface
}

func tcpparse(buf []byte) error {
	if len(buf) < 14+20+20 {
		return errors.New("buf too small")
	}
	eth := dgrams.DecodeEthernetHeader(buf[:])
	if dgrams.EtherType(eth.SizeOrEtherType) != dgrams.EtherTypeIPv4 {
		return errors.New("only support IPv4")
	}
	ip := dgrams.DecodeIPv4Header(buf[14:])

	fmt.Printf("%+v\n%q", ip, string(buf[24:]))
	if ip.Protocol != 6 {
		return errors.New("not TCP IPv4: " + ip.String())
	}
	tcp := dgrams.DecodeTCPHeader(buf[14+20:])
	fmt.Println(tcp.String())
	return nil
}
