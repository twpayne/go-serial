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
