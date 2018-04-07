package twistr

import "bytes"
import "fmt"
import "io"
import "io/ioutil"
import "log"
import "net"
import "os"
import "os/user"
import "path/filepath"

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

// This package will evolve into functions around managing multiple games and
// connections.

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

func isServer(ui UI) bool {
	var reply string
	input(ui, &reply, "Are you the host or the guest?", "host", "guest")
	return reply == "host"
}

func choosePlayer(ui UI) (player Aff) {
	input(ui, &player, "Who are you playing as?", "usa", "ussr")
	return
}

func NameGame() string {
	return "server"
}

func connectHost() string {
	return "localhost"
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
	ServerMode bool
	Name       string
	Game       *Game
	State      *State
	closers    []io.Closer
	Who        Aff
	// Connected? Synced?
	// Peer address / ports?
}

func NewMatch(ui UI) *Match {
	return &Match{
		UI:       ui,
		Port:     1550,
		SyncPort: 1551,
		Name:     NameGame(),
		Game:     NewGame(),
		closers:  []io.Closer{}}
}

func (m *Match) AofPath() string {
	return fmt.Sprintf("%s.aof", filepath.Join(DataDir, m.Name))
}

func (m *Match) setupServer() error {
	// In
	b, err := loadAof(m.AofPath())
	if err != nil {
		return err
	}
	var history *History
	if len(b) > 0 {
		history = NewHistoryBacklog(m.UI, string(b))
	} else {
		history = NewHistory(m.UI)
	}
	// Out
	out, err := os.OpenFile(m.AofPath(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	m.State = NewState(history, m.Game, true, m.Who, out)
	m.closers = append(m.closers, out)
	return nil
}

func (m *Match) setupClient(aof string) {
	var history *History
	if len(aof) > 0 {
		history = NewHistoryBacklog(m.UI, aof)
	} else {
		history = NewHistory(m.UI)
	}
	m.State = NewState(history, m.Game, false, m.Who, ioutil.Discard)
}

func (m *Match) connect() (net.Conn, error) {
	switch {
	case m.ServerMode:
		return Server(m.Port)
	default:
		return Client(fmt.Sprintf("%s:%d", connectHost(), m.Port))
	}
}

func (m *Match) syncConnect() (net.Conn, error) {
	switch {
	case m.ServerMode:
		return Server(m.SyncPort)
	default:
		return Client(fmt.Sprintf("%s:%d", connectHost(), m.SyncPort))
	}
}

func (m Match) sendAof() (err error) {
	// fails if aof doesn't exist
	var in io.Reader
	in, err = os.Open(m.AofPath())
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Failed to open aof to sync ... %s\n", err.Error())
			return
		}
		in = new(bytes.Buffer) // what is this
	}
	log.Println("Server connecting to sync aof")
	syncConn, err := m.syncConnect()
	if err != nil {
		log.Printf("Failed to accept client conn to sync ... %s\n", err.Error())
		return
	}
	defer syncConn.Close()
	log.Println("Server syncing aof")
	if _, err = io.Copy(syncConn, in); err != nil {
		log.Printf("Failed while sending sync ... %s\n", err.Error())
		return
	}
	log.Printf("Server sent aof")
	return
}

func (m *Match) receiveAof() (string, error) {
	b := new(bytes.Buffer)
	log.Println("Client connecting to sync aof")
	syncConn, err := m.syncConnect()
	if err != nil {
		log.Printf("Failed to dial server to sync ... %s\n", err.Error())
		return "", err
	}
	defer syncConn.Close()
	log.Println("Client receiving aof sync")
	if _, err := io.Copy(b, syncConn); err != nil {
		log.Printf("Failed while reading sync ... %s\n", err.Error())
		return "", err
	}
	log.Printf("Client received aof")
	return b.String(), nil
}

func (m *Match) Run() (err error) {
	m.ServerMode = isServer(m.UI)
	// Need to tell the opponent who they are!
	m.Who = choosePlayer(m.UI)
	if m.ServerMode {
		if err = m.sendAof(); err != nil {
			return
		}
		if err = m.setupServer(); err != nil {
			return
		}
	} else {
		aof, err := m.receiveAof()
		if err != nil {
			return err
		}
		m.setupClient(aof)
	}
	var conn net.Conn
	log.Printf("Connecting")
	if conn, err = m.connect(); err != nil {
		return
	}
	log.Printf("Connected")
	m.closers = append(m.closers, conn)
	m.State.LinkIn = NewCmdIn(conn)
	m.State.LinkOut = NewCmdOut(conn)
	log.Printf("Starting")
	Start(m.State)
	return nil
}

func (m *Match) Close() {
	for _, c := range m.closers {
		c.Close()
	}
}
