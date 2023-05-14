//go:build taptest

package dgrams_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/songgao/water"
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
	sock := tcpctl.Socket{}
	// go sendTCP()
	for {
		n, err := iface.Read(buf[:])
		if err != nil {
			t.Fatal(err)
		}
		start, end, err := sock.RecvEthernet(buf[:n])
		if err != nil {
			fmt.Print(".")
		} else {
			fmt.Println("[RCV] ", string(buf[start:end]))
		}
	}
}

func sendTCP() {
	for {
		time.Sleep(200 * time.Millisecond)
		conn, err := net.DialTimeout("tcp", "192.168.0.2:80", 10*time.Second)
		if err != nil {
			fmt.Println("Dial error:", err)
			continue
		}
		defer conn.Close()
		for {
			conn.Write([]byte("Hello world!\n"))
			time.Sleep(4 * time.Second)
		}
	}

}
