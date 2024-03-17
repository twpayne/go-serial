//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package serial

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

var (
	standardBaudRateFlags = map[int]uint64{
		50:     unix.B50,
		75:     unix.B75,
		110:    unix.B110,
		134:    unix.B134,
		150:    unix.B150,
		200:    unix.B200,
		300:    unix.B300,
		600:    unix.B600,
		1200:   unix.B1200,
		1800:   unix.B1800,
		2400:   unix.B2400,
		4800:   unix.B4800,
		9600:   unix.B9600,
		19200:  unix.B19200,
		38400:  unix.B38400,
		57600:  unix.B57600,
		115200: unix.B115200,
		230400: unix.B230400,
	}

	dataBitsFlags = map[int]uint64{
		5: unix.CS5,
		6: unix.CS6,
		7: unix.CS7,
		8: unix.CS8,
	}

	parityBitsFlags = map[Parity]uint64{
		ParityNone: 0,
		ParityEven: unix.PARENB,
		ParityOdd:  unix.PARENB | unix.PARODD,
	}

	stopBitsFlags = map[StopBits]uint64{
		StopBits1: 0,
		StopBits2: unix.CSTOPB,
	}
)

// Flush flushes p.
func (p *Port) Flush() error {
	return unix.IoctlSetInt(int(p.file.Fd()), unix.TIOCFLUSH, unix.TCIOFLUSH)
}

// Reconfigure reconfigures p.
func (p *Port) Reconfigure(config *Config) error {
	termios, err := config.termios()
	if err != nil {
		return err
	}
	if err := unix.IoctlSetTermios(int(p.file.Fd()), unix.TIOCSETA, &termios); err != nil {
		return err
	}
	return nil
}

// termios returns the unix.Termios structure for c.
func (c Config) termios() (unix.Termios, error) {
	termios := unix.Termios{
		Iflag: unix.IGNPAR,
		Cflag: unix.CREAD | unix.CLOCAL,
	}

	if baudRate, ok := standardBaudRateFlags[c.BaudRate]; ok {
		termios.Ispeed = baudRate
		termios.Ospeed = baudRate
	} else {
		termios.Ispeed = uint64(c.BaudRate)
		termios.Ospeed = uint64(c.BaudRate)
	}

	if dataBitsFlag, ok := dataBitsFlags[c.DataBits]; ok {
		termios.Cflag |= dataBitsFlag
	} else {
		return unix.Termios{}, fmt.Errorf("%d: invalid data bits", c.DataBits)
	}

	if parityBitsFlag, ok := parityBitsFlags[c.Parity]; ok {
		termios.Cflag |= parityBitsFlag
	} else {
		return unix.Termios{}, fmt.Errorf("%d: invalid parity", c.Parity)
	}

	if stopBitsFlag, ok := stopBitsFlags[c.StopBits]; ok {
		termios.Cflag |= stopBitsFlag
	} else {
		return unix.Termios{}, fmt.Errorf("%d: invalid stop bits", c.StopBits)
	}

	if c.ReadTimeout > 0 {
		readTimeoutDeciSeconds := int(c.ReadTimeout / (100 * time.Millisecond))
		termios.Cc[unix.VTIME] = uint8(max(1, min(readTimeoutDeciSeconds, 255)))
	} else {
		termios.Cc[unix.VMIN] = 1
	}

	return termios, nil
}
