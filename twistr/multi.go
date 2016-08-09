package twistr

import (
	"fmt"
	"net"
)

// This package will evolve into functions around managing multiple games and
// connections.

// Startup:
// master syncs existing aof to slave

// Slave:
// read from (synced) AOF if data remains
// remote player? read from conn
// else get input, buffer
// flush/commit to conn

// Master:
// write AOF to conn on startup (sync)
// remote player? read from conn
// else get input, buffer
// flush/commit to AOF+conn

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
