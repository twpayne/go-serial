package serial

// See https://learn.microsoft.com/en-us/windows/win32/devio/communications-resources.
// See https://learn.microsoft.com/en-us/previous-versions/ms810467(v=msdn.10).

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Extra Windows parities.
const (
	ParityMark  Parity = 3
	ParitySpace Parity = 4
)

// Extra Windows stop bits.
const (
	StopBits1Point5 StopBits = 15
)

// DCB flags.
const (
	fBinary           = 0x01
	fParity           = 0x02
	fOutxCtsFlow      = 0x04
	fOutxDsrFlow      = 0x08
	fDtrControl       = 0x30
	fDsrSensitivity   = 0x40
	fTXContinueOnXoff = 0x80
	fOutX             = 0x01
	fInX              = 0x02
	fErrorChar        = 0x04
	fNull             = 0x08
	fRtsControl       = 0x30
	fAbortOnError     = 0x40
)

var (
	parities = map[Parity]uint8{
		ParityNone:  windows.NOPARITY,
		ParityOdd:   windows.ODDPARITY,
		ParityEven:  windows.EVENPARITY,
		ParityMark:  windows.MARKPARITY,
		ParitySpace: windows.SPACEPARITY,
	}

	stopBits = map[StopBits]uint8{
		StopBits1:       windows.ONESTOPBIT,
		StopBits1Point5: windows.ONE5STOPBITS,
		StopBits2:       windows.TWOSTOPBITS,
	}
)

type Port struct {
	handle windows.Handle
	file   *os.File
}

// Open opens the serial path at path with the given config.
func Open(path string, config *Config) (p *Port, err error) {
	p = &Port{}
	if !strings.HasPrefix(path, `\`) {
		path = `\\.\` + path
	}
	name, err := windows.UTF16FromString(path)
	if err != nil {
		return nil, err
	}
	p.handle, err = windows.CreateFile(&name[0], windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return nil, err
	}
	p.file = os.NewFile(uintptr(p.handle), path)
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

// Flush flushes p.
func (p *Port) Flush() error {
	if err := windows.PurgeComm(p.handle, windows.PURGE_TXABORT|windows.PURGE_RXABORT|windows.PURGE_TXCLEAR|windows.PURGE_RXCLEAR); err != nil {
		return err
	}
	return nil
}

// Reconfigure reconfigures p.
func (p *Port) Reconfigure(config *Config) error {
	// See https://learn.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-dcb.
	var dcb windows.DCB
	dcb.DCBlength = uint32(unsafe.Sizeof(dcb))
	dcb.Flags[0] = fBinary
	dcb.Flags[0] |= windows.DTR_CONTROL_ENABLE << 4
	dcb.BaudRate = uint32(config.BaudRate)
	dcb.ByteSize = uint8(config.DataBits)
	if parity, ok := parities[config.Parity]; ok {
		dcb.Parity = parity
	} else {
		return fmt.Errorf("%d: invalid parity", config.Parity)
	}
	if stopBits, ok := stopBits[config.StopBits]; ok {
		dcb.StopBits = stopBits
	} else {
		return fmt.Errorf("%d: invalid stop bits", config.StopBits)
	}
	if err := windows.SetCommState(p.handle, &dcb); err != nil {
		return err
	}

	commTimeouts := windows.CommTimeouts{
		ReadIntervalTimeout:        math.MaxUint32,
		ReadTotalTimeoutMultiplier: math.MaxUint32,
		ReadTotalTimeoutConstant:   uint32(config.ReadTimeout / time.Millisecond), // FIXME clamp
	}
	if err := windows.SetCommTimeouts(p.handle, &commTimeouts); err != nil {
		return err
	}

	if err := windows.SetCommMask(p.handle, windows.EV_RXCHAR); err != nil {
		return err
	}

	return nil
}
