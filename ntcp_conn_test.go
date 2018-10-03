package ipfsi2pntcp

import (
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"log"
	"testing"

	crypto "github.com/libp2p/go-libp2p-crypto"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
)

func TestGarlicConn(t *testing.T) {
	k, e := createEepServiceKey()
	if e != nil {
		log.Println(e)
		t.Fatal(e.Error())
	}
	raddr, e := i2pma.NewI2PMultiaddr("/ntcp/"+k.String(), true)
	if e != nil {
		t.Fatal(e.Error())
	}
	log.Println(raddr.String())

	var lPrivKey crypto.PrivKey
	var lPubKey crypto.PubKey
	var rPubKey crypto.PubKey
	var laddr ma.Multiaddr
	var transport GarlicTransport

	garlicConn, err := NewGarlicConn(
		tpt.Transport(transport),
		&laddr,
		lPrivKey,
		lPubKey,
		raddr,
		rPubKey,
	)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(garlicConn.IsClosed())
}
