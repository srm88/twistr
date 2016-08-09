package twistr

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

type CmdOut struct {
	*bytes.Buffer
	w io.Writer
}

func NewCmdOut(w io.Writer) *CmdOut {
	return &CmdOut{
		Buffer: new(bytes.Buffer),
		w:      w,
	}
}

func (co *CmdOut) Flush() {
	// XXX errors
	// Flush buffer to writer
	co.WriteTo(co.w)
}

func (co *CmdOut) Log(thing interface{}) (err error) {
	var b []byte
	if b, err = Marshal(thing); err != nil {
		log.Println(err)
		return
	}
	if _, err = fmt.Fprintf(co, "%s\n", b); err != nil {
		log.Println(err)
		return
	}
	return
}

type CmdIn struct {
	*bufio.Scanner
	in   io.Reader
	done bool
}

func NewCmdIn(r io.Reader) *CmdIn {
	return &CmdIn{
		Scanner: bufio.NewScanner(r),
		in:      r,
		done:    false,
	}
}

func (ci *CmdIn) ReadInto(thing interface{}) bool {
	if ci.done || !ci.Scan() {
		ci.done = true
		return false
	}
	line := ci.Text()
	if err := Unmarshal(line, thing); err != nil {
		log.Printf("Corrupt log! Tried to parse '%s' into %s\n", line, thing)
		return false
	}
	return true
}
