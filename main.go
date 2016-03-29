package main

import (
	"fmt"
	"github.com/srm88/twistr/twistr"
)

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	state, err := twistr.NewState(ui, "/tmp/twistr.aof")
	if err != nil {
		panic(fmt.Sprintf("Failed to start game: %s\n", err.Error()))
	}
	defer state.Close()
	twistr.Start(state)
	fmt.Println("Nice.")
}
