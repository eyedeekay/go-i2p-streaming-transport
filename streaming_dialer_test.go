package i2pStreaming

import (
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"log"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

func TestGarlicDialer(t *testing.T) {
	log.Println("\n+++ Testing ntcp_dialer.go\n ")
	laddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/7899")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(laddr.String())
	if transport, err := NewGarlicTransport("127.0.0.1", testPort, "", "", true); err != nil {
		t.Fatal(err)
	} else {
		log.Println("transport generated")
		raddr, err := i2pma.NewI2PMultiaddr("/ntcp/"+transport.I2PKeys.String(), true, "/sam/127.0.0.1"+testPort)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(raddr.String(), transport.I2PKeys.Addr().Base32())
		dialer, err := NewGarlicDialer(transport, laddr, raddr)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(dialer.GarlicConn.I2PMultiaddr.String(), "\n ")
	}
}
