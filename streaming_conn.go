package i2pStreaming

import (
	"github.com/eyedeekay/sam3"
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"github.com/libp2p/go-stream-muxer"
	"log"
	n "net"
	"time"

	crypto "github.com/libp2p/go-libp2p-crypto"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
)

// GarlicConn implements go-libp2p-transport's Conn interface
type GarlicConn struct {
	conn n.Conn
	i2pma.I2PMultiaddr
	*sam3.StreamSession
	*GarlicTransport

	laddr    *ma.Multiaddr
	lPrivKey crypto.PrivKey
	lPubKey  crypto.PubKey

	rPubKey crypto.PubKey
}

// Close ends a SAM session associated with a transport
func (c GarlicConn) Close() error {
	err := c.StreamSession.Close()
	if err == nil {
		c.StreamSession = nil
	}
	return err
}

// Conn converts a GarlicConn to a net.Conn
func (c GarlicConn) Conn() net.Conn {
	return c
}

// Transport returns the GarlicTransport associated
// with this GarlicConn
func (c GarlicConn) Transport() tpt.Transport {
	return c.GarlicTransport
}

// LocalMultiaddr returns the local multiaddr for this connection
func (c GarlicConn) LocalMultiaddr() ma.Multiaddr {
	return *c.laddr
}

// RemoteMultiaddr returns the remote multiaddr for this connection
func (c GarlicConn) RemoteMultiaddr() ma.Multiaddr {
	return ma.Multiaddr(c.I2PMultiaddr)
}

// IsClosed says a connection is close if a session hasn't been opened for now.
func (c GarlicConn) IsClosed() bool {
	if c.laddr == nil {
		return true
	}
	if c.StreamSession == nil {
		return true
	}
	return false
}

// LocalPeer returns the local peer.ID used for IPFS
func (c GarlicConn) LocalPeer() peer.ID {
	lpeer, _ := peer.IDFromPrivateKey(c.LocalPrivateKey())
	return lpeer
}

// LocalPrivateKey returns the local private key used for the peer.ID
func (c GarlicConn) LocalPrivateKey() crypto.PrivKey {
	return c.lPrivKey
}

// RemotePeer returns the remote peer.ID used for IPFS
func (c GarlicConn) RemotePeer() peer.ID {
	rpeer, _ := peer.IDFromPublicKey(c.RemotePublicKey())
	return rpeer
}

//RemotePublicKey returns the remote public key used for the peer.ID
func (c GarlicConn) RemotePublicKey() crypto.PubKey {
	return c.rPubKey
}

//Read finishes implementing a net.Conn
func (c GarlicConn) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

//Write finishes implementing a net.Conn
func (c GarlicConn) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

// Reset lets us streammux
func (c GarlicConn) Reset() error {
	return c.Close()
}

//
func (c GarlicConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

//
func (c GarlicConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

//
func (c GarlicConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

//
func (c GarlicConn) Protocol() protocol.ID {
	r := protocol.ID("/sam/")
	return r
}

//
func (c GarlicConn) SetProtocol(protocol.ID) {

}

// OpenStream lets us streammux
func (c GarlicConn) OpenStream() (streammux.Stream, error) {
	return c, nil
}

// AcceptStream lets us streammux
func (c GarlicConn) AcceptStream() (streammux.Stream, error) {
	return c, nil
}

func (c GarlicConn) NewStream() (net.Stream, error) {
	r := c
	return r, nil
}

// GetStreams lets us streammux
func (c GarlicConn) GetStreams() []net.Stream {
	var r []net.Stream
	return r
}

func (c GarlicConn) Stat() net.Stat {
	var r net.Stat
	return r
}

func NewGarlicConn(t *GarlicTransport, laddr *ma.Multiaddr, lPrivKey crypto.PrivKey, lPubKey crypto.PubKey, raddr i2pma.I2PMultiaddr, rPubKey crypto.PubKey) (GarlicConn, error) {
	name := RandTunName()
	log.Println("Creating a new GarlicConn", name)
	garlicConn := GarlicConn{
		GarlicTransport: t,
		laddr:           laddr,
		lPrivKey:        lPrivKey,
		lPubKey:         lPubKey,
		I2PMultiaddr:    raddr,
		rPubKey:         rPubKey,
	}
	conn, err := sam3.NewSAM(raddr.SAMAddressString())
	if err != nil {
		return garlicConn, nil
	}
	keys, err := conn.NewKeys()
	if err != nil {
		return garlicConn, nil
	}
	garlicConn.StreamSession, err = conn.NewStreamSession(name, keys, sam3.Options_Small)
	if err != nil {
		return garlicConn, nil
	}
	return garlicConn, nil
}
