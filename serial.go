// Package serial handles serial ports.
package serial

import (
	"errors"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

// A Parity is a parity.
type Parity int

// Parities.
const (
	ParityNone Parity = 0
	ParityOdd  Parity = 1
	ParityEven Parity = 2
)

// A Config is a serial port configuration.
type Config struct {
	BaudRate    int
	DataBits    int
	Parity      Parity
	StopBits    int
	ReadTimeout time.Duration
}

// A Port is a serial port.
type Port struct {
	file *os.File
}

// Open opens the serial path at path with the given config.
func Open(path string, config *Config) (p *Port, err error) {
	p = &Port{}
	p.file, err = os.OpenFile(path, unix.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, p.file.Close())
		}
	}()

	if err := p.Reconfigure(config); err != nil {
		return nil, err
	}

	return p, nil
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
