package ipfsi2pntcp

import (
	"context"
	i2pma "github.com/eyedeekay/sam3-multiaddr"
	"strings"

	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

// GarlicDialer implements go-libp2p-transport's Dialer interface
type GarlicDialer struct {
	garlicConn *GarlicConn
	laddr      *ma.Multiaddr
	transport  *GarlicTransport
}

// Dial connects to the specified multiaddr and returns
// a go-libp2p-transport Conn interface
func (d *GarlicDialer) Dial(raddr i2pma.I2PMultiaddr) (tpt.Conn, error) {
	netaddr, err := manet.ToNetAddr(raddr)
	var garlicAddress string
	if err != nil {
		garlicAddress, err = raddr.ValueForProtocol(ma.P_ONION)
		if err != nil {
			return nil, err
		}
	}
	garlicConn := GarlicConn{
		transport: tpt.Transport(d.transport),
		laddr:     d.laddr,
		raddr:     raddr,
	}
	if garlicAddress != "" {
		split := strings.Split(garlicAddress, ":")
		garlicConn.Conn, err = d.transport.garlicDialer.garlicConn.session.Dial("ntcp", split[0]+".b32.i2p:"+split[1])
	} else {
		garlicConn.Conn, err = d.transport.garlicDialer.garlicConn.session.Dial(netaddr.Network(), raddr.I2PAddr.Base32())
	}
	if err != nil {
		return nil, err
	}
	return &garlicConn, nil
}

func (d GarlicDialer) ClosePeer(id peer.ID) error {
	return nil
}

func (d GarlicDialer) Connectedness(id peer.ID) net.Connectedness {
	var n net.Connectedness
	return n
}

func (d GarlicDialer) Conns() []net.Conn {
	return nil
}

func (d GarlicDialer) Peers() []peer.ID {
	return nil
}

func (d GarlicDialer) Peerstore() peerstore.Peerstore {
	return nil
}

func (d GarlicDialer) ConnsToPeer(p peer.ID) []net.Conn {
	return nil
}

func (d GarlicDialer) DialPeer(ctx context.Context, p peer.ID) (net.Conn, error) {
	return nil, nil
}

func (d *GarlicDialer) DialContext(ctx context.Context, raddr i2pma.I2PMultiaddr) (tpt.Conn, error) {
	return d.Dial(raddr)
}

func (d GarlicDialer) LocalPeer() peer.ID {
	var p peer.ID
	return p
}

func (d GarlicDialer) Notify(net.Notifiee) {

}

func (d GarlicDialer) StopNotify(net.Notifiee) {

}

func (d *GarlicDialer) Matches(a ma.Multiaddr) bool {
	return IsValidGarlicMultiAddr(a)
}