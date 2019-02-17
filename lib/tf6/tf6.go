package tf6

import (
	"fmt"
	"log"

	"github.com/jacobsa/go-serial/serial"
)

type GraphicProps struct {
	D, W, H int16
}

func ExecuteHex(b []byte, options serial.OpenOptions) {
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

func PrintString(inputString string, options serial.OpenOptions) {
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

func PrintGraphic(props GraphicProps, data []byte, options serial.OpenOptions) {
	commandPrefix := []byte{0x1b, 0x2a}
	wholeCommand := append(commandPrefix, []byte{byte(props.D)}...)   // add doubleprint
	wholeCommand = append(wholeCommand, []byte{byte(props.W)}...)     // add width
	wholeCommand = append(wholeCommand, []byte{byte(props.H / 8)}...) // add height
	wholeCommand = append(wholeCommand, data...)                      // and the data
	fmt.Printf("%2x\n", wholeCommand)
	ExecuteHex(wholeCommand, options)
}
