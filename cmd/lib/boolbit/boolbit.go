package boolbit

import "math"

type BoolBit struct {
	Raw [8]bool
}

func (bb *BoolBit) ToHex() byte {
	var doubleByte byte

	for i := 0; i < 8; i++ {
		if i == 0 && bb.Raw[7-i] == true {
			doubleByte = 0x01
		} else {
			if bb.Raw[7-i] == true {
				doubleByte += byte(math.Pow(2, float64(i)))
			}
		}
	}

	return doubleByte
}

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
