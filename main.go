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
	port int
)

func init() {
	flag.IntVar(&port, "port", 1551, "Server port number")
}

// XXX: should also include game name in the file path
func aofPath(server bool) string {
	var fname string
	switch {
	case server:
		fname = "server.aof"
	default:
		fname = "client.aof"
	}
	return fmt.Sprintf("%s/%s", AofDir, fname)
}

func setup(aofPath string) (*twistr.State, error) {
	in, err := os.Open(aofPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		in, err = os.OpenFile(aofPath, os.O_CREATE|os.O_RDONLY, 0666)
		if err != nil {
			return nil, err
		}
	}
	txn, err := twistr.OpenTxnLog(aofPath)
	if err != nil {
		return nil, err
	}
	aof := twistr.NewAof(in, txn)
	s := twistr.NewState()
	s.Aof = aof
	s.Txn = txn
	return s, nil
}

func connectHost(nc *twistr.NCursesUI) string {
	return "localhost"
}

func isServer(nc *twistr.NCursesUI) bool {
	var reply string
	for {
		reply = nc.GetInput("server or client")
		reply = strings.ToLower(reply)
		switch reply {
		case "server":
			return true
		case "client":
			return false
		}
	}

}
func connect(nc *twistr.NCursesUI, server bool) (net.Conn, error) {
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
	toClose := []io.Closer{ui}

	// XXX: revisit 'aof path' and 'is server' for multiple game support
	server := isServer(ui)
	path := aofPath(server)

	state, err := setup(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to start game: %s\n", err.Error()))
	}
	state.UI = ui
	toClose = append(toClose, state)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	done := make(chan int, 1)
	// This goroutine cleans up
	go func() {
		exitcode := <-done
		for _, thing := range toClose {
			thing.Close()
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

	state.Conn, err = connect(ui, server)
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect %s\n", err.Error()))
	}
	toClose = append(toClose, state.Conn)

	twistr.Start(state)
}
