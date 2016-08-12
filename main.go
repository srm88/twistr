package main

import (
	"fmt"
	"github.com/srm88/twistr/twistr"
	"os"
	"os/signal"
)

// Temp:
func main() {
	game := twistr.NewGame()
	ui := twistr.MakeNCursesUI()
	state, err := twistr.NewState(ui, "/tmp/twistr.aof", game)
	if err != nil {
		panic(fmt.Sprintf("Failed to start game: %s\n", err.Error()))
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	done := make(chan int, 1)
	// This goroutine cleans up
	go func() {
		exitcode := <-done
		state.Close()
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
	twistr.Start(state)
}
