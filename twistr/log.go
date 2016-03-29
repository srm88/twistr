package twistr

import (
	"bufio"
	"io"
	"log"
	"os"
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
	in io.ReadCloser
	io.WriteCloser
	done bool
}

func OpenAof(path string) (*Aof, error) {
	in, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		in, err = os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0666)
		if err != nil {
			return nil, err
		}
	}
	out, err2 := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	if err2 != nil {
		in.Close()
		return nil, err2
	}
	return NewAof(in, out), nil
}

func NewAof(r io.ReadCloser, w io.WriteCloser) *Aof {
	return &Aof{
		Scanner:     bufio.NewScanner(r),
		in:          r,
		WriteCloser: w,
		done:        false,
	}
}

func (aof *Aof) Close() error {
	if err := aof.in.Close(); err != nil {
		return err
	}
	return aof.Close()
}

func (aof *Aof) ReadInto(thing interface{}) bool {
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
