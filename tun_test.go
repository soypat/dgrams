//go:build tuntest

package dgrams_test

import (
	"fmt"
	"testing"

	"github.com/songgao/water"
	"github.com/soypat/dgrams/tcpctl"
)

func TestTun(t *testing.T) {
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf [1504]byte
	sock := tcpctl.Socket{}
	// go sendTCP()
	for {
		n, err := iface.Read(buf[:])
		if err != nil {
			t.Fatal(err)
		}
		start, end, err := sock.RecvTCP(buf[:n])
		if err != nil {
			fmt.Print("[err]", err, "\n\n")
		} else {
			fmt.Println("[RCV] ", string(buf[start:end]))
		}
	}
}
