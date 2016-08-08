package twistr

import (
	"fmt"
	"net"
)

// This package will evolve into functions around managing multiple games and
// connections.

func Server(port int) (conn net.Conn, err error) {
	var ln net.Listener
	ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	return ln.Accept()
}

func Client(url string) (conn net.Conn, err error) {
	return net.Dial("tcp", url)
}
