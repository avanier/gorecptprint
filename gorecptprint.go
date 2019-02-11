package main

import (
	"fmt"
	"image"
	"strconv"

	"github.com/boombuler/barcode/datamatrix"
	"github.com/jacobsa/go-serial/serial"
	"itkettle.org/avanier/gorecptprint/lib/extras"
	"itkettle.org/avanier/gorecptprint/lib/tf6"

	"github.com/Workiva/go-datastructures/bitarray"
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
	// extras.PrintDummyGraphic(options)
	// tf6.ExecuteHex([]byte{0x1b, 0x64, 0x02}, options) // Feed N lines
	// tf6.PrintString("Hello World", options)
	// tf6.ExecuteHex([]byte{0x1b, 0x64, 0x02}, options) // Feed N lines

	dmtxCode, _ := datamatrix.Encode("Hello World")
	// dmtxCode, _ = barcode.Scale(dmtxCode, 432, 432) // 432 is 3 times 144

	// bounds := dmtxCode.Bounds()

	// dmtxProps := tf6.GraphicProps{
	// 	D: 2,
	// 	W: int16(bounds.Max.X),
	// 	H: int16(bounds.Max.Y),
	// }

	pixels, _ := getPixels(dmtxCode)

	fmt.Println(pixels)

	// var dmtxData []byte

	// tf6.PrintGraphic(dmtxProps, dmtxData)

	// tf6.ExecuteHex(cmdCut, options)
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
	extras.ReadyTune(options)
}

// NOTE TO SELF
// Implement bwImage with dimensions automatically conformed to nearest 8 pixels
// Allow it to take an Image type for conversion
// Make it have dimensions properties and a hexdump instance method

// Converts an Image to a list of black and white pixels
func getPixels(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	fmt.Println("width: " + strconv.Itoa(width) + " pixels")

	var bytePixels []byte
	for y := 0; y < height; y++ {
		var row []bool
		var byteWidth int
		if width%8 == 0 {
			byteWidth = width / 8
		} else {
			byteWidth = (width / 8) + 1
		}
		// all of this is double conversion and should be merged with boolSlice2byteSlice
		for x := 0; x < byteWidth; x++ { // always round up to rows of 8 pixels
			fmt.Println("parsing byte: " + strconv.Itoa(x))
			for a := 0; a < 8; a++ {
				fmt.Println("parsing pixel " + strconv.Itoa(a) + " to bit")
				if x < width/8 {
					row = append(row, []bool{rgbaToBW(img.At(x, y).RGBA())}...)
				} else {
					row = append(row, []bool{false}...)
				}
			}
		}
		byteList := boolSlice2byteSlice(row)
		bytePixels = append(bytePixels, byteList...)
	}

	var pixelList []byte

	return pixelList, nil
}

func boolSlice2byteSlice(s []bool) []byte {
	var c []byte
	fmt.Println("length of current boolSlice:", strconv.Itoa(len(s)), "bits")
	for y := 0; y < (len(s) / 8); y++ {
		fmt.Println("current slice range " + strconv.Itoa(y))
		var x []bool
		if y == 0 {
			fmt.Println("I'm in the first conditional")
			x, s = s[0:7], s[8:]
		} else if y < 7 {
			fmt.Println("I'm in the second conditional")
			x, s = s[y:(y*8)-1], s[(y*8):]
		} else {
			x = s[y : (y*8)-1]
		}
		var b = bitarray.NewBitArray(8)
		for z := 0; z < 8; z++ {
			fmt.Println("trying bit", strconv.Itoa(z), "of", len(x))
			if x[z] == true {
				fmt.Println("setting bit", strconv.Itoa(z), "to 1")
				b.SetBit(uint64(z))
			} else {
				fmt.Println("setting bit", strconv.Itoa(z), "to 0")
			}
		}
		d, _ := bitarray.Marshal(b)
		fmt.Println("resulting byte: ", d)
		c = append(c, d...)
	}
	return c
}

func rgbaToBW(r uint32, g uint32, b uint32, a uint32) bool {
	return bool(r != 0) && bool(g != 0) && bool(b != 0)
}

// Check out https://github.com/grantae/certinfo
// openssl x509 -in "$2" -text -noout -certopt no_pubkey,no_sigdump
