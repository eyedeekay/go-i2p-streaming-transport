package ipfsi2pntcp

import (
	"crypto/rand"
	"github.com/eyedeekay/sam3"
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"log"
	"net"

	crypto "github.com/libp2p/go-libp2p-crypto"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

// GarlicListener implements go-libp2p-transport's Listener interface
type GarlicListener struct {
	*sam3.StreamListener
	*sam3.StreamSession
	*GarlicTransport
	GarlicConn

	raddr    i2pma.I2PMultiaddr
	laddr    ma.Multiaddr
	lPrivKey crypto.PrivKey
	lPubKey  crypto.PubKey
}

// Accept blocks until a connection is received returning
// go-libp2p-transport's Conn interface or an error if
// something went wrong
func (l *GarlicListener) Accept() (tpt.Conn, error) {
	log.Println("GarlicListener.Accept()", l.StreamListener.Addr())
	var err error
	l.GarlicConn, err = NewGarlicConn(l.GarlicTransport, &l.laddr, l.lPrivKey, l.lPubKey, l.raddr, nil)
	if err != nil {
		return nil, err
	}
	return l.GarlicConn, nil
}

func (l *GarlicListener) Listen() (tpt.Listener, error) {
	var err error
	l.StreamListener, err = l.StreamSession.Listen()
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Close shuts down the listener
func (l *GarlicListener) Close() error {
	if err := l.StreamListener.Close(); err != nil {
		return err
	}
	if err := l.StreamSession.Close(); err != nil {
		return err
	}
	return nil
}

// Addr returns the net.Addr interface which represents
// the local multiaddr we are listening on
func (l *GarlicListener) Addr() net.Addr {
	netaddr, _ := manet.ToNetAddr(l.laddr)
	return netaddr
}

// Multiaddr returns the local multiaddr we are listening on
func (l *GarlicListener) Multiaddr() ma.Multiaddr {
	return l.laddr
}

// NewGarlicListener
func NewGarlicListener(t *GarlicTransport, key sam3.I2PKeys, laddr ma.Multiaddr) (*GarlicListener, error) {
	name := RandTunName()
	log.Println("Creating a new GarlicListener", name)
	sk, pk, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	garlicAddr, err := i2pma.NewI2PMultiaddr("/ntcp/"+key.String(), true, t.SAMAddr)
	if err != nil {
		log.Println(" \n  ", garlicAddr.String(), " \n  ")
		return nil, err
	}
	StreamSession, err := t.SAM.NewStreamSession(name, key, sam3.Options_Small)
	if err != nil {
		return nil, err
	}
	g := &GarlicListener{
		raddr:           garlicAddr,
		laddr:           laddr,
		lPrivKey:        sk,
		lPubKey:         pk,
		GarlicTransport: t,
		StreamSession:   StreamSession,
	}
	return g, nil
}
