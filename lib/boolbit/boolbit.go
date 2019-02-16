package boolbit

type BoolBit struct {
	Raw [8]bool
}

func New(boolSlice [8]bool) *BoolBit {
	return &BoolBit{Raw: boolSlice}
}

func ToByte(bb *BoolBit) byte {
	var oneByte byte

	for i := 0; i < 8; i++ {
		if i == 0 && bb.Raw[7-i] == true {
			oneByte = 0x01
		} else {
			if bb.Raw[7-i] == true {
				oneByte += byte(2 * i)
			}
		}
	}

	return oneByte
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
