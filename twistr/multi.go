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
import "regexp"
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
	DataDir   string
	ValidName = regexp.MustCompile(`^[a-z0-9-$_.]+$`)
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

func choosePlayer(ui UI) (player Aff) {
	Input(ui, &player, "Who are you playing as?", "usa", "ussr")
	return
}

func chooseName(ui UI) string {
	var reply string
	Input(ui, &reply, "Choose a name for this game")
	for !ValidName.MatchString(reply) {
		Input(ui, &reply, "Choose a name for this game (a-z0-9-_.$ chars allowed)")
	}
	return reply
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
	UI         UI
	Port       int
	SyncPort   int
	Name       string
	Game       *Game
	State      *State
	closers    []io.Closer
	Who        Aff
	Conn	net.Conn
	// Connected? Synced?
	// Peer address / ports?
}

func NewMatch(ui UI) *Match {
	return &Match{
		UI:       ui,
		Port:     1550,
		SyncPort: 1551,
		Game:     NewGame(),
		closers:  []io.Closer{}}
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

func NewHostMatch(ui UI) *HostMatch {
	return &HostMatch{
		Match: NewMatch(ui),
	}
}

func (h *HostMatch) Run() (err error) {
	return h.GameSelect()
}

func (h *HostMatch) GameSelect() error {
	h.Name = chooseName(h.UI)
	// Need to tell the opponent who they are!
	h.Who = choosePlayer(h.UI)
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
	return h.SendAof()
}

func (h *HostMatch) SendAof() (err error) {
	var in io.Reader
	in, err = os.Open(h.AofPath())
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Failed to open aof to sync ... %s\n", err.Error())
			return
		}
		in = new(bytes.Buffer)
	}
	log.Println("Server connecting to sync aof")
	syncConn, err := Server(h.SyncPort)
	if err != nil {
		log.Printf("Failed to accept client conn to sync ... %s\n", err.Error())
		return
	}
	defer syncConn.Close()
	log.Println("Server syncing aof")
	if _, err = syncConn.Write([]byte("$ BEGIN AOF\n")); err != nil {
		log.Printf("Failed to send aof header: %s\n", err.Error())
		return
	}
	if _, err = io.Copy(syncConn, in); err != nil {
		log.Printf("Failed while sending sync ... %s\n", err.Error())
		return
	}
	if _, err = syncConn.Write([]byte("$ END AOF\n")); err != nil {
		log.Printf("Failed to send aof footer: %s\n", err.Error())
		return
	}
	log.Printf("Server sent aof")
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
	h.State.LinkIn = NewCmdIn(h.Conn)
	h.State.LinkOut = NewCmdOut(h.Conn)
	return h.Match.Start()
}

type GuestMatch struct {
	*Match
	Aof	[]string
}

func NewGuestMatch(ui UI) *GuestMatch {
	return &GuestMatch{
		Match: NewMatch(ui),
		Aof: []string{},
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
	return g.ReceiveAof()
}

func (g *GuestMatch) ReceiveAof() error {
	// XXX: this is not re-entrant
	log.Println("Client connecting to sync aof")
	syncConn, err := Client(fmt.Sprintf("%s:%d", g.ConnectHost(), g.SyncPort))
	if err != nil {
		log.Printf("Failed to dial server to sync ... %s\n", err.Error())
		return err
	}
	defer syncConn.Close()
	log.Println("Client receiving aof sync")
	scanner := bufio.NewScanner(syncConn)
	var line string
ReadAof:
	for scanner.Scan() {
		line = scanner.Text()
		switch {
		case line == "$ BEGIN AOF":
		case line == "$ END AOF":
			break ReadAof
		default:
			g.Aof = append(g.Aof, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Failed while reading sync ... %s\n", err.Error())
		return err
	}
	log.Printf("Client received aof")
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
	g.State.LinkIn = NewCmdIn(g.Conn)
	g.State.LinkOut = NewCmdOut(g.Conn)
	return g.Match.Start()
}
