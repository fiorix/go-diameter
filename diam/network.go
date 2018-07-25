package diam

import (
	"net"
	"time"

	"github.com/ishidawataru/sctp"
)

const MaxInboundSCTPStreams = 3 // see https://tools.ietf.org/html/rfc4960#page-25
const MaxOutboundSCTPStreams = 3

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type sctpDialer struct {
	LocalAddr *sctp.SCTPAddr
}

func (d sctpDialer) Dial(network, address string) (net.Conn, error) {
	sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
	if err != nil {
		return nil, err
	}

	return sctp.DialSCTPExt(
		network,
		d.LocalAddr,
		sctpAddr,
		sctp.InitMsg{
			NumOstreams:  MaxOutboundSCTPStreams,
			MaxInstreams: MaxInboundSCTPStreams})
}

func getDialer(network string, timeout time.Duration, laddr net.Addr) Dialer {
	switch network {
	case "sctp", "sctp4", "sctp6":
		la, _ := laddr.(*sctp.SCTPAddr)
		return sctpDialer{LocalAddr: la}
	default:
		return &net.Dialer{Timeout: timeout, LocalAddr: laddr}
	}
}

func resolveAddress(network, addr string) (net.Addr, error) {
	switch network {
	case "sctp", "sctp4", "sctp6":
		return sctp.ResolveSCTPAddr(network, addr)
	case "":
		network = "tcp"
		fallthrough
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, addr)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}

func Listen(network, address string) (net.Listener, error) {
	switch network {
	case "sctp", "sctp4", "sctp6":
		sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return sctp.ListenSCTPExt(
			network,
			sctpAddr,
			sctp.InitMsg{
				NumOstreams:  MaxOutboundSCTPStreams,
				MaxInstreams: MaxInboundSCTPStreams})
	default:
		return net.Listen(network, address)
	}
}
