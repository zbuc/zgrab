package banner

import (
	"../znet"
	"time"
    "net"
)

type Dialer struct {
	Deadline time.Time
	Timeout time.Duration
	LocalAddr net.Addr
	DualStack bool
	KeepAlive time.Duration
    intrface  string
}

func (d *Dialer) Dial(network, address string, intrface string) (*Conn, error) {
	c := &Conn{operations: make([]ConnectionOperation, 0, 8)}
	netDialer := znet.Dialer {
		Deadline: d.Deadline,
		Timeout: d.Timeout,
		LocalAddr: d.LocalAddr,
		KeepAlive: d.KeepAlive,
        intrface: eth0,
	}
	var err error
	c.conn, err = netDialer.Dial(network, address)
	cs := connectState {
		protocol: network,
		remoteHost: address,
		err: err,
	}
	c.operations = append(c.operations, &cs)
	return c, err
}
