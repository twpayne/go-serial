//go:build !windows

package serial

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

// A Port is a serial port.
type Port struct {
	file *os.File
}

// Open opens the serial port at path.
func Open(path string) (*Port, error) {
	file, err := os.OpenFile(path, unix.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		return nil, err
	}
	return &Port{
		file: file,
	}, nil
}

// OpenAndConfigure opens the serial port at path and configures it with the
// given config.
func OpenAndConfigure(path string, config *Config) (*Port, error) {
	p, err := Open(path)
	if err != nil {
		return nil, err
	}

	if err := p.Configure(config); err != nil {
		return nil, errors.Join(err, p.Close())
	}

	return p, nil
}
