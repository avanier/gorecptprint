package main

import (
	"fmt"
	"log"

	"github.com/jacobsa/go-serial/serial"
)

var cmdCut = []byte{0x0c}
var cmdSize0 = []byte{0x1d, 0x21, 0x00}
var cmdSize1 = []byte{0x1d, 0x21, 0x01}

var certString = `Certificate:

Data:
  Version: 3 (0x2)
  Serial Number:
    b7:ed:d8:49:9d:0a:35:14
  Signature Algorithm:
    sha256WithRSAEncryption
  Issuer:
    CN = openvpn.itkettle.org

  Validity
    Not Before:
      Sep  2 21:04:07 2018 GMT
    Not After :
      Aug 30 21:04:07 2028 GMT

  Subject:
    CN = openvpn.itkettle.org

  X509v3 extensions:
    X509v3 Subject Key Identifier:
      CF:5C:C7:44
      2B:E8:69:CC
      9C:28:C2:42
      8F:4D:57:FE
      07:BD:05:43

    X509v3 Authority Key Identifier:
      keyid
        CF:5C:C7:44
        2B:E8:69:CC
        9C:28:C2:42
        8F:4D:57:FE
        07:BD:05:43

      DirName:
        /CN=openvpn.itkettle.org

      serial:
        B7:ED:D8:49
        9D:0A:35:14

    X509v3 Basic Constraints:
      CA:TRUE

    X509v3 Key Usage: 
      Certificate Sign, CRL Sign
`

var options = serial.OpenOptions{
	PortName:        "/dev/ttyS0",
	BaudRate:        19200,
	DataBits:        8,
	StopBits:        1,
	MinimumReadSize: 4,
}

func main() {
	initialize()
	// someStuff()
	printString(certString)
	executeHex(cmdCut)
	byeTune()
}

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
