package ipfsi2pntcp

import (
	"context"
	"github.com/eyedeekay/sam3"
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"os"
	"path/filepath"
	"strings"

	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
)

// GarlicTransport implements go-libp2p-transport's Transport interface
type GarlicTransport struct {
	*sam3.SAM
	*sam3.I2PKeys
	*GarlicListener
	*GarlicDialer
	SAMAddr string

	keysPath string
}

func (t GarlicTransport) CanDial(m ma.Multiaddr) bool {
	return t.GarlicDialer.Matches(m)
}

func (t GarlicTransport) CanDialI2P(m i2pma.I2PMultiaddr) bool {
	return t.GarlicDialer.MatchesI2P(m)
}

func (t GarlicTransport) Dial(c context.Context, m ma.Multiaddr, p peer.ID) (tpt.Conn, error) {
	var conn GarlicConn
	return conn, nil
}

// Protocols need only return this I think
func (t GarlicTransport) Protocols() []int {
	return []int{i2pma.P_GARLIC_NTCP}
}

// Proxy always returns false, we're using the SAM bridge to make our requests
func (t GarlicTransport) Proxy() bool {
	return false
}

// Dialer creates and returns a go-libp2p-transport Dialer
func (t GarlicTransport) Dialer(laddr ma.Multiaddr) (net.Dialer, error) {
	var err error
	t.GarlicDialer, err = NewGarlicDialer(&t, laddr, t.GarlicDialer.I2PMultiaddr)
	if err != nil {
		return nil, err
	}
	return t.GarlicDialer, nil
}

// Listen creates and returns a go-libp2p-transport Listener
func (t GarlicTransport) Listen(laddr ma.Multiaddr) (tpt.Listener, error) {
	var err error
	//name := RandTunName()
	t.GarlicListener, err = NewGarlicListener(&t, *t.I2PKeys, laddr)
	if err != nil {
		return nil, err
	}
	return t.GarlicListener.Listen()
}

// Matches returns true if the address is a valid onion multiaddr
func (t *GarlicTransport) Matches(a i2pma.I2PMultiaddr) bool {
	return IsValidGarlicMultiAddr(a)
}

// loadKeys loads keys into our keys from files in the keys directory
func LoadKeys(keysPath string) (*sam3.I2PKeys, error) {
	absPath, err := filepath.EvalSymlinks(keysPath)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(absPath, ".i2pkeys") {
		file, err := os.Open(absPath)
		defer file.Close()
		if err != nil {
			return nil, err
		}
		keys, err := sam3.LoadKeysIncompat(file)
		if err != nil {
			return nil, err
		}
		return &keys, nil
	}

	return createEepServiceKey()
}

// NewGarlicTransport initializes a GarlicTransport for libp2p
func NewGarlicTransport(SAMAddr, SAMPort, SANPass string, keysPath string, onlyGarlic bool) (*GarlicTransport, error) {
	conn, err := sam3.NewSAM(SAMAddr + ":" + SAMPort)
	if err != nil {
		return nil, err
	}
	keys, err := LoadKeys(keysPath)
	if err != nil {
		return nil, err
	}
	g := GarlicTransport{
		SAMAddr:        "/sam/" + SAMAddr + ":" + SAMPort,
		SAM:            conn,
		I2PKeys:        keys,
		GarlicListener: &GarlicListener{},
		GarlicDialer:   &GarlicDialer{},
		keysPath:       keysPath,
	}
	return &g, nil
}
