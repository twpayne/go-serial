// Package serial handles serial ports.
package serial

import (
	"time"
)

// A Parity is a parity.
type Parity int

// Parities.
const (
	ParityNone Parity = 0
	ParityOdd  Parity = 1
	ParityEven Parity = 2
)

// A StopBits is the number of stop bits.
type StopBits int

// Stop bits.
const (
	StopBits1 StopBits = 1
	StopBits2 StopBits = 2
)

// A Config is a serial port configuration.
type Config struct {
	BaudRate    int
	DataBits    int
	Parity      Parity
	StopBits    StopBits
	ReadTimeout time.Duration
}

// Close closes p.
func (p *Port) Close() error {
	return p.file.Close()
}

// Read reads up to len(data) bytes from p and stores them in data. It returns
// the number of bytes read and any error encountered.
func (p *Port) Read(data []byte) (int, error) {
	return p.file.Read(data)
}

// Write writes len(data) bytes from data to the port. It returns the number of
// bytes written and any error. Write returns a non-nil error when n != len(b).
func (p *Port) Write(data []byte) (int, error) {
	return p.file.Write(data)
}
