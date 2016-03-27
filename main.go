package main

import (
	"fmt"
	"github.com/srm88/twistr/twistr"
	"log"
	"os"
)

// Temp:
func main() {
	ui := twistr.MakeTerminalUI()
	aofR, err := os.Open("/tmp/twistr.aof")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal("Failed to open AOF: ", err.Error())
		}
		aofR, err = os.OpenFile("/tmp/twistr.aof", os.O_CREATE|os.O_RDONLY, 0666)
		if err != nil {
			log.Fatal("Failed to create AOF: ", err.Error())
		}
	}
	defer aofR.Close()
	aofW, err2 := os.OpenFile("/tmp/twistr.aof", os.O_WRONLY|os.O_APPEND, 0666)
	if err2 != nil {
		log.Fatal("Failed to open AOF: ", err2.Error())
	}
	defer aofW.Close()
	state := twistr.NewState(ui, aofR, aofW)
	twistr.Start(state)
	fmt.Println("Nice.")
}
