package twistr

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

type Aof struct {
	*bufio.Scanner
	in io.ReadCloser
	io.WriteCloser
	exhaust io.Writer
	done    bool
}

type TxnLog struct {
	*bytes.Buffer
	wc io.WriteCloser
}

func NewTxnLog(wc io.WriteCloser) *TxnLog {
	return &TxnLog{
		Buffer: new(bytes.Buffer),
		wc:     wc,
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

func NewAof(r io.ReadCloser, w io.WriteCloser, e io.Writer) *Aof {
	return &Aof{
		Scanner:     bufio.NewScanner(r),
		in:          r,
		WriteCloser: w,
		exhaust:     e,
		done:        false,
	}
}

func (aof *Aof) Close() error {
	if err := aof.in.Close(); err != nil {
		return err
	}
	return aof.Close()
}

func (aof *Aof) Next() (bool, string) {
	if aof.done || !aof.Scan() {
		aof.done = true
		return false, ""
	}
	line := aof.Text()
	aof.exhaust.Write([]byte(line))
	return true, line
}

func (aof *Aof) Log(thing interface{}) (err error) {
	var b []byte
	if b, err = Marshal(thing); err != nil {
		log.Println(err)
		return
	}
	if _, err = aof.Write(append(b, '\n')); err != nil {
		log.Println(err)
		return
	}
	aof.exhaust.Write(b)
	return
}
