//go:build aix || linux || solaris || zos

package serial

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

var (
	standardBaudRateFlags = map[int]uint32{
		50:      unix.B50,
		75:      unix.B75,
		110:     unix.B110,
		134:     unix.B134,
		150:     unix.B150,
		200:     unix.B200,
		300:     unix.B300,
		600:     unix.B600,
		1200:    unix.B1200,
		1800:    unix.B1800,
		2400:    unix.B2400,
		4800:    unix.B4800,
		9600:    unix.B9600,
		19200:   unix.B19200,
		38400:   unix.B38400,
		57600:   unix.B57600,
		115200:  unix.B115200,
		230400:  unix.B230400,
		460800:  unix.B460800,
		500000:  unix.B500000,
		576000:  unix.B576000,
		921600:  unix.B921600,
		1000000: unix.B1000000,
		1152000: unix.B1152000,
		1500000: unix.B1500000,
		2000000: unix.B2000000,
		2500000: unix.B2500000,
		3000000: unix.B3000000,
		3500000: unix.B3500000,
		4000000: unix.B4000000,
	}

	dataBitsFlags = map[int]uint32{
		5: unix.CS5,
		6: unix.CS6,
		7: unix.CS7,
		8: unix.CS8,
	}

	parityBitsFlags = map[Parity]uint32{
		ParityNone: 0,
		ParityEven: unix.PARENB,
		ParityOdd:  unix.PARENB | unix.PARODD,
	}

	stopBitsFlags = map[StopBits]uint32{
		StopBits1: 0,
		StopBits2: unix.CSTOPB,
	}
)

// Flush flushes p.
func (p *Port) Flush() error {
	return unix.IoctlSetInt(int(p.file.Fd()), unix.TCFLSH, unix.TCIOFLUSH)
}

// Reconfigure reconfigures p.
func (p *Port) Reconfigure(config *Config) error {
	termios, err := config.termios()
	if err != nil {
		return err
	}
	if err := unix.IoctlSetTermios(int(p.file.Fd()), unix.TCSETS, &termios); err != nil {
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
		termios.Ispeed = uint32(c.BaudRate)
		termios.Ospeed = uint32(c.BaudRate)
	}

	if dataBitsFlag, ok := dataBitsFlags[c.DataBits]; ok {
		termios.Cflag |= dataBitsFlag
	} else {
		return unix.Termios{}, fmt.Errorf("%d: invalid data bits", c.DataBits)
	}

	if parityBitsFlag, ok := parityBitsFlags[c.Parity]; ok {
		termios.Cflag |= parityBitsFlag
	} else {
		return unix.Termios{}, fmt.Errorf("%d: invalid parity bits", c.Parity)
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
