package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/srm88/twistr/twistr"
	"io"
	"io/ioutil"
	"log"
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

func loadAof(aofPath string) ([]byte, error) {
	in, err := os.Open(aofPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("Error reading existing AOF: %s", err.Error())
		}
		return nil, nil
	} else {
		b := new(bytes.Buffer)
		if _, err := io.Copy(b, in); err != nil {
			return nil, err
		}
		if err := in.Close(); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}
}

func setupMaster(ui twistr.UI, game *twistr.Game, aofPath string) error {
	// In
	var history *twistr.History
	b, err := loadAof(aofPath)
	if err != nil {
		return err
	}
	if len(b) > 0 {
		history = twistr.NewHistoryBacklog(ui, string(b))
	} else {
		history = twistr.NewHistory(ui)
	}
	// Out
	out, err := os.OpenFile(aofPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	state = twistr.NewState(history, game, true, twistr.SOV, out)
	closers = append(closers, out)
	return nil
}

// Return writer with which to accept sync'd aof
func setupSlave(ui twistr.UI, game *twistr.Game, replay string) {
	var history *twistr.History
	if len(replay) > 0 {
		history = twistr.NewHistoryBacklog(ui, replay)
	} else {
		history = twistr.NewHistory(ui)
	}
	state = twistr.NewState(history, game, false, twistr.USA, ioutil.Discard)
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

func syncAof(aofPath string) {
	if !serverMode {
		return
	}
	// fails if aof doesn't exist
	var in io.Reader
	var err error
	in, err = os.Open(aofPath)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Sprintf("Failed to open aof to sync ... %s\n", err.Error()))
		}
		in = new(bytes.Buffer)
	}
	secretConn, err := secretConnect()
	if err != nil {
		panic(fmt.Sprintf("Failed to accept client conn to sync ... %s\n", err.Error()))
	}
	log.Println("Master syncing aof")
	if _, err := io.Copy(secretConn, in); err != nil {
		panic(fmt.Sprintf("Failed while sending sync ... %s\n", err.Error()))
	}
	secretConn.Close()
}

func receiveAof() string {
	if serverMode {
		return ""
	}
	clientReplay := new(bytes.Buffer)
	secretConn, err := secretConnect()
	if err != nil {
		panic(fmt.Sprintf("Failed to dial master to sync ... %s\n", err.Error()))
	}
	log.Println("Client receiving aof sync")
	if _, err := io.Copy(clientReplay, secretConn); err != nil {
		panic(fmt.Sprintf("Failed while reading sync ... %s\n", err.Error()))
	}
	secretConn.Close()
	return clientReplay.String()
}

// Temp:
func main() {
	flag.Parse()

	game := twistr.NewGame()
	ui := twistr.MakeNCursesUI()

	closers = append(closers, ui)
	// XXX: revisit 'aof path' and 'is server' for multiple game support
	serverMode = isServer(ui)
	path := aofPath()

	var err error
	if serverMode {
		syncAof(path)
		if err = setupMaster(ui, game, path); err != nil {
			panic(fmt.Sprintf("Failed to start game: %s\n", err.Error()))
		}
	} else {
		syncedAof := receiveAof()
		setupSlave(ui, game, syncedAof)
	}

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

	twistr.Start(state)
}
