package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/twpayne/go-serial"
)

var defaultPath = map[string]string{
	"linux":   "/dev/ttyUSB0",
	"windows": "COM1",
}

func run() error {
	path := flag.String("p", defaultPath[runtime.GOOS], "port")
	baudRate := flag.Int("b", 57600, "baud rate")
	flag.Parse()

	port, err := serial.Open(*path, &serial.Config{
		BaudRate: *baudRate,
		DataBits: 8,
		Parity:   serial.ParityNone,
		StopBits: 1,
	})
	if err != nil {
		return err
	}
	defer port.Close()

	_, err = io.Copy(os.Stdout, port)
	return err
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
