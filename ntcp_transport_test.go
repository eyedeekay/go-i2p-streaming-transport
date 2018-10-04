package ipfsi2pntcp

import (
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"log"
	"testing"
)

func TestGarlicTransport(t *testing.T) {
	log.Println("\n+++ Testing ntcp_transport.go\n ")
	if key, err := createEepServiceKey(); err != nil {
		t.Fatal(err)
	} else {

		if invalidAddr, err := i2pma.NewI2PMultiaddr("/ip4/0.0.0.0/tcp/4001", true); err == nil {
			t.Fatal(err)
		} else {
			if valid := IsValidGarlicMultiAddr(invalidAddr); valid {
				t.Fatal("IsValidMultiAddr failed")
			}
		}

		if validAddr, err := i2pma.NewI2PMultiaddr("/ntcp/"+key.String(), true); err != nil {
			t.Fatal(err)
		} else {
			if valid := IsValidGarlicMultiAddr(validAddr); !valid {
				t.Fatal("IsValidMultiAddr failed")
			}
			if addr, err := validAddr.ValueForProtocol(i2pma.P_GARLIC_NTCP); err != nil {
				t.Fatal(err)
			} else {
				log.Println(addr)
			}
		}

		if transport, err := NewGarlicTransport("127.0.0.1", "7656", "", "", true); err != nil {
			t.Fatal(err)
		} else {
			log.Println(transport.keys.Addr().Base32())
		}
	}
	log.Println("\n ")
}
