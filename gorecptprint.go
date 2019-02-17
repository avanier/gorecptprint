package main

import (
	"fmt"
	"image"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/jacobsa/go-serial/serial"
	"itkettle.org/avanier/gorecptprint/lib/boolbit"
	"itkettle.org/avanier/gorecptprint/lib/tf6"
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
	initialize()
	// extras.PrintDummyGraphic(options)
	// tf6.ExecuteHex([]byte{0x1b, 0x64, 0x02}, options) // Feed N lines
	tf6.PrintString("Hello World", options)
	// tf6.ExecuteHex([]byte{0x1b, 0x64, 0x02}, options) // Feed N lines

	dmtxCode, _ := datamatrix.Encode("Hello World")
	dmtxCode, _ = barcode.Scale(dmtxCode, dmtxCode.Bounds().Max.X*4, dmtxCode.Bounds().Max.Y*4)

	pixels, byteWidth, height := getPixels(dmtxCode)

	fmt.Println(pixels)
	fmt.Printf("%2x\n", pixels)
	fmt.Println(byteWidth, height)

	dmtxProps := tf6.GraphicProps{D: 2, W: int16(byteWidth), H: int16(height)}
	tf6.PrintGraphic(dmtxProps, pixels, options)

	tf6.ExecuteHex(cmdCut, options)
	// extras.ByeTune(options)
}

// Data buffer on the printer is 16KB
// Check out pages 115 and 157 for uploading and printing pixels

func initialize() {
	var initCmds = []byte{
		0x1B, 0x40, // Reinitialize the printer <p.142>
		0x1B, 0x43, 0xFF, // Set the number of feed lines before cut to 255 (FF) steps, default 160 (A0) <p.138>
	}
	tf6.ExecuteHex(initCmds, options)
	// extras.ReadyTune(options)
}

// Converts an Image to a list of black and white pixels
func getPixels(img image.Image) ([]byte, int, int) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	fmt.Println("width: " + strconv.Itoa(width) + " pixels")

	var byteWidth int
	var bytePixels []byte
	for y := 0; y < height; y++ {
		var row []bool
		var stringRow string
		if width%8 == 0 {
			byteWidth = width / 8
		} else {
			byteWidth = (width / 8) + 1
		}
		// all of this is double conversion and should be merged with boolSlice2byteSlice
		for x := 0; x < byteWidth; x++ { // always round up to rows of 8 pixels
			// fmt.Println("parsing byte: " + strconv.Itoa(x))
			for a := 0; a < 8; a++ {
				// fmt.Println("parsing pixel x:"+strconv.Itoa((x*8)+a)+", y:", y, "to bit")
				if x < width/8 {
					row = append(row, []bool{rgbaToBW(img.At(x+a, y).RGBA())}...)
				} else {
					row = append(row, []bool{false}...)
				}
			}
		}
		for i := 0; i < byteWidth; i++ {
			var x [8]bool
			copy(x[:], row[:8])
			row = row[8:]
			b := boolbit.BoolBit{Raw: x}
			stringRow += b.ToBin()
			bytePixels = append(bytePixels, []byte{b.ToHex()}...)
		}
		fmt.Println(stringRow)
	}

	return bytePixels, byteWidth, height
}

func rgbaToBW(r uint32, g uint32, b uint32, a uint32) bool {
	black := bool(r == uint32(0))
	if black == true {
		// fmt.Println("pixel is black")
	} else {
		// fmt.Println("pixel is white")
	}
	return black
}

// Check out https://github.com/grantae/certinfo
// openssl x509 -in "$2" -text -noout -certopt no_pubkey,no_sigdump
