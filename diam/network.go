package diam

import (
	"net"
	"time"

	"github.com/ishidawataru/sctp"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type sctpDialer struct {
	laddr string
}

func (d sctpDialer) Dial(network, raddr string) (net.Conn, error) {
	var (
		sctpLAddr *sctp.SCTPAddr
		err       error
	)

	if d.laddr != "" {
		sctpLAddr, err = sctp.ResolveSCTPAddr(network, d.laddr)
		if err != nil {
			return nil, err
		}
	}

	sctpRAddr, err := sctp.ResolveSCTPAddr(network, raddr)
	if err != nil {
		return nil, err
	}
	return sctp.DialSCTP(network, sctpLAddr, sctpRAddr)
}

func getDialer(network string, laddr string, timeout time.Duration) Dialer {
	switch network {
	case "sctp", "sctp4", "sctp6":
		return sctpDialer{laddr: laddr}
	default:
		return &net.Dialer{Timeout: timeout}
	}
}

func Listen(network, address string) (net.Listener, error) {
	switch network {
	case "sctp", "sctp4", "sctp6":
		sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return sctp.ListenSCTP(network, sctpAddr)
	default:
		return net.Listen(network, address)
	}
}
