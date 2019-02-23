package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"itkettle.org/avanier/gorecptprint/lib/extras"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/jacobsa/go-serial/serial"
	"itkettle.org/avanier/gorecptprint/lib/boolbit"
	"itkettle.org/avanier/gorecptprint/lib/tf6"
)

var playSong bool
var cmdCut = []byte{0x0c}
var cmdFeed = []byte{0x1b, 0x4a, 0x20} // Print and feed paper using minimum units, 2mm in this case <p.156>
var cmdSize0 = []byte{0x1d, 0x21, 0x00}
var cmdSize1 = []byte{0x1d, 0x21, 0x01}

var options = serial.OpenOptions{
	PortName:        "/dev/ttyS0",
	BaudRate:        19200,
	DataBits:        8,
	StopBits:        1,
	MinimumReadSize: 4,
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s pemEncodedCertFile\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	if len(os.Args) <= 1 {
		usage()
	}

	if os.Getenv("PLAY_TUNE") == "false" {
		playSong = false
	} else {
		playSong = true
	}

	certData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	pairs, err := extras.CertToKVPairs(certData)

	initialize()
	for i := 0; i < len(pairs); i++ {
		tf6.PrintTitleValues(pairs[i].Key, pairs[i].Value+"\n", options)
	}
	time.Sleep(1 * time.Second)

	tf6.ExecuteHex(cmdFeed, options)
	tf6.ExecuteHex(cmdFeed, options)
	tf6.ExecuteHex(cmdFeed, options)

	stringArray := extras.SplitString(string(certData), 174)

	// zero-pad to maximum byte length for the last symbol

	if len(stringArray[len(stringArray)-1]) < 174 {
		shortString := stringArray[len(stringArray)-1]
		paddedString := shortString + "\n" + strings.Repeat(string("\x00"), (173-len(shortString)))
		stringArray[len(stringArray)-1] = paddedString
	}

	scaleFactor := 2
	symbPerLine := 2
	offset := 16 // px @ 8px/mm

	// Print x DMTX per line

	var numGroups int
	if len(stringArray)%symbPerLine != 0 {
		numGroups = (len(stringArray) / symbPerLine) + 1
	} else {
		numGroups = (len(stringArray) / symbPerLine)
	}

	for grInd := 0; grInd < numGroups; grInd++ {
		stringBatch := make([]string, symbPerLine)
		var masterRectangle image.Rectangle
		var masterImage draw.Image
		var baseWidth int

		for sInd := 0; (sInd < symbPerLine) && ((grInd*symbPerLine)+sInd < len(stringArray)); sInd++ {
			stringBatch[sInd] = stringArray[(grInd*symbPerLine)+sInd]
			oneBarcode, _ := datamatrix.Encode(stringBatch[sInd])
			oneBarcode, _ = barcode.Scale(oneBarcode, oneBarcode.Bounds().Dx()*scaleFactor, oneBarcode.Bounds().Dy()*scaleFactor)
			if sInd == 0 {
				baseWidth = oneBarcode.Bounds().Dx()
			}
			sp := image.Point{(baseWidth * sInd) + (offset * sInd), 0}    // place image starting point with offset
			dr := image.Rectangle{sp, sp.Add(oneBarcode.Bounds().Size())} // make a drawing rectangle for holding them pixels
			if sInd == 0 {
				fullWidth := (oneBarcode.Bounds().Dx() * symbPerLine) + (offset * symbPerLine)
				masterRectangle = image.Rectangle{image.ZP, image.Point{fullWidth + (1 * scaleFactor), oneBarcode.Bounds().Dy() + (1 * scaleFactor)}}
				masterImage = image.NewRGBA(masterRectangle)
				white := color.RGBA{255, 255, 255, 255}
				draw.Draw(masterImage, masterImage.Bounds(), &image.Uniform{white}, image.ZP, draw.Src) // apply primer on the whole rectangle
			}
			draw.Draw(masterImage, dr, oneBarcode, image.ZP, draw.Src)
		}

		pixels, width, height := getPixels(masterImage)

		dmtxProps := tf6.GraphicProps{D: 2, W: int16(width / 8), H: int16(height / 8)}
		tf6.PrintGraphic(dmtxProps, pixels, options)
		tf6.ExecuteHex(cmdFeed, options)
		// To prevent the buffer from exploding, we do this ugly thing which is simpler than
		// implementing flow control...
		// at 15 lps * 20 px == 300 px high per second is printed
		time.Sleep(time.Duration((height/300)+2) * time.Second) // add 1 second for wiggle room
	}

	tf6.ExecuteHex(cmdCut, options)
	if playSong {
		extras.ByeTune(options)
	}
}

func initialize() {
	var initCmds = []byte{
		0x1B, 0x40, // Reinitialize the printer <p.142>
		0x1B, 0x43, 0xFF, // Set the number of feed lines before cut to 255 (FF) steps, default 160 (A0) <p.138>
	}
	tf6.ExecuteHex(initCmds, options)
	if playSong {
		extras.ReadyTune(options)
	}
}

// Converts an Image to a list of black and white pixels
func getPixels(img image.Image) ([]byte, int, int) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	fmt.Println("width: " + strconv.Itoa(width) + " pixels")

	roundedWidth := width + (8 - width%8) // always round up to multiples of 8 pixels

	var bytePixels []byte
	for y := 0; y < height; y++ {
		var row []bool
		var stringRow string
		for x := 0; x < roundedWidth; x++ {
			if x < width {
				row = append(row, []bool{rgbaToBW(img.At(x, y).RGBA())}...)
			} else {
				row = append(row, []bool{false}...)
			}
		}
		for i := 0; i < roundedWidth/8; i++ {
			var x [8]bool
			copy(x[:], row[:8])
			row = row[8:]
			b := boolbit.BoolBit{Raw: x}
			stringRow += b.ToBin()
			bytePixels = append(bytePixels, []byte{b.ToHex()}...)
		}
		fmt.Println(stringRow)
	}

	return bytePixels, roundedWidth, height
}

func rgbaToBW(r uint32, g uint32, b uint32, a uint32) bool {
	black := bool(r == uint32(0))
	return black
}

// Check out https://github.com/grantae/certinfo
// openssl x509 -in "$2" -text -noout -certopt no_pubkey,no_sigdump
// https://commandlinefanatic.com/cgi-bin/showarticle.cgi?article=art030
// openssl asn1parse
