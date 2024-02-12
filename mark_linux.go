package proxyguard

import "syscall"

func socketFWMark(fd int, mark int) error {
	return syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, mark)
}
