package main

import (
	"flag"
	"fmt"
	"github.com/srm88/twistr/twistr"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
)

const (
	AofDir = "/tmp/twistr"
)

var (
	port    int
	server  bool
	state   *twistr.State
	closers []io.Closer
)

func init() {
	flag.IntVar(&port, "port", 1551, "Server port number")
	closers = []io.Closer{}
}

// XXX: should also include game name in the file path
func aofPath() string {
	var fname string
	switch {
	case server:
		fname = "server.aof"
	default:
		fname = "client.aof"
	}
	return fmt.Sprintf("%s/%s", AofDir, fname)
}

func setup(aofPath string) error {
	in, err := os.Open(aofPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		in, err = os.OpenFile(aofPath, os.O_CREATE|os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
	}
	aof := twistr.NewCommandStream(in)
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		in.Close()
		return err
	}
	state = twistr.NewState()
	state.Aof = aof
	state.Txn = twistr.NewTxnLog(out)
	closers = append(closers, in, out)
	return nil
}

func connectHost(nc *twistr.NCursesUI) string {
	return "localhost"
}

func isServer(nc *twistr.NCursesUI) bool {
	var reply string
	for {
		reply = twistr.Solicit(nc, "What mode?", []string{"server", "client"})
		reply = strings.ToLower(reply)
		switch reply {
		case "server":
			return true
		case "client":
			return false
		}
	}

}
func connect(nc *twistr.NCursesUI) (net.Conn, error) {
	switch {
	case server:
		return twistr.Server(port)
	default:
		return twistr.Client(fmt.Sprintf("%s:%d", connectHost(nc), port))
	}
}

// Temp:
func main() {
	flag.Parse()

	ui := twistr.MakeNCursesUI()
	closers = append(closers, ui)

	// XXX: revisit 'aof path' and 'is server' for multiple game support
	server = isServer(ui)
	path := aofPath(server)

	if err := setup(path); err != nil {
		panic(fmt.Sprintf("Failed to start game: %s\n", err.Error()))
	}
	state.UI = ui

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	done := make(chan int, 1)
	// This goroutine cleans up
	go func() {
		exitcode := <-done
		for _, c := range closers {
			c.Close()
		}
		os.Exit(exitcode)
	}()
	// This one forwards the signal to the cleanup routine
	go func() {
		<-sigs
		done <- 1
	}()
	// This one forwards the signal during a normal exit
	defer func() {
		done <- 0
	}()

	if conn, err := connect(ui, server); err != nil {
		panic(fmt.Sprintf("Couldn't connect %s\n", err.Error()))
	}
	closers = append(closers, conn)
	state.Link = twistr.NewCommandStream(conn)

	twistr.Start(state)
}
