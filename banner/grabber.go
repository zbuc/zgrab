package banner

import (
	"crypto/x509"
	"log"
	"net"
	"strconv"
	"time"

	"../zcrypto/ztls"
)

type GrabConfig struct {
	Tls          bool
	TlsVersion   uint16
	Banners      bool
	SendMessage  bool
	ReadResponse bool
	Smtp         bool
	Ehlo         bool
	SmtpHelp     bool
	StartTls     bool
	Imap         bool
	Pop3         bool
	Heartbleed   bool
	Port         uint16
	Timeout      time.Duration
	Message      []byte
	EhloDomain   string
	Protocol     string
	ErrorLog     *log.Logger
	LocalAddr    net.Addr
	RootCAPool   *x509.CertPool
	CbcOnly      bool
}

type Grab struct {
	Host string     `json:"host"`
	Port uint16     `json:"port"`
	Time time.Time  `json:"timestamp"`
	Log  []StateLog `json:"log"`
}

type Progress struct {
	Success uint
	Error   uint
	Total   uint
}

func makeDialer(c *GrabConfig) func(string) (*Conn, error) {
	proto := c.Protocol
	timeout := c.Timeout
	return func(addr string) (*Conn, error) {
		deadline := time.Now().Add(timeout)
		d := Dialer{
			Deadline: deadline,
		}
		conn, err := d.Dial(proto, addr)
		conn.maxTlsVersion = c.TlsVersion
		if err == nil {
			conn.SetDeadline(deadline)
		}
		return conn, err
	}
}

func makeGrabber(config *GrabConfig) func(*Conn) ([]StateLog, error) {
	// Do all the hard work here
	g := func(c *Conn) error {
		//banner := make([]byte, 1024)
		//response := make([]byte, 65536)
		c.SetCAPool(config.RootCAPool)
		c.SetCbcOnly() // XXX THIS IS TERRIBLE
		c.maxTlsVersion = ztls.VersionSSL30
		if err := c.TlsHandshake(); err == nil {
			c.Close()
		}
		if err := c.ReDial(); err != nil {
			log.Print("whoops")
			log.Print(err)
			return err
		}
		c.maxTlsVersion = ztls.VersionTLS12
		if err := c.TlsHandshake(); err != nil {
			return err
		}
		return nil
	}
	// Wrap the whole thing in a logger
	return func(c *Conn) ([]StateLog, error) {
		rh := c.RemoteAddr()
		err := g(c)
		if err != nil {
			config.ErrorLog.Printf("Conversation error with remote host %s: %s",
				rh, err.Error())
		}
		return c.States(), err
	}
}

func GrabBanner(addrChan chan net.IP, grabChan chan Grab, doneChan chan Progress, config *GrabConfig) {
	dial := makeDialer(config)
	grabber := makeGrabber(config)
	port := strconv.FormatUint(uint64(config.Port), 10)
	p := Progress{}
	for ip := range addrChan {
		p.Total += 1
		addr := ip.String()
		rhost := net.JoinHostPort(addr, port)
		t := time.Now()
		conn, dialErr := dial(rhost)
		if dialErr != nil {
			// Could not connect to host
			config.ErrorLog.Printf("Could not connect to remote host %s: %s",
				addr, dialErr.Error())
			grabChan <- Grab{addr, config.Port, t, conn.States()}
			p.Error += 1
			continue
		}
		grabStates, err := grabber(conn)
		if err != nil {
			p.Error += 1
		} else {
			p.Success += 1
		}
		grabChan <- Grab{addr, config.Port, t, grabStates}
	}
	doneChan <- p
}
