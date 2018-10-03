package ipfsi2pntcp

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/eyedeekay/sam3"
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"golang.org/x/net/proxy"
	"os"
	"path/filepath"
	"strings"

	crypto "github.com/libp2p/go-libp2p-crypto"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
)

// GarlicTransport implements go-libp2p-transport's Transport interface
type GarlicTransport struct {
	SAMConn        *sam3.SAM
	garlicDialer   *GarlicDialer
	garlicListener *GarlicListener
	keysDir        string
	keys           *sam3.I2PKeys
}

func (t GarlicTransport) CanDial(m ma.Multiaddr) bool {
	return t.garlicDialer.Matches(m)
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

// loadKeys loads keys into our keys from files in the keys directory
func (t *GarlicTransport) loadKeys() (*sam3.I2PKeys, error) {
	var keys sam3.I2PKeys
	absPath, err := filepath.EvalSymlinks(t.keysDir)
	if err != nil {
		return nil, err
	}
	walkpath := func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".i2pkeys") {
			file, err := os.Open(path)
			defer file.Close()
			if err != nil {
				return err
			}
			privKey, err := sam3.LoadKeysIncompat(file)
			if err != nil {
				return err
			}
			keys = privKey
		}
		return nil
	}
	err = filepath.Walk(absPath, walkpath)
	return &keys, err
}

// Dialer creates and returns a go-libp2p-transport Dialer
func (t GarlicTransport) Dialer(laddr ma.Multiaddr) (net.Dialer, error) {
	dialer, err := NewGarlicDialer(&t, laddr)
	if err != nil {
		return nil, err
	}
	return dialer, nil
}

// Listen creates and returns a go-libp2p-transport Listener
func (t GarlicTransport) Listen(laddr ma.Multiaddr) (tpt.Listener, error) {

	garlicAddr, err := i2pma.NewI2PMultiaddr(t.keys.String(), true)
	if err != nil {
		return nil, fmt.Errorf("failed to generate I2PMultiaddr")
	}

	sk, pk, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}

	listener := GarlicListener{
		key:       garlicAddr,
		laddr:     laddr,
		lPrivKey:  sk,
		lPubKey:   pk,
		transport: &t,
	}

	listener.session, err = t.SAMConn.NewStreamSession(RandTunName(), *t.keys, sam3.Options_Medium)
	if err != nil {
		return nil, err
	}

	listener.StreamListener, err = listener.session.Listen()
	if err != nil {
		return nil, err
	}

	return &listener, nil
}

// Matches returns true if the address is a valid onion multiaddr
func (t *GarlicTransport) Matches(a ma.Multiaddr) bool {
	return IsValidGarlicMultiAddr(a)
}

// NewGarlicTransport initializes a GarlicTransport for libp2p
func NewGarlicTransport(SAMAddr, SAMPort, SANPass string, auth *proxy.Auth, keysDir string, onlyGarlic bool) (*GarlicTransport, error) {
	conn, err := sam3.NewSAM(SAMAddr + ":" + SAMPort)
	if err != nil {
		return nil, err
	}
	g := GarlicTransport{
		SAMConn: conn,
		keysDir: keysDir,
	}
	/*g.garlicDialer, err = NewGarlicDialer(&g)
	if err != nil {
		return nil, err
	}*/
	keys, err := g.loadKeys()
	if err != nil {
		return nil, err
	}
	g.keys = keys
	g.garlicListener.session, err = conn.NewStreamSession(RandTunName(), *g.keys, sam3.Options_Small)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
