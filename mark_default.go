//go:build !linux

package proxyguard

func socketReuseSport(fd int) error {
	panic("reusing a source port is not supported on this OS")
}

func socketFWMark(fd int, mark int) error {
	panic("setting fwmark is not supported on this OS")
}
