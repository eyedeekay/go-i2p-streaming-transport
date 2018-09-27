package ipfsi2pntcp

import (
	"github.com/eyedeekay/sam3"
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"github.com/libp2p/go-stream-muxer"
	"io"
	"net"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
)

// GarlicConn implement's go-libp2p-transport's Conn interface
type GarlicConn struct {
	net.Conn
	io.Reader
	io.Writer
	io.Closer

	transport tpt.Transport
	laddr     *ma.Multiaddr
	raddr     i2pma.I2PMultiaddr
	session   *sam3.StreamSession
}

// Close ends a SAM session associated with a transport
func (c GarlicConn) Close() error {
	//c.transport.Close()
	err := c.session.Close()
    if err == nil {
        c.session = nil
    }
	return err
}

// Transport returns the GarlicTransport associated
// with this GarlicConn
func (c GarlicConn) Transport() tpt.Transport {
	return c.transport
}

// LocalMultiaddr returns the local multiaddr for this connection
func (c GarlicConn) LocalMultiaddr() ma.Multiaddr {
	return *c.laddr
}

// RemoteMultiaddr returns the remote multiaddr for this connection
func (c GarlicConn) RemoteMultiaddr() ma.Multiaddr {
	return ma.Multiaddr(c.raddr)
}

func (c GarlicConn) IsClosed() bool {
	if c.laddr == nil {
		return true
	}
	if c.session == nil {
		return true
	}
	return false
}

func (c GarlicConn) LocalPeer() peer.ID {
	var p peer.ID
	return p
}

func (c GarlicConn) LocalPrivateKey() crypto.PrivKey {
	var pk crypto.PrivKey
	return pk
}

func (c GarlicConn) RemotePeer() peer.ID {
	var p peer.ID
	return p
}

func (c GarlicConn) RemotePublicKey() crypto.PubKey {
	var p crypto.PubKey
	return p
}

func (c GarlicConn) OpenStream() (streammux.Stream, error) {
	var s streammux.Stream
	return s, nil
}

func (c GarlicConn) AcceptStream() (streammux.Stream, error) {
	var s streammux.Stream
	return s, nil
}
