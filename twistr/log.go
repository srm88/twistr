package twistr

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

type TxnLog struct {
	*bytes.Buffer
	w io.Writer
}

func NewTxnLog(w io.Writer) *TxnLog {
	return &TxnLog{
		Buffer: new(bytes.Buffer),
		w:      w,
	}
}

func (log *TxnLog) Flush() {
	// Flush buffer to writer
	log.WriteTo(log.w)
}

func Log(thing interface{}, w io.Writer) (err error) {
	var b []byte
	if b, err = Marshal(thing); err != nil {
		log.Println(err)
		return
	}
	if _, err = fmt.Fprintf(w, "%s\n", b); err != nil {
		log.Println(err)
		return
	}
	return
}

type CommandStream struct {
	*bufio.Scanner
	in   io.Reader
	done bool
}

func NewCommandStream(r io.Reader) *CommandStream {
	return &CommandStream{
		Scanner: bufio.NewScanner(r),
		in:      r,
		done:    false,
	}
}

func (cs *CommandStream) ReadInto(thing interface{}) bool {
	if cs.done || !cs.Scan() {
		cs.done = true
		return false
	}
	line := cs.Text()
	if err := Unmarshal(line, thing); err != nil {
		log.Printf("Corrupt log! Tried to parse '%s' into %s\n", line, thing)
		return false
	}
	return true
}
