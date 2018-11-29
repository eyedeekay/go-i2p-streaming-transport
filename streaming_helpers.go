package i2pStreaming

import (
	"fmt"
	"github.com/eyedeekay/sam3"
	"math/rand"

	i2pma "github.com/eyedeekay/sam3-multiaddr"
	//ma "github.com/multiformats/go-multiaddr"
)

// IsValidGarlicMultiAddr is used to validate that a multiaddr
// is representing a I2P garlic service
func IsValidGarlicMultiAddr(a i2pma.I2PMultiaddr) bool {
	if len(a.Protocols()) < 2 {
		return false
	}

	// check for correct network type
	if a.Protocols()[0].Name != "ntcp" {
		fmt.Println("Protocol != ntcp")
		return false
	}

	// split into garlic address
	addr, err := a.ValueForProtocol(i2pma.P_GARLIC_NTCP)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//kinda crude, but if it's bigger than this it's at least possible that
	//it's a valid kind of i2p address.
	if len(addr) < 30 {
		fmt.Println(addr)
		return false
	}

	return true
}

// RandTunName generates a random tunnel names to avoid collisions
func RandTunName() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = "abcdefghijklmnopqrstuvwxyz"[rand.Intn(len("abcdefghijklmnopqrstuvwxyz"))]
	}
	return string(b)
}

func createEepServiceKey() (*sam3.I2PKeys, error) {
	sam, err := sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	defer sam.Close()
	k, err := sam.NewKeys()
	if err != nil {
		return nil, err
	}
	return &k, err
}
