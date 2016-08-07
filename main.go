package main

import (
	"fmt"
	"github.com/srm88/twistr/twistr"
	"os"
	"os/signal"
)

const (
	aofPath = "/tmp/twistr.aof"
)

func setup() (*twistr.State, error) {
	ui := twistr.MakeNCursesUI()
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
	s.UI = ui
	s.Aof = aof
	s.Txn = txn
	return s, nil
}

// Temp:
func main() {
	state, err := setup()
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
