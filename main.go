package main

import (
    "fmt"
    "strings"
    "twistr"
)

// Temp:
func main() {
    ui := twistr.MakeTerminalUI()
    cname := ui.Solicit("Adj to?", []string{"SKorea", "Mexico", "etc"})
    c := twistr.ByName(strings.TrimSpace(cname))
    fmt.Println(twistr.CountryNames(c.AdjCountries))
}

