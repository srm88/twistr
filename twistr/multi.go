package twistr

import (
	"fmt"
	"io"
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

type multiWriteCloser struct {
	io.Writer
	wcs []io.WriteCloser
}

func (m *multiWriteCloser) Close() error {
	for _, wc := range m.wcs {
		if err := wc.Close(); err != nil {
			return err
		}
	}
	return nil
}

func MultiWriteCloser(wcs ...io.WriteCloser) io.WriteCloser {
	writers
	writeClosers := make([]io.WriteCloser, len(wcs))
	copy(writeClosers, wcs)
	return &multiWriteCloser{
		Writer: io.MultiWriter(wcs),
		wcs:    writeClosers,
	}
}

type Link struct {
	conn net.Conn
}
