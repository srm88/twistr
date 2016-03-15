package main

import (
    "log"
    "twistr"
)

func main() {
    ui := twistr.MakeTerminalUI()
    log.Println(ui.Solicit("What is your name?", []string{"bob", "alice"}))
}

