//go:build !linux

package proxyguard

func socketFWMark(fd int, mark int) error {
	panic("setting fwmark is not supported on this OS")
}
