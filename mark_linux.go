package proxyguard

import (
	"net"
	"syscall"
)

// markedDial creates a TCP dial with fwmark/SO_MARK set
func markedDial(mark int, laddr *net.TCPAddr, to string) (net.Conn, error) {
	d := net.Dialer{
		Control: func(network, address string, conn syscall.RawConn) error {
			var seterr error
			err := conn.Control(func(fd uintptr) {
				seterr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, mark)
			})
			if err != nil {
				return err
			}
			return seterr
		},
		LocalAddr: &net.TCPAddr{
			Port: laddr.Port,
		},
	}
	return d.Dial("tcp", to)
}
