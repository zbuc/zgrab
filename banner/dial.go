package banner

import (
	"../znet"
	"time"
)

type Dialer struct {
	Deadline time.Time
	Timeout time.Duration
	LocalAddr znet.Addr
	DualStack bool
	KeepAlive time.Duration
    	Intrface  string
}

func (d *Dialer) Dial(network, address string, intface string) (*Conn, error) {
	c := &Conn{operations: make([]ConnectionOperation, 0, 8)}
	netDialer := znet.Dialer {
		Deadline: d.Deadline,
		Timeout: d.Timeout,
		LocalAddr: d.LocalAddr,
		KeepAlive: d.KeepAlive,
        	Intrface: "eth0",
	}
	var err error
	c.conn, err = netDialer.Dial(network, address, "eth0")
	cs := connectState {
		protocol: network,
		remoteHost: address,
		err: err,
	}
	c.operations = append(c.operations, &cs)
	return c, err
}
