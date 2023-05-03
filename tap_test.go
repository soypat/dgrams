//go:build taptest

package dgrams_test

import (
	"fmt"
	"testing"

	"github.com/songgao/water"
	"github.com/soypat/dgrams"
)

func Test(t *testing.T) {
	iface, err := water.New(water.Config{
		DeviceType: water.TAP,
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf [1504]byte
	for {
		n, err := iface.Read(buf[:])
		if err != nil {
			t.Fatal(err)
		}
		eth := dgrams.DecodeEthernetHeader(buf[:n])
		if eth.SizeOrEtherType != uint16(dgrams.EtherTypeIPv4) {
			// fmt.Println("ignoring weird frame", eth.String()) // ignore non IPv4 frames.
			continue
		}
		fmt.Printf("(%d)interpret eth header: %s\n", n, eth.String())

	}
	_ = iface
}
