package twistr

import (
	"bufio"
	"io"
	"log"
)

type CoupLog struct {
	Country *Country
	Roll    int
}

type RealignLog struct {
	Country *Country
	RollUSA int
	RollSOV int
}

type Aof struct {
	*bufio.Scanner
	io.Writer
	done bool
}

func NewAof(r io.Reader, w io.Writer) *Aof {
	return &Aof{
		Scanner: bufio.NewScanner(r),
		Writer:  w,
		done:    false,
	}
}

func (aof *Aof) Next(thing interface{}) bool {
	if aof.done || !aof.Scan() {
		aof.done = true
		return false
	}
	line := aof.Text()
	if err := Unmarshal(line, thing); err != nil {
		log.Printf("Corrupt log! Tried to parse '%s' into %v\n", line, thing)
		return false
	}
	return true
}

func (aof *Aof) Log(thing interface{}) (err error) {
	var b []byte
	if b, err = Marshal(thing); err != nil {
		log.Println(err)
		return
	}
	if _, err = aof.Write(b); err != nil {
		log.Println(err)
		return
	}
	b = []byte{'\n'}
	if _, err = aof.Write(b); err != nil {
		log.Println(err)
		return
	}
	return
}
