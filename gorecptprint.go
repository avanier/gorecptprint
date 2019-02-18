package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"strings"

	"itkettle.org/avanier/gorecptprint/lib/extras"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"github.com/jacobsa/go-serial/serial"
	"itkettle.org/avanier/gorecptprint/lib/boolbit"
	"itkettle.org/avanier/gorecptprint/lib/tf6"
)

var cmdCut = []byte{0x0c}
var cmdFeed = []byte{0x1b, 0x4a, 0x10} // Print and feed paper using minimum units, 2mm in this case <p.156>
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
	daString := `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Eu volutpat odio facilisis mauris sit amet massa vitae. Pellentesque habitant morbi tristique senectus et netus. Sed lectus vestibulum mattis ullamcorper velit. Vitae sapien pellentesque habitant morbi tristique senectus et. Ac turpis egestas integer eget aliquet nibh praesent tristique. Aliquet eget sit amet tellus cras adipiscing enim. Netus et malesuada fames ac. Eget sit amet tellus cras adipiscing enim. Elit eget gravida cum sociis natoque penatibus et magnis. Diam volutpat commodo sed egestas egestas. Diam quam nulla porttitor massa. Condimentum lacinia quis vel eros donec ac odio. Eget duis at tellus at urna condimentum. Pharetra massa massa ultricies mi quis hendrerit dolor magna. Lectus proin nibh nisl condimentum id venenatis a condimentum vitae. Fames ac turpis egestas integer eget aliquet nibh. Faucibus pulvinar elementum integer enim neque.

	Nec sagittis aliquam malesuada bibendum arcu vitae. Ipsum nunc aliquet bibendum enim facilisis gravida neque. Egestas diam in arcu cursus euismod. Metus aliquam eleifend mi in nulla posuere sollicitudin aliquam ultrices. Mattis ullamcorper velit sed ullamcorper morbi tincidunt ornare. Duis ultricies lacus sed turpis tincidunt. Eget felis eget nunc lobortis mattis aliquam faucibus purus in. Proin fermentum leo vel orci porta. Eget dolor morbi non arcu risus quis varius quam quisque. Rhoncus aenean vel elit scelerisque mauris. Imperdiet massa tincidunt nunc pulvinar sapien et ligula. Tellus integer feugiat scelerisque varius morbi enim. Sem et tortor consequat id porta nibh venenatis cras. Aliquam sem fringilla ut morbi. Tellus orci ac auctor augue mauris. Lectus mauris ultrices eros in cursus turpis massa tincidunt. Senectus et netus et malesuada. Proin nibh nisl condimentum id venenatis a.

	Aliquam purus sit amet luctus venenatis lectus magna fringilla urna. Volutpat lacus laoreet non curabitur gravida. Nulla pellentesque dignissim enim sit amet venenatis urna cursus eget. At in tellus integer feugiat scelerisque varius morbi enim nunc. A scelerisque purus semper eget duis at. Sed lectus vestibulum mattis ullamcorper. Scelerisque purus semper eget duis at. Ut porttitor leo a diam sollicitudin. Sit amet aliquam id diam maecenas ultricies. Scelerisque mauris pellentesque pulvinar pellentesque habitant morbi tristique. Posuere lorem ipsum dolor sit amet. Arcu vitae elementum curabitur vitae nunc sed velit dignissim. In est ante in nibh mauris cursus mattis molestie. In ante metus dictum at tempor commodo. Eu augue ut lectus arcu bibendum at varius vel. Eu scelerisque felis imperdiet proin fermentum leo vel. Sagittis nisl rhoncus mattis rhoncus urna. Dui faucibus in ornare quam viverra orci sagittis eu volutpat. Elit duis tristique sollicitudin nibh sit amet commodo nulla.
	`
	// daString := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Eu volutpat odio facilisis mauris sit amet massa vitae. Pellentesque habitant morbi tristique senectus et netus. Sed lectus vestibulum mattis ullamcorper velit. Vitae sapien pellentesque habitant morbi tristique senectus et. Ac turpis egestas integer eget aliquet nibh praesent tristique. Aliquet eget sit amet tellus cras adipiscing enim. Netus et malesuada fames ac. Eget sit amet tellus cras adipiscing enim. Elit eget gravida cum sociis natoque penatibus et magnis. Diam volutpat commodo sed egestas egestas. Diam quam nulla porttitor massa. Condimentum lacinia quis vel eros donec ac odio. Eget duis at tellus at urna condimentum. Pharetra massa massa ultricies mi quis hendrerit dolor magna. Lectus proin nibh nisl condimentum id venenatis a condimentum vitae. Fames ac turpis egestas integer eget aliquet nibh. Faucibus pulvinar elementum integer enim neque."
	stringArray := extras.SplitString(daString, 174)

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
				masterRectangle = image.Rectangle{image.ZP, image.Point{fullWidth, oneBarcode.Bounds().Dy()}}
				masterImage = image.NewRGBA(masterRectangle)
				white := color.RGBA{255, 255, 255, 255}
				draw.Draw(masterImage, masterImage.Bounds(), &image.Uniform{white}, image.ZP, draw.Src) // apply primer on the whole rectangle
			}
			draw.Draw(masterImage, dr, oneBarcode, image.ZP, draw.Src)
		}

		tf6.PrintString(strings.Join(stringBatch, ""), options)
		pixels, width, height := getPixels(masterImage)

		dmtxProps := tf6.GraphicProps{D: 2, W: int16(width / 8), H: int16(height / 8)}
		tf6.PrintGraphic(dmtxProps, pixels, options)
		tf6.ExecuteHex(cmdFeed, options)
	}

	// for i := 0; i < len(stringArray); i++ {
	// 	oneCode := stringArray[i]

	// 	tf6.PrintString(oneCode, options)
	// 	dmtxCode, _ := datamatrix.Encode(oneCode)
	// 	dmtxCode, _ = barcode.Scale(dmtxCode, dmtxCode.Bounds().Max.X*scaleFactor, dmtxCode.Bounds().Max.Y*scaleFactor)

	// 	pixels, width, height := getPixels(dmtxCode)

	// 	dmtxProps := tf6.GraphicProps{D: 2, W: int16(width / 8), H: int16(height / 8)}
	// 	tf6.PrintGraphic(dmtxProps, pixels, options)
	// 	tf6.ExecuteHex(cmdFeed, options)
	// }

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
	if black == true {
		// fmt.Println("pixel is black")
	} else {
		// fmt.Println("pixel is white")
	}
	return black
}

// Check out https://github.com/grantae/certinfo
// openssl x509 -in "$2" -text -noout -certopt no_pubkey,no_sigdump
