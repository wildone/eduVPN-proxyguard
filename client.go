package proxyguard

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"syscall"
	"time"
)

// GotClientFD is a function that is called when the Client file descriptor has been obtained
var GotClientFD func(fd int)

// configureSocket creates a TCP dial with fwmark/SO_MARK set
// it also calls the GotClientFD updater
func configureSocket(mark int, sport int) net.Dialer {
	d := net.Dialer{
		Control: func(network, address string, conn syscall.RawConn) error {
			var seterr error
			err := conn.Control(func(fd uintptr) {
				if mark != -1 && runtime.GOOS == "linux" {
					seterr = socketFWMark(int(fd), mark)
				}
				if GotClientFD != nil {
					GotClientFD(int(fd))
				}
			})
			if err != nil {
				return err
			}
			return seterr
		},
		LocalAddr: &net.TCPAddr{
			Port: sport,
		},
	}
	return d
}

// Client runs doClient in a retry loop with a 5 second pause
func Client(ctx context.Context, listen string, tcpsp int, to string, pips []string, fwmark int) error {
	for {
		err := doClient(ctx, listen, tcpsp, to, pips, fwmark)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
		if err != nil {
			log.Logf("Retrying as client exited with error: %v", err)
		} else {
			log.Logf("Retrying as client exited cleanly but context is not canceled yet")
		}
	}
}

// doClient creates a client that forwards UDP to TCP
// listen is the IP:PORT port
// tcpsp is the TCP source port
// to is the IP:PORT string for the TCP proxy on the other end
// fwmark is the mark to set on the TCP socket such that we do not get a routing loop, use -1 to disable setting fwmark
func doClient(ctx context.Context, listen string, tcpsp int, to string, pips []string, fwmark int) (err error) {
	log.Log("Connecting to HTTP server...")
	if tcpsp == -1 {
		laddr, err := net.ResolveTCPAddr("tcp", listen)
		if err != nil {
			return err
		}
		tcpsp = laddr.Port
	}

	u, err := url.Parse(to)
	if err != nil {
		return err
	}

	peerhost := u.Host

	// set fwmark on the socket
	dialer := configureSocket(fwmark, tcpsp)
	c := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
				// no peer ips defined
				// just do the dial context with the configured dialer
				if len(pips) == 0 {
					return dialer.DialContext(ctx, network, addr)
				}

				// there are hardcoded ips given
				// use that instead of a DNS request

				// the address is given in hostname:port
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				// the hostname is not the peer hostname
				// return the default dialcontext
				if host != peerhost {
					return dialer.DialContext(ctx, network, addr)
				}

				// otherwise loop over the ips and return if one succeeds
				for _, ip := range pips {
					conn, err = dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
					if err == nil {
						return conn, nil
					}
				}
				return conn, err
			},
		},
	}

	req, err := http.NewRequestWithContext(ctx, "GET", to, nil)
	if err != nil {
		return err
	}

	// upgrade the connection to wireguard
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", UpgradeProto)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	// TODO: why does nhooyr.io/websocket set the body to nil and make a rb copy?
	// is this needed?
	rb := resp.Body
	resp.Body = nil

	// TODO: clean this up?
	cancel := make(chan struct{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				rb.Close()
				return
			case <-cancel:
				rb.Close()
				return
			}
		}
	}()
	defer close(cancel)

	if resp.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("status is not switching protocols, got: '%v'", resp.StatusCode)
	}

	if resp.Header.Get("Connection") != "Upgrade" {
		return fmt.Errorf("the 'Connection' header is not 'Upgrade', got: '%v'", resp.Header.Get("Connection"))
	}

	if resp.Header.Get("Upgrade") != UpgradeProto {
		return fmt.Errorf("upgrade header is not '%v', got: '%v'", UpgradeProto, resp.Header.Get("Upgrade"))
	}

	rwc, ok := rb.(io.ReadWriteCloser)
	if !ok {
		return fmt.Errorf("response body is not of type io.ReadWriteCloser: %T", rb)
	}
	log.Log("Connected to HTTP server")

	udpaddr, err := net.ResolveUDPAddr("udp", listen)
	if err != nil {
		return err
	}
	log.Log("Waiting for first UDP packet...")
	wgaddr, first, err := inferUDPAddr(ctx, udpaddr)
	if err != nil {
		return err
	}
	log.Logf("First UDP packet received with address: %s", wgaddr.String())
	wgconn, err := net.DialUDP("udp", udpaddr, wgaddr)
	if err != nil {
		return err
	}
	defer wgconn.Close()
	log.Log("Client is ready for converting UDP<->HTTP")

	// create a buffered read writer
	rw := bufio.NewReadWriter(bufio.NewReader(rwc), bufio.NewWriter(rwc))

	// first forward the outstanding packet
	err = writeTCP(rw.Writer, first, len(first)-hdrLength)
	if err != nil {
		log.Logf("Failed forwarding first outstanding packet: %v", err)
	}

	return tunnel(ctx, wgconn, rw)
}
