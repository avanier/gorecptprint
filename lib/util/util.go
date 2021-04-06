package util

import (
	"image"
	"log"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/avanier/gorecptprint/lib/boolbit"
)

// Bytes2Bits converts bytes into a human readable string of 0s and 1s.
func Bytes2Bits(data []byte) []int {
	dst := make([]int, 0)
	for _, v := range data {
		for i := 0; i < 8; i++ {
			move := uint(7 - i)
			dst = append(dst, int((v>>move)&1))
		}
	}
	// log.Println(len(dst))
	return dst
}

// IsASCII checks whether a string only contains ASCII characters.
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// SplitString splits strings into a list of strings of a maxmum length to fit
// fit within a single datamatrix 2D symbol.
func SplitString(longString string, maxLen int) []string {
	splits := []string{}

	var l, r int
	for l, r = 0, maxLen; r < len(longString); l, r = r, r+maxLen {
		for !utf8.RuneStart(longString[r]) {
			r--
		}
		splits = append(splits, longString[l:r])
	}
	splits = append(splits, longString[l:])
	return splits
}

// GetPixels converts an Image to a list of black and white pixels
func GetPixels(img image.Image) ([]byte, int, int) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	log.Println("width: " + strconv.Itoa(width) + " pixels")

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
		// log.Println(stringRow)
	}

	return bytePixels, roundedWidth, height
}

func rgbaToBW(r uint32, g uint32, b uint32, a uint32) bool {
	black := bool(r == uint32(0))
	return black
}
