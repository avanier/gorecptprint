package main

import (
	"image"

	"github.com/boombuler/barcode"
	"github.com/jacobsa/go-serial/serial"
	"itkettle.org/avanier/gorecptprint/lib/extras"
	"itkettle.org/avanier/gorecptprint/lib/tf6"
	boolslice "github.com/mkideal/pkg/container/boolslice"
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
	extras.PrintDummyGraphic(options)
	tf6.ExecuteHex([]byte{0x1b, 0x64, 0x02}, options) // Feed N lines
	tf6.PrintString("Hello World", options)
	tf6.ExecuteHex([]byte{0x1b, 0x64, 0x02}, options) // Feed N lines

	dmtxCode := dataMatrixCode.Encode("Hello World")
	dmtxCode, _ = barcode.Scale(dmtxCode, 432, 432) // 432 is 3 times 144

	dmtxProps := tf6.GraphicProps{d: 2, w: dmtxCode.dimension, h: dmtxCode.dimension}

	dmtxData := nil

	tf6.PrintGraphic(dmtxProps, dmtxData)

	tf6.ExecuteHex(cmdCut, options)
	extras.ByeTune(options)
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

// Converts an Image to a list of black and white pixels
func getPixels(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var boolPixels [][]byte
	for y := 0; y < height; y++ {
		row := []bool
		for x := 0; x < ((width / 8) + (( width % 8 ) + 8)) ; x++ { // always round to rows of 8 pixels
			if x < width {
				row = append(row, []bool{rgbaToPixel(img.At(x, y).RGBA())})
			}	else {
				row = append(row, []bool{false})
			}
		}
		boolPixels = append(boolPixels, row)
		pixels = append(pixels, row)
	}

	pixelList = []byte

	return pixelList, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) bool {
	return bool(r != 0) && bool(g != 0) && bool(b != 0)
}

// Check out https://github.com/grantae/certinfo
// openssl x509 -in "$2" -text -noout -certopt no_pubkey,no_sigdump
