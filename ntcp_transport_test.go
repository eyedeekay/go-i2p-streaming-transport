package ipfsi2pntcp

import (
	"github.com/eyedeekay/sam3"
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"testing"
)

func TestNTCPConn(t *testing.T) {
	// Test valid
	key, err := createEepServiceKey()
	if err != nil {
		t.Fatal(err)
	}
	validAddr, err := i2pma.NewI2PMultiaddr("/ntcp/" + key.String())
	if err != nil {
		t.Fatal(err)
	}
	valid := IsValidGarlicMultiAddr(validAddr)
	if !valid {
		t.Fatal("IsValidMultiAddr failed")
	}

	// Test wrong protocol
	invalidAddr, err := i2pma.NewI2PMultiaddr("/ip4/0.0.0.0/tcp/4001")
	if err == nil {
		t.Fatal(err)
	}
	valid = IsValidGarlicMultiAddr(invalidAddr)
	if valid {
		t.Fatal("IsValidMultiAddr failed")
	}
}

func createEepServiceKey() (*sam3.I2PKeys, error) {
	sam, err := sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	defer sam.Close()
	k, err := sam.NewKeys()
	return &k, err
}
