package serial_test

import (
	"io"

	"github.com/twpayne/go-serial"
)

var (
	_ io.Closer = &serial.Port{}
	_ io.Reader = &serial.Port{}
	_ io.Writer = &serial.Port{}

	_ interface {
		Close() error
		Configure(*serial.Config) error
		Flush() error
		Read([]byte) (int, error)
		Write([]byte) (int, error)
	} = &serial.Port{}
)
