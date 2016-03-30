package twistr

import (
	"bufio"
	"io"
	"log"
	"os"
    "bytes"
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

type TxnLog struct {
    *bytes.Buffer
    wc io.WriteCloser
}

func NewTxnLog(wc io.WriteCloser) *TxnLog{
    return &TxnLog{
        Buffer: new(bytes.Buffer),
        wc: wc,
    }
}

func (log TxnLog) Close() error {
    return log.wc.Close()
}

func (log *TxnLog) Flush() {
    // Writes the contents of its internal buffer into the
    // WriteCloser
    log.WriteTo(log.wc)
}

func OpenTxnLog(path string) (*TxnLog, error) {
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }
    return NewTxnLog(out), nil
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
