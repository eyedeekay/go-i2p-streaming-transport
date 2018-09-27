package ipfsi2pntcp

import (
	"github.com/eyedeekay/sam3"
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"net"

	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

// GarlicListener implements go-libp2p-transport's Listener interface
type GarlicListener struct {
	key       i2pma.I2PMultiaddr
	laddr     ma.Multiaddr
	session   *sam3.StreamSession
	listener  net.Listener
	transport tpt.Transport
}

// Accept blocks until a connection is received returning
// go-libp2p-transport's Conn interface or an error if
// something went wrong
func (l *GarlicListener) Accept() (tpt.Conn, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}
	raddr, err := manet.FromNetAddr(conn.RemoteAddr())
	if err != nil {
		return nil, err
	}
	garlicConn := GarlicConn{
		Conn:      conn,
		transport: l.transport,
		laddr:     &l.laddr,
		raddr:     raddr.(i2pma.I2PMultiaddr),
	}
	return &garlicConn, nil
}

// Close shuts down the listener
func (l *GarlicListener) Close() error {
	return l.listener.Close()
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
