package main

import (
	"fmt"
	"log"

	"github.com/jacobsa/go-serial/serial"
	dmtx "itkettle.org/avanier/gorecptprint/lib"
)

var cmdCut = []byte{0x0c}
var cmdSize0 = []byte{0x1d, 0x21, 0x00}
var cmdSize1 = []byte{0x1d, 0x21, 0x01}

var options = serial.OpenOptions{
	PortName:        "/dev/ttyS0",
	BaudRate:        19200,
	DataBits:        8,
	StopBits:        1,
	MinimumReadSize: 4,
}

func main() {
	// initialize()
	// // someStuff()
	// printString(certString)
	// executeHex(cmdCut)
	// byeTune()
	err, data := dmtx.GenDMXT()

	if err != nil {
		panic(err)
	}

	// fmt.Println(buf.Bytes())
	// pixelarray.Test()

}

// Data buffer on the printer is 16KB
// Check out pages 115 and 157 for uploading and printing pixels

func initialize() {
	var initCmds = []byte{
		0x1B, 0x40, // Reinitialize the printer <p.142>
		0x1B, 0x43, 0xFF, // Set the number of feed lines before cut to 255 (FF) steps, default 160 (A0) <p.138>
	}
	executeHex(initCmds)
	readyTune()
}

func executeHex(b []byte) {
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	defer port.Close()

	_, err = port.Write(b)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
}

func readyTune() {
	// Plays some beeps to signal end of initialization
	// See <p.144>
	var readyTune = []byte{
		0x1b, 0x07, // Start the sequence
		0x02, // Set the duration from 01 - FF times 0.1 seconds
		0x90, // Binary conversion of 10010000 - (10)<soft>(01)<octave 2>(0000)<note c>
		0x1b, 0x07,
		0x01,
		0x95,
		0x1b, 0x07,
		0x01,
		0x99,
	}
	executeHex(readyTune)
}

func byeTune() {
	var readyTune = []byte{
		0x1b, 0x07,
		0x02,
		0x9a,
		0x1b, 0x07,
		0x01,
		0x99,
		0x1b, 0x07,
		0x01,
		0x95,
	}
	executeHex(readyTune)
}

func someStuff() {
	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	b := append(cmdSize1, []byte("This is a string.\n")...)
	b = append(b, cmdCut...)

	n, err := port.Write(b)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}

	fmt.Println("Wrote", n, "bytes.")
}

func printString(inputString string) {
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	defer port.Close()

	_, err = port.Write([]byte(inputString))
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
}

// Check out https://github.com/grantae/certinfo
// openssl x509 -in "$2" -text -noout -certopt no_pubkey,no_sigdump
