package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/srm88/twistr/twistr"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strings"
)

const (
	AofDir = "/tmp/twistr"
)

var (
	port       int
	secretPort int
	serverMode bool
	state      *twistr.State
	closers    []io.Closer
)

func init() {
	flag.IntVar(&port, "port", 1551, "Server port number")
	secretPort = port + 1
	closers = []io.Closer{}
}

// XXX: should include game name in the file path
func aofPath() string {
	return fmt.Sprintf("%s/%s", AofDir, "server.aof")
}

func setupMaster(aofPath string) error {
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
	out, err := os.OpenFile(aofPath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		in.Close()
		return err
	}
	state = twistr.NewState()
	state.Master = true
	state.Replay = twistr.NewCmdIn(in)
	state.Aof = twistr.NewCmdOut(out)
	// XXX hardcode
	state.LocalPlayer = twistr.SOV
	closers = append(closers, in, out)
	return nil
}

// Return writer with which to accept sync'd aof
func setupSlave() io.Writer {
	replayBuf := new(bytes.Buffer)
	state = twistr.NewState()
	state.Master = false
	state.Replay = twistr.NewCmdIn(replayBuf)
	state.Aof = twistr.NewCmdOut(ioutil.Discard)
	// XXX hardcode
	state.LocalPlayer = twistr.USA
	return replayBuf
}

func connectHost() string {
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

func connect() (net.Conn, error) {
	switch {
	case serverMode:
		return twistr.Server(port)
	default:
		return twistr.Client(fmt.Sprintf("%s:%d", connectHost(), port))
	}
}

func secretConnect() (net.Conn, error) {
	switch {
	case serverMode:
		return twistr.Server(secretPort)
	default:
		return twistr.Client(fmt.Sprintf("%s:%d", connectHost(), secretPort))
	}
}

// Temp:
func main() {
	flag.Parse()

	ui := twistr.MakeNCursesUI()
	closers = append(closers, ui)

	// XXX: revisit 'aof path' and 'is server' for multiple game support
	serverMode = isServer(ui)
	path := aofPath()

	var err error
	var clientReplay io.Writer
	if serverMode {
		if err = setupMaster(path); err != nil {
			panic(fmt.Sprintf("Failed to start game: %s\n", err.Error()))
		}
	} else {
		clientReplay = setupSlave()
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

	var conn net.Conn
	if conn, err = connect(); err != nil {
		panic(fmt.Sprintf("Couldn't connect %s\n", err.Error()))
	}
	closers = append(closers, conn)
	state.LinkIn = twistr.NewCmdIn(conn)
	state.LinkOut = twistr.NewCmdOut(conn)

	// XXX this all needs to move out of main ...
	// sync aof to slave
	if serverMode {
		in, err := os.Open(path)
		if err != nil {
			panic(fmt.Sprintf("Failed to open aof to sync ... %s\n", err.Error()))
		}
		secretConn, err := secretConnect()
		if err != nil {
			panic(fmt.Sprintf("Failed to accept client conn to sync ... %s\n", err.Error()))
		}
		twistr.Debug("Master syncing aof")
		if _, err := io.Copy(secretConn, in); err != nil {
			panic(fmt.Sprintf("Failed while sending sync ... %s\n", err.Error()))
		}
		secretConn.Close()
	} else {
		secretConn, err := secretConnect()
		if err != nil {
			panic(fmt.Sprintf("Failed to dial master to sync ... %s\n", err.Error()))
		}
		twistr.Debug("Client receiving aof sync")
		if _, err := io.Copy(clientReplay, secretConn); err != nil {
			panic(fmt.Sprintf("Failed while reading sync ... %s\n", err.Error()))
		}
		secretConn.Close()
	}
	twistr.Start(state)
}
