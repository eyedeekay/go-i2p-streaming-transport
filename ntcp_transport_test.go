package ipfsi2pntcp

import (
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	ma "github.com/multiformats/go-multiaddr"

	"log"
	"testing"
)

func TestGarlicTransport(t *testing.T) {
	log.Println("\n+++ Testing ntcp_transport.go\n ")
	if key, err := createEepServiceKey(); err != nil {
		t.Fatal(err)
	} else {
		log.Println("\n ++ Testing validators")
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
				log.Println("validated ntcp addr", addr)
			}
		}
		log.Println("\n ++ Testing constructors")
		if transport, err := NewGarlicTransport("127.0.0.1", "7657", "", "", true); err != nil {
			t.Fatal(err)
		} else {
			laddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/7899")
			if err != nil {
				t.Fatal(err)
			}
			log.Println("Creating connection on local", laddr.String())
			listener, err := transport.Listen(laddr)
			if err != nil {
				t.Fatal(err)
			}
			log.Println("Generated listener", listener.Addr())
			conn, err := listener.Accept()
			if err != nil {
				t.Fatal(err)
			}
			log.Println("Generated connection", conn.RemoteMultiaddr().String())
		}
	}
	log.Println("\n ")
}
