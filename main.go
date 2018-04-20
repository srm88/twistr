package main

import "github.com/srm88/twistr/twistr"
import "log"
import "os"
import "os/signal"
import "path/filepath"
import "regexp"

func isServer(ui twistr.UI) bool {
	var reply string
	twistr.Input(ui, &reply, "Are you the host or the guest?", "host", "guest")
	return reply == "host"
}

func choosePlayer(ui twistr.UI) (player twistr.Aff) {
	twistr.Input(ui, &player, "Who are you playing as?", "usa", "ussr")
	return
}

var ValidName = regexp.MustCompile(`^[a-z0-9-$_.]+$`)

func chooseName(ui twistr.UI) string {
	var reply string
	twistr.Input(ui, &reply, "Choose a name for this game")
	for !ValidName.MatchString(reply) {
		twistr.Input(ui, &reply, "Choose a name for this game (a-z0-9-_.$ chars allowed)")
	}
	return reply
}

type Match interface {
	Run() error
	Close()
}

// Temp:
func main() {
	logFile, err := os.OpenFile(filepath.Join(twistr.DataDir, "twistr.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	ui := twistr.MakeNCursesUI()

	var match Match

	if isServer(ui) {
		log.SetPrefix("host  ")
		// Need to tell the opponent who they are!
		// Where should this happen? Should address when we formalize loading previous games
		match = twistr.NewHostMatch(ui, chooseName(ui), choosePlayer(ui))
	} else {
		log.SetPrefix("guest ")
		match = twistr.NewGuestMatch(ui)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	done := make(chan int, 1)
	// This goroutine cleans up
	go func() {
		exitcode := <-done
		ui.Close()
		match.Close()
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
	match.Run()
}
