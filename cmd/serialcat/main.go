package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/twpayne/go-serial"
)

func run() error {
	path := flag.String("p", "/dev/ttyUSB0", "port")
	baudRate := flag.Int("b", 57600, "baud rate")
	flag.Parse()

	port, err := serial.OpenAndConfigure(*path, &serial.Config{
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
