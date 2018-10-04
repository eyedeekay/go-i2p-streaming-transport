package ipfsi2pntcp

import (
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"log"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

func TestGarlicDialer(t *testing.T) {
	log.Println("\n+++ Testing ntcp_dialer.go\n ")
	var laddr ma.Multiaddr
	if transport, err := NewGarlicTransport("127.0.0.1", "7656", "", "", true); err != nil {
		t.Fatal(err)
	} else {
		raddr, err := i2pma.NewI2PMultiaddr("/ntcp/"+transport.keys.String(), true)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(transport.keys.Addr().Base32())
		dialer, err := NewGarlicDialer(transport, laddr, raddr)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(dialer.GarlicConn.raddr.String(), "\n ")
	}
}
