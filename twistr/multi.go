// This package will evolve into functions around managing multiple games and
// connections.
package twistr

import "bufio"
import "bytes"
import "fmt"
import "io"
import "io/ioutil"
import "log"
import "net"
import "os"
import "os/user"
import "path/filepath"
import "strings"

// Startup:
// server syncs existing aof to client

// client:
// read from (synced) AOF if data remains
// remote player? read from conn
// else get input, buffer
// flush/commit to conn

// server:
// write AOF to conn on startup (sync)
// remote player? read from conn
// else get input, buffer
// flush/commit to AOF+conn

var (
	DataDir string
)

func init() {
	u, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	DataDir = filepath.Join(u.HomeDir, ".twistr")
	if err := os.MkdirAll(DataDir, os.ModePerm); err != nil {
		panic(err)
	}
}

func Server(port int) (conn net.Conn, err error) {
	var ln net.Listener
	ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Error listening on %d: %s", port, err.Error())
		return
	}
	conn, err = ln.Accept()
	if err != nil {
		log.Printf("Error accepting conn: %s", err.Error())
	}
	return
}

func Client(url string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", url)
	if err != nil {
		log.Printf("Error connecting to server: %s", err.Error())
	}
	return
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
			return nil, fmt.Errorf("Error copying aof bytes to buffer: %s", err.Error())
		}
		if err := in.Close(); err != nil {
			return nil, fmt.Errorf("Error closing aof: %s", err.Error())
		}
		return b.Bytes(), nil
	}
}

type Match struct {
	UI      UI
	Port    int
	Name    string
	Game    *Game
	State   *State
	closers []io.Closer
	Who     Aff
	Conn    net.Conn
	// Connected? Synced?
	// Peer address / ports?
}

func NewMatch(ui UI) *Match {
	return &Match{
		UI:      ui,
		Port:    1550,
		Game:    NewGame(),
		closers: []io.Closer{}}
}

func (m *Match) AofPath() string {
	return fmt.Sprintf("%s.aof", filepath.Join(DataDir, m.Name))
}

func (m *Match) Start() error {
	log.Printf("Starting")
	Start(m.State)
	return nil
}

func (m *Match) Close() {
	for _, c := range m.closers {
		c.Close()
	}
}

type HostMatch struct {
	*Match
}

func NewHostMatch(ui UI, name string, who Aff) *HostMatch {
	m := NewMatch(ui)
	m.Name = name
	m.Who = who
	return &HostMatch{
		Match: m,
	}
}

func (h *HostMatch) Run() (err error) {
	return h.Connect()
}

func (h *HostMatch) Connect() (err error) {
	log.Println("Connecting")
	h.Conn, err = Server(h.Port)
	if err == nil {
		h.closers = append(h.closers, h.Conn)
		log.Println("Connected")
	} else {
		log.Printf("Failed to connect to guest: %s\n", err.Error())
		return
	}
	return h.Sync()
}

func (h *HostMatch) Sync() (err error) {
	if _, err = h.Conn.Write([]byte("$ BEGIN AOF\n")); err != nil {
		log.Printf("Failed to send aof header: %s\n", err.Error())
		return
	}
	var in io.Reader
	// Don't we need to close this file?
	in, err = os.Open(h.AofPath())
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Failed to open aof to sync ... %s\n", err.Error())
			return
		}
	} else {
		log.Println("Server syncing aof")
		if _, err = io.Copy(h.Conn, in); err != nil {
			log.Printf("Failed while sending sync ... %s\n", err.Error())
			return
		}
	}
	if _, err = h.Conn.Write([]byte("$ END AOF\n")); err != nil {
		log.Printf("Failed to send aof footer: %s\n", err.Error())
		return
	}
	log.Printf("Server synced aof")
	return h.Setup()
}

func (h *HostMatch) Setup() error {
	// In
	b, err := loadAof(h.AofPath())
	if err != nil {
		return err
	}
	var history *History
	if len(b) > 0 {
		history = NewHistoryBacklog(h.UI, strings.Split(string(b), "\n"))
	} else {
		history = NewHistory(h.UI)
	}
	// Out
	out, err := os.OpenFile(h.AofPath(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	h.closers = append(h.closers, out)
	h.State = NewState(history, h.Game, true, h.Who, out)
	h.State.LinkIn = NewCmdIn(bufio.NewScanner(h.Conn))
	h.State.LinkOut = NewCmdOut(h.Conn)
	return h.Match.Start()
}

type GuestMatch struct {
	*Match
	HostFeed *bufio.Scanner
	Aof      []string
}

func NewGuestMatch(ui UI) *GuestMatch {
	return &GuestMatch{
		Match: NewMatch(ui),
		Aof:   []string{},
	}
}

func (g *GuestMatch) ConnectHost() string {
	return "localhost"
}

func (g *GuestMatch) Run() (err error) {
	return g.Connect()
}

func (g *GuestMatch) Connect() (err error) {
	log.Println("Connecting")
	g.Conn, err = Client(fmt.Sprintf("%s:%d", g.ConnectHost(), g.Port))
	if err == nil {
		g.closers = append(g.closers, g.Conn)
		log.Println("Connected")
	} else {
		log.Printf("Failed to connect to host: %s\n", err.Error())
		return
	}
	return g.Sync()
}

func (g *GuestMatch) Sync() error {
	// XXX: this is not re-entrant
	log.Println("Client receiving aof sync")
	g.HostFeed = bufio.NewScanner(g.Conn)
	var line string
ReadAof:
	for g.HostFeed.Scan() {
		line = g.HostFeed.Text()
		switch {
		case line == "$ BEGIN AOF":
		case line == "$ END AOF":
			break ReadAof
		default:
			g.Aof = append(g.Aof, line)
		}
	}
	if err := g.HostFeed.Err(); err != nil {
		log.Printf("Failed while reading sync ... %s\n", err.Error())
		return err
	}
	log.Println("Client received aof")
	return g.Setup()
}

func (g *GuestMatch) Setup() error {
	var history *History
	if len(g.Aof) > 0 {
		history = NewHistoryBacklog(g.UI, g.Aof)
	} else {
		history = NewHistory(g.UI)
	}
	g.State = NewState(history, g.Game, false, g.Who, ioutil.Discard)
	g.State.LinkIn = NewCmdIn(g.HostFeed)
	g.State.LinkOut = NewCmdOut(g.Conn)
	return g.Match.Start()
}
