package boolbit

import "math"

// BoolBit purpose is to simply express a black and white bitmaps line of eight
// pixels. Handy for working with monochrome serial printers.
type BoolBit struct {
	Raw [8]bool
}

// ToHex returns a single byte representing a 8 pixels wide chunk of image.
func (bb *BoolBit) ToHex() byte {
	var oneByte byte

	for i := 0; i < 8; i++ {
		if i == 0 && bb.Raw[7-i] == true {
			oneByte = 0x01
		} else {
			if bb.Raw[7-i] == true {
				oneByte += byte(math.Pow(2, float64(i)))
			}
		}
	}

	return oneByte
}

// ToBin returns a human readable representation of 0s and 1s in a string. Handy
// for "previewing" your image in a terminal.
func (bb *BoolBit) ToBin() string {
	var binString string

	for i := 0; i < 8; i++ {
		if bb.Raw[i] == true {
			binString = binString + "1"
		} else {
			binString = binString + "0"
		}
	}

	return binString
}
