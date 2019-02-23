package tf6

import (
	"fmt"
	"log"

	"github.com/jacobsa/go-serial/serial"
)

type GraphicProps struct {
	D, W, H int16
}

var alignLeft = []byte{0x1b, 0x61, 0x00}
var alignCenter = []byte{0x1b, 0x61, 0x01}
var setTitleFont = []byte{0x1d, 0x21, 0x11} // make font twice width, twice as high
var setParaFont = []byte{0x1d, 0x21, 0x01}  // make font default size

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
	ExecuteHex([]byte{0x1b, 0x2F, 0x03}, options) // set the print speed to 15 lps for better graphics
	ExecuteHex(alignCenter, options)
	commandPrefix := []byte{0x1b, 0x2a}
	wholeCommand := append(commandPrefix, []byte{byte(props.D)}...) // add doubleprint
	wholeCommand = append(wholeCommand, []byte{byte(props.W)}...)   // add width
	wholeCommand = append(wholeCommand, []byte{byte(props.H)}...)   // add height
	wholeCommand = append(wholeCommand, data...)                    // and the data
	fmt.Printf("%2x\n", wholeCommand)
	ExecuteHex(wholeCommand, options)
	ExecuteHex([]byte{0x1b, 0x2F, 0x00}, options) // set the print speed back to default 52 lps
	ExecuteHex(alignLeft, options)
}

// PrintTitle prints a center aligned title
func PrintTitle(title string, options serial.OpenOptions) {
	ExecuteHex(alignCenter, options)
	ExecuteHex(setTitleFont, options)
	PrintString(title, options)
}

// PrintParagraph prints a center aligned title, followed by a left aligned paragraph
func PrintParagraph(values string, options serial.OpenOptions) {
	ExecuteHex(alignLeft, options)
	ExecuteHex(setParaFont, options)
	PrintString(values, options)
}

// PrintTitleValues print a title with a bunch of values
func PrintTitleValues(title string, values string, options serial.OpenOptions) {
	PrintTitle(title, options)
	PrintParagraph(values, options)
}

// 576 dots per line @ 80mm

// Font sizes
// A 10 + 2 x 20
// B 12 + 2 x 24
// C 8 + 2 x 20
// D
