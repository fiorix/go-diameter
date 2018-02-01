package diam

import (
	"github.com/ishidawataru/sctp"
	"net"
	"time"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type sctpDialer struct{}

func (_ sctpDialer) Dial(network, address string) (net.Conn, error) {
	sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
	if err != nil {
		return nil, err
	}
	return sctp.DialSCTP(network, nil, sctpAddr)
}

func getDialer(network string, timeout time.Duration) Dialer {
	switch network {
	case "sctp", "sctp4", "sctp6":
		return sctpDialer{}
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
